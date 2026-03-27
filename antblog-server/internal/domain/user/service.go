package user

import (
	"context"
	"time"

	apperrors "antblog/pkg/errors"
)

// IDomainService 用户领域服务接口
type IDomainService interface {
	// ValidateRegister 校验注册信息（值对象验证 + 唯一性检查）
	ValidateRegister(ctx context.Context, username, email string) error

	// ValidateCredentials 校验登录凭据，返回用户实体
	ValidateCredentials(ctx context.Context, loginID, password string, verifyFn func(plain, hash string) bool) (*User, error)

	// BuildNewUser 构建新用户实体（值对象验证，未持久化）
	BuildNewUser(username, email, hashedPwd string) (*User, error)
}

// DomainService 用户领域服务实现
type DomainService struct {
	repo IUserRepository
}

// NewDomainService 创建领域服务
func NewDomainService(repo IUserRepository) IDomainService {
	return &DomainService{repo: repo}
}

// ValidateRegister 校验注册信息：值对象合法性 + 唯一性
func (s *DomainService) ValidateRegister(ctx context.Context, username, email string) error {
	if _, err := NewUsername(username); err != nil {
		return apperrors.ErrInvalidParams(err.Error())
	}
	if _, err := NewEmail(email); err != nil {
		return apperrors.ErrInvalidParams(err.Error())
	}

	exists, err := s.repo.ExistsByUsername(ctx, username)
	if err != nil {
		return apperrors.ErrInternalError(err)
	}
	if exists {
		return apperrors.FromCode(apperrors.CodeUserAlreadyExists)
	}

	exists, err = s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return apperrors.ErrInternalError(err)
	}
	if exists {
		return apperrors.FromCode(apperrors.CodeEmailAlreadyExists)
	}
	return nil
}

// ValidateCredentials 验证登录凭据，loginID 可以是 username 或 email
func (s *DomainService) ValidateCredentials(ctx context.Context, loginID, password string, verifyFn func(plain, hash string) bool) (*User, error) {
	var (
		u   *User
		err error
	)

	// 判断 loginID 是 email 还是 username
	if _, emailErr := NewEmail(loginID); emailErr == nil {
		u, err = s.repo.FindByEmail(ctx, loginID)
	} else {
		u, err = s.repo.FindByUsername(ctx, loginID)
	}

	if err != nil {
		return nil, apperrors.ErrUserNotFound()
	}
	if !u.IsActive() {
		return nil, apperrors.ErrUserDisabled()
	}
	if !verifyFn(password, u.Password) {
		return nil, apperrors.ErrPasswordIncorrect()
	}
	return u, nil
}

// BuildNewUser 构建待创建的用户实体（不含 ID，由仓储层赋值）
func (s *DomainService) BuildNewUser(username, email, hashedPwd string) (*User, error) {
	if _, err := NewUsername(username); err != nil {
		return nil, apperrors.ErrInvalidParams(err.Error())
	}
	if _, err := NewEmail(email); err != nil {
		return nil, apperrors.ErrInvalidParams(err.Error())
	}

	now := time.Now()
	return &User{
		Username:  username,
		Email:     email,
		Password:  hashedPwd,
		Nickname:  username, // 默认昵称与用户名相同
		Role:      RoleUser,
		Status:    StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
