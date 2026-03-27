// Package user 用户领域层 —— 纯业务实体，零外部依赖。
package user

import "time"

// Role 用户角色
type Role int

const (
	RoleUser  Role = 1 // 普通用户
	RoleAdmin Role = 2 // 管理员
)

func (r Role) IsAdmin() bool { return r == RoleAdmin }
func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	default:
		return "user"
	}
}

// Status 用户状态
type Status int

const (
	StatusActive   Status = 1 // 正常
	StatusDisabled Status = 2 // 禁用
)

func (s Status) IsActive() bool { return s == StatusActive }

// User 用户聚合根实体
type User struct {
	ID        uint64
	UUID      string
	Username  string
	Email     string
	Password  string // bcrypt hash，不对外暴露明文
	Nickname  string
	Avatar    string
	Bio       string
	Role      Role
	Status    Status
	LastLogin *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsAdmin 是否为管理员
func (u *User) IsAdmin() bool { return u.Role.IsAdmin() }

// IsActive 账号是否正常可用
func (u *User) IsActive() bool { return u.Status.IsActive() }

// Disable 禁用账号
func (u *User) Disable() { u.Status = StatusDisabled }

// Enable 启用账号
func (u *User) Enable() { u.Status = StatusActive }

// UpdateProfile 更新个人资料（仅允许修改安全字段）
func (u *User) UpdateProfile(nickname, avatar, bio string) {
	if nickname != "" {
		u.Nickname = nickname
	}
	if avatar != "" {
		u.Avatar = avatar
	}
	u.Bio = bio
}

// RecordLogin 记录登录时间
func (u *User) RecordLogin(t time.Time) {
	u.LastLogin = &t
}

// UserToken 用户刷新令牌实体
type UserToken struct {
	ID           uint64
	UserID       uint64
	RefreshToken string
	UserAgent    string
	ClientIP     string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

// IsExpired 判断 token 是否过期
func (t *UserToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
