package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"

	domain "antblog/internal/domain/user"
	"antblog/pkg/crypto"
	apperrors "antblog/pkg/errors"
	"antblog/pkg/jwt"
)

// ─── 依赖声明 ────────────────────────────────────────────────────────────────

// Deps fx 注入依赖
type Deps struct {
	fx.In
	Repo          domain.IUserRepository
	DomainService domain.IDomainService
	TokenManager  jwt.ITokenManager
	Logger        *zap.Logger
}

// userUseCase 用户用例实现
type userUseCase struct {
	repo   domain.IUserRepository
	svc    domain.IDomainService
	token  jwt.ITokenManager
	logger *zap.Logger
}

// NewUserUseCase 创建用户用例（fx provider）
func NewUserUseCase(deps Deps) IUserUseCase {
	return &userUseCase{
		repo:   deps.Repo,
		svc:    deps.DomainService,
		token:  deps.TokenManager,
		logger: deps.Logger,
	}
}

// ─── Register ────────────────────────────────────────────────────────────────

func (uc *userUseCase) Register(ctx context.Context, req *RegisterReq) (*UserResp, error) {
	// 1. 领域校验：值对象合法性 + 唯一性
	if err := uc.svc.ValidateRegister(ctx, req.Username, req.Email); err != nil {
		return nil, err
	}

	// 2. 密码强度校验
	if _, err := domain.NewPassword(req.Password); err != nil {
		return nil, apperrors.ErrInvalidParams(err.Error())
	}

	// 3. 密码哈希
	hashed, err := crypto.HashPassword(req.Password)
	if err != nil {
		uc.logger.Error("hash password failed", zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}

	// 4. 构建领域实体
	u, err := uc.svc.BuildNewUser(req.Username, req.Email, hashed)
	if err != nil {
		return nil, err
	}
	u.UUID = uuid.NewString()

	// 5. 持久化
	created, err := uc.repo.Create(ctx, u)
	if err != nil {
		uc.logger.Error("create user failed", zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}

	uc.logger.Info("user registered", zap.Uint64("user_id", created.ID), zap.String("username", created.Username))
	return toUserResp(created), nil
}

// ─── Login ───────────────────────────────────────────────────────────────────

func (uc *userUseCase) Login(ctx context.Context, req *LoginReq, userAgent, clientIP string) (*TokenResp, error) {
	// 1. 领域验证凭据
	u, err := uc.svc.ValidateCredentials(ctx, req.Email, req.Password, crypto.CheckPassword)
	if err != nil {
		return nil, err
	}

	// 2. 生成 Token 对
	resp, err := uc.issueTokenPair(ctx, u, userAgent, clientIP)
	if err != nil {
		return nil, err
	}

	// 3. 更新最后登录时间（异步，不影响主流程）
	go func() {
		if err := uc.repo.UpdateLastLogin(context.Background(), u.ID); err != nil {
			uc.logger.Warn("update last login failed", zap.Uint64("user_id", u.ID), zap.Error(err))
		}
	}()

	uc.logger.Info("user logged in", zap.Uint64("user_id", u.ID), zap.String("ip", clientIP))
	return resp, nil
}

// ─── Logout ──────────────────────────────────────────────────────────────────

func (uc *userUseCase) Logout(ctx context.Context, refreshToken string) error {
	if err := uc.repo.DeleteToken(ctx, refreshToken); err != nil {
		uc.logger.Warn("delete token failed", zap.Error(err))
		// 登出失败不报错，保证幂等
	}
	return nil
}

// ─── RefreshToken ────────────────────────────────────────────────────────────

func (uc *userUseCase) RefreshToken(ctx context.Context, req *RefreshTokenReq, userAgent, clientIP string) (*TokenResp, error) {
	// 1. 解析 Refresh Token
	claims, err := uc.token.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 2. 数据库校验 Token 是否存在（防重放）
	tokenRecord, err := uc.repo.FindToken(ctx, req.RefreshToken)
	if err != nil || tokenRecord == nil {
		return nil, apperrors.ErrTokenInvalid()
	}
	if tokenRecord.IsExpired() {
		_ = uc.repo.DeleteToken(ctx, req.RefreshToken)
		return nil, apperrors.FromCode(apperrors.CodeRefreshTokenExpired)
	}

	// 3. 查询用户
	u, err := uc.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound()
	}
	if !u.IsActive() {
		return nil, apperrors.ErrUserDisabled()
	}

	// 4. 旧 Token 轮换（删除旧的，签发新的）
	_ = uc.repo.DeleteToken(ctx, req.RefreshToken)
	return uc.issueTokenPair(ctx, u, userAgent, clientIP)
}

// ─── GetProfile ──────────────────────────────────────────────────────────────

func (uc *userUseCase) GetProfile(ctx context.Context, userID uint64) (*UserResp, error) {
	u, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound()
	}
	return toUserResp(u), nil
}

