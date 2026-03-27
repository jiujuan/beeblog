package user

import "context"

// IUserUseCase 用户用例接口
type IUserUseCase interface {
	// Register 用户注册
	Register(ctx context.Context, req *RegisterReq) (*UserResp, error)

	// Login 用户登录，返回 Token 对
	Login(ctx context.Context, req *LoginReq, userAgent, clientIP string) (*TokenResp, error)

	// Logout 用户登出（撤销 Refresh Token）
	Logout(ctx context.Context, refreshToken string) error

	// RefreshToken 刷新 Access Token
	RefreshToken(ctx context.Context, req *RefreshTokenReq, userAgent, clientIP string) (*TokenResp, error)

	// GetProfile 获取用户个人资料
	GetProfile(ctx context.Context, userID uint64) (*UserResp, error)

	// UpdateProfile 更新个人资料
	UpdateProfile(ctx context.Context, userID uint64, req *UpdateProfileReq) (*UserResp, error)

	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, userID uint64, req *ChangePasswordReq) error
}
