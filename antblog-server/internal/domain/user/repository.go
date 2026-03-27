package user

import "context"

// IUserRepository 用户仓储接口（由 infrastructure 层实现）
type IUserRepository interface {
	// ── 用户 CRUD ────────────────────────────────────────────────────────────

	// Create 创建用户，返回含 ID 的完整实体
	Create(ctx context.Context, u *User) (*User, error)

	// FindByID 按 ID 查找用户（含软删除过滤）
	FindByID(ctx context.Context, id uint64) (*User, error)

	// FindByUUID 按 UUID 查找用户
	FindByUUID(ctx context.Context, uuid string) (*User, error)

	// FindByEmail 按邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindByUsername 按用户名查找用户
	FindByUsername(ctx context.Context, username string) (*User, error)

	// Update 全量更新用户（仅更新允许修改的字段）
	Update(ctx context.Context, u *User) error

	// UpdateLastLogin 更新最后登录时间
	UpdateLastLogin(ctx context.Context, id uint64) error

	// Delete 软删除用户
	Delete(ctx context.Context, id uint64) error

	// ExistsByEmail 检查邮箱是否已存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUsername 检查用户名是否已存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ── 令牌管理 ─────────────────────────────────────────────────────────────

	// SaveToken 保存 Refresh Token
	SaveToken(ctx context.Context, token *UserToken) error

	// FindToken 查找 Refresh Token
	FindToken(ctx context.Context, refreshToken string) (*UserToken, error)

	// DeleteToken 删除指定 Refresh Token（登出）
	DeleteToken(ctx context.Context, refreshToken string) error

	// DeleteUserTokens 删除用户所有 Token（强制下线）
	DeleteUserTokens(ctx context.Context, userID uint64) error

	// CleanExpiredTokens 清理过期 Token
	CleanExpiredTokens(ctx context.Context) error
}
