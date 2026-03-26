// Package cache 基础设施缓存层，为用户模块提供 Redis 缓存加速。
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	domain "antblog/internal/domain/user"
	"antblog/pkg/cache"
	apperrors "antblog/pkg/errors"
)

const (
	userCacheTTL      = 30 * time.Minute
	userTokenBlockTTL = 24 * time.Hour // 登出后 access token 黑名单时长
)

// IUserCache 用户缓存接口
type IUserCache interface {
	// SetUser 缓存用户信息
	SetUser(ctx context.Context, u *domain.User) error
	// GetUser 获取缓存的用户信息
	GetUser(ctx context.Context, userID uint64) (*domain.User, error)
	// DeleteUser 删除用户缓存
	DeleteUser(ctx context.Context, userID uint64) error
	// BlockToken 将 Access Token 加入黑名单（登出/改密用）
	BlockToken(ctx context.Context, token string, ttl time.Duration) error
	// IsTokenBlocked 检查 Access Token 是否在黑名单中
	IsTokenBlocked(ctx context.Context, token string) (bool, error)
}

// userCache IUserCache Redis 实现
type userCache struct {
	cache  cache.ICache
	logger *zap.Logger
}

// NewUserCache 创建用户缓存实例
func NewUserCache(c cache.ICache, logger *zap.Logger) IUserCache {
	return &userCache{cache: c, logger: logger}
}

func userKey(userID uint64) string {
	return fmt.Sprintf(cache.KeyUserInfo, userID)
}

func tokenBlacklistKey(token string) string {
	// 使用 token 的 SHA256 摘要作为 key，避免 key 过长
	h := fmt.Sprintf("token:blacklist:%x", token[:min(32, len(token))])
	return h
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// SetUser 缓存用户信息（JSON 序列化）
func (c *userCache) SetUser(ctx context.Context, u *domain.User) error {
	// 缓存的用户信息不含密码
	data := &cachedUser{
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
	return c.cache.Set(ctx, userKey(u.ID), data, userCacheTTL)
}

// GetUser 获取缓存用户信息
func (c *userCache) GetUser(ctx context.Context, userID uint64) (*domain.User, error) {
	raw, err := c.cache.Get(ctx, userKey(userID))
	if err != nil {
		if cache.IsMiss(err) {
			return nil, apperrors.ErrUserNotFound()
		}
		return nil, err
	}

	var data cachedUser
	if err := json.Unmarshal(raw, &data); err != nil {
		c.logger.Warn("unmarshal user cache failed", zap.Error(err))
		return nil, apperrors.ErrUserNotFound()
	}

	return &domain.User{
		ID:        data.ID,
		UUID:      data.UUID,
		Username:  data.Username,
		Email:     data.Email,
		Nickname:  data.Nickname,
		Avatar:    data.Avatar,
		Bio:       data.Bio,
		Role:      domain.Role(data.Role),
		Status:    domain.Status(data.Status),
		LastLogin: data.LastLogin,
		CreatedAt: data.CreatedAt,
	}, nil
}

// DeleteUser 删除用户缓存
func (c *userCache) DeleteUser(ctx context.Context, userID uint64) error {
	return c.cache.Delete(ctx, userKey(userID))
}

// BlockToken 将 Access Token 加入黑名单
func (c *userCache) BlockToken(ctx context.Context, token string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = userTokenBlockTTL
	}
	return c.cache.Set(ctx, tokenBlacklistKey(token), "1", ttl)
}

// IsTokenBlocked 检查 Access Token 是否被加入黑名单
func (c *userCache) IsTokenBlocked(ctx context.Context, token string) (bool, error) {
	return c.cache.Exists(ctx, tokenBlacklistKey(token))
}

// cachedUser 缓存中的用户数据结构（不含密码）
type cachedUser struct {
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
