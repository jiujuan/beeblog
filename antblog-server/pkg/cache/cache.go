// Package cache 提供统一缓存接口，支持 Redis（分布式）和
// Ristretto（本地内存）两种实现，通过接口抽象屏蔽底层细节。
package cache

import (
	"context"
	"time"
)

// ICache 通用缓存接口
type ICache interface {
	// Set 设置缓存，ttl=0 使用默认过期时间
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	// Get 获取缓存，返回序列化后的字节切片
	Get(ctx context.Context, key string) ([]byte, error)
	// Delete 删除缓存
	Delete(ctx context.Context, key string) error
	// Exists 判断 key 是否存在
	Exists(ctx context.Context, key string) (bool, error)
	// SetNX 仅在 key 不存在时设置（原子操作，用于分布式锁）
	SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error)
	// Expire 重置过期时间
	Expire(ctx context.Context, key string, ttl time.Duration) error
	// Incr key 自增 1，返回自增后的值
	Incr(ctx context.Context, key string) (int64, error)
	// IncrBy key 自增 delta
	IncrBy(ctx context.Context, key string, delta int64) (int64, error)
	// Close 关闭缓存连接
	Close() error
	// Ping 检测缓存服务连通性
	Ping(ctx context.Context) error
}

// ErrCacheMiss 缓存未命中错误
type ErrCacheMiss struct {
	Key string
}

func (e *ErrCacheMiss) Error() string {
	return "cache: key not found: " + e.Key
}

// IsMiss 判断是否为缓存未命中错误
func IsMiss(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*ErrCacheMiss)
	return ok
}

// BuildKey 统一构建带前缀的缓存 key
func BuildKey(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + key
}

// ─── 常用 Key 模板 ───────────────────────────────────────────────────────────

const (
	KeyUserInfo        = "user:info:%d"         // %d = userID
	KeyUserToken       = "user:token:%d"        // %d = userID
	KeyArticleDetail   = "article:detail:%d"    // %d = articleID
	KeyArticleList     = "article:list:%s"      // %s = page_pageSize_categoryID
	KeyCategoryAll     = "category:all"
	KeyTagAll          = "tag:all"
	KeyArticleLike     = "article:like:%d:%d"   // articleID:userID
	KeyArticleBookmark = "article:bookmark:%d:%d"
	KeyCommentList     = "comment:list:%d"      // %d = articleID
	KeyRateLimit       = "rate:limit:%s:%s"     // IP:路由
)
