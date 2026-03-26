// Package cache 文章缓存层（Redis / 本地 Ristretto）。
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pkgcache "antblog/pkg/cache"
)

const (
	articleCacheTTL = 10 * time.Minute
	articleListTTL  = 5 * time.Minute
	archiveCacheTTL = 30 * time.Minute
	viewCountTTL    = 60 * time.Second // 阅读数写回节流
)

// IArticleCache 文章缓存接口
type IArticleCache interface {
	// SetArticle 缓存文章详情（slug → JSON）
	SetArticle(ctx context.Context, slug string, data any) error
	// GetArticle 读取文章详情缓存
	GetArticle(ctx context.Context, slug string) ([]byte, error)
	// DeleteArticle 删除文章详情缓存（内容更新时调用）
	DeleteArticle(ctx context.Context, slug string) error

	// SetArticleByID 缓存文章详情（id → JSON）
	SetArticleByID(ctx context.Context, id uint64, data any) error
	// GetArticleByID 读取文章详情缓存
	GetArticleByID(ctx context.Context, id uint64) ([]byte, error)
	// DeleteArticleByID 删除 ID 维度缓存
	DeleteArticleByID(ctx context.Context, id uint64) error

	// IncrViewCount 阅读数 Redis 自增（返回当前值，用于节流写回 DB）
	IncrViewCount(ctx context.Context, id uint64) (int64, error)
	// GetViewCount 读取阅读数缓存
	GetViewCount(ctx context.Context, id uint64) (int64, error)

	// InvalidateListCache 删除所有列表缓存（文章写操作后调用）
	// InvalidateListCache(ctx context.Context) error
}

// articleCache IArticleCache 实现
type articleCache struct {
	cache pkgcache.ICache
}

// NewArticleCache 创建文章缓存
func NewArticleCache(c pkgcache.ICache) IArticleCache {
	return &articleCache{cache: c}
}

func articleSlugKey(slug string) string { return fmt.Sprintf("article:slug:%s", slug) }
func articleIDKey(id uint64) string     { return fmt.Sprintf("article:id:%d", id) }
func articleViewKey(id uint64) string   { return fmt.Sprintf("article:view:%d", id) }

const articleListPrefix = "article:list:"

func (c *articleCache) SetArticle(ctx context.Context, slug string, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, articleSlugKey(slug), string(b), articleCacheTTL)
}

func (c *articleCache) GetArticle(ctx context.Context, slug string) ([]byte, error) {
	val, err := c.cache.Get(ctx, articleSlugKey(slug))
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

func (c *articleCache) DeleteArticle(ctx context.Context, slug string) error {
	return c.cache.Delete(ctx, articleSlugKey(slug))
}

func (c *articleCache) SetArticleByID(ctx context.Context, id uint64, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, articleIDKey(id), string(b), articleCacheTTL)
}

func (c *articleCache) GetArticleByID(ctx context.Context, id uint64) ([]byte, error) {
	val, err := c.cache.Get(ctx, articleIDKey(id))
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

func (c *articleCache) DeleteArticleByID(ctx context.Context, id uint64) error {
	return c.cache.Delete(ctx, articleIDKey(id))
}

func (c *articleCache) IncrViewCount(ctx context.Context, id uint64) (int64, error) {
	return c.cache.Incr(ctx, articleViewKey(id))
}

func (c *articleCache) GetViewCount(ctx context.Context, id uint64) (int64, error) {
	val, err := c.cache.Get(ctx, articleViewKey(id))
	if err != nil {
		return 0, err
	}
	var count int64
	fmt.Sscanf(string(val), "%d", &count)
	return count, nil
}

// func (c *articleCache) InvalidateListCache(ctx context.Context) error {
// 	// 使用 pattern delete（Redis SCAN + DEL）
// 	// 若使用本地缓存可直接忽略错误
// 	//return c.cache.DeleteByPrefix(ctx, articleListPrefix)
// }