// ─── UpdateProfile ───────────────────────────────────────────────────────────

func (uc *userUseCase) UpdateProfile(ctx context.Context, userID uint64, req *UpdateProfileReq) (*UserResp, error) {
	u, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound()
	}

	u.UpdateProfile(req.Nickname, req.Avatar, req.Bio)

	if err = uc.repo.Update(ctx, u); err != nil {
		uc.logger.Error("update profile failed", zap.Uint64("user_id", userID), zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}
	return toUserResp(u), nil
}

// ─── ChangePassword ──────────────────────────────────────────────────────────

func (uc *userUseCase) ChangePassword(ctx context.Context, userID uint64, req *ChangePasswordReq) error {
	u, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		return apperrors.ErrUserNotFound()
	}

	// 验证旧密码
	if !crypto.CheckPassword(req.OldPassword, u.Password) {
		return apperrors.ErrPasswordIncorrect()
	}

	// 校验新密码强度
	if _, err := domain.NewPassword(req.NewPassword); err != nil {
		return apperrors.ErrInvalidParams(err.Error())
	}

	// 新旧密码不能相同
	if crypto.CheckPassword(req.NewPassword, u.Password) {
		return apperrors.ErrInvalidParams("新密码不能与旧密码相同")
	}

	hashed, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		return apperrors.ErrInternalError(err)
	}
	u.Password = hashed

	if err = uc.repo.Update(ctx, u); err != nil {
		return apperrors.ErrInternalError(err)
	}

	// 修改密码后强制所有设备下线
	_ = uc.repo.DeleteUserTokens(ctx, userID)
	uc.logger.Info("user changed password", zap.Uint64("user_id", userID))
	return nil
}

// ─── 内部辅助 ────────────────────────────────────────────────────────────────

// issueTokenPair 签发 Access + Refresh Token 并持久化 Refresh Token
func (uc *userUseCase) issueTokenPair(ctx context.Context, u *domain.User, userAgent, clientIP string) (*TokenResp, error) {
	// 生成 Access Token
	claims := &jwt.UserClaims{
		UserID:   u.ID,
		Username: u.Username,
		Role:     int(u.Role),
	}
	accessToken, err := uc.token.GenerateAccessToken(claims)
	if err != nil {
		return nil, apperrors.ErrInternalError(err)
	}

	// 生成 Refresh Token
	refreshToken, err := uc.token.GenerateRefreshToken(u.ID, int(u.Role))
	if err != nil {
		return nil, apperrors.ErrInternalError(err)
	}

	// 持久化 Refresh Token（7天有效）
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	tokenRecord := &domain.UserToken{
		UserID:       u.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		ExpiresAt:    expiresAt,
	}
	if err = uc.repo.SaveToken(ctx, tokenRecord); err != nil {
		uc.logger.Warn("save refresh token failed", zap.Error(err))
		// 不阻断登录，Token 已签发
	}

	return &TokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(2 * time.Hour), // access token 2小时
		User:         toUserResp(u),
	}, nil
}

// toUserResp 领域实体 → 响应 DTO（脱敏）
func toUserResp(u *domain.User) *UserResp {
	return &UserResp{
		ID:        u.ID,
		UUID:      u.UUID,
		Username:  u.Username,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Role:      int(u.Role),
		Status:    int(u.Status),
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
	}
}

// ─── fx Module ───────────────────────────────────────────────────────────────

// Module fx 用户应用模块
var Module = fx.Options(
	fx.Provide(
		NewUserUseCase,
		domain.NewDomainService,
	),
)
