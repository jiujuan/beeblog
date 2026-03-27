// Package user 用户应用层 —— 用例编排，依赖领域接口，不依赖基础设施。
package user

import "time"

// ─── 请求 DTO ────────────────────────────────────────────────────────────────

// RegisterReq 注册请求
type RegisterReq struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

// LoginReq 登录请求
type LoginReq struct {
	Email    string `json:"email"  validate:"required"` // 用户名 or 邮箱
	Password string `json:"password"  validate:"required"`
}

// RefreshTokenReq 刷新 Token 请求
type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UpdateProfileReq 更新个人资料请求
type UpdateProfileReq struct {
	Nickname string `json:"nickname" validate:"max=64"`
	Avatar   string `json:"avatar"   validate:"max=512,omitempty,url"`
	Bio      string `json:"bio"      validate:"max=512"`
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=72"`
}

// ─── 响应 DTO ────────────────────────────────────────────────────────────────

// UserResp 用户信息响应（脱敏，不含密码）
type UserResp struct {
	ID        uint64     `json:"id"`
	UUID      string     `json:"uuid"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Nickname  string     `json:"nickname"`
	Avatar    string     `json:"avatar"`
	Bio       string     `json:"bio"`
	Role      int        `json:"role"`
	Status    int        `json:"status"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// TokenResp 登录/刷新 Token 响应
type TokenResp struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"` // access token 过期时间
	User         *UserResp `json:"user"`
}
