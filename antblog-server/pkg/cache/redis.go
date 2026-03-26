package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"antblog/pkg/config"
)

// ─── Redis 实现 ──────────────────────────────────────────────────────────────

// RedisCache Redis 缓存实现，满足 ICache 接口
type RedisCache struct {
	client     *redis.Client
	defaultTTL time.Duration
	keyPrefix  string
}

// NewRedis 创建 Redis 缓存实例
func NewRedis(opts ...Option) (*RedisCache, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	client := redis.NewClient(&redis.Options{
		Addr:         o.redisAddr,
		Password:     o.redisPassword,
		DB:           o.redisDB,
		PoolSize:     o.redisPoolSize,
		MinIdleConns: o.redisMinIdleConns,
		DialTimeout:  o.redisDialTimeout,
		ReadTimeout:  o.redisReadTimeout,
		WriteTimeout: o.redisWriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("cache: redis ping failed: %w", err)
	}

	return &RedisCache{
		client:     client,
		defaultTTL: o.defaultTTL,
		keyPrefix:  o.keyPrefix,
	}, nil
}

// NewRedisFromConfig 从配置构建 Redis 实例
func NewRedisFromConfig(cfg config.RedisConfig, opts ...Option) (*RedisCache, error) {
	baseOpts := []Option{
		WithRedisAddr(cfg.Addr),
		WithRedisPassword(cfg.Password),
		WithRedisDB(cfg.DB),
		WithRedisPoolSize(cfg.PoolSize),
		WithRedisMinIdleConns(cfg.MinIdleConns),
		WithRedisDialTimeout(time.Duration(cfg.DialTimeout) * time.Second),
	}
	return NewRedis(append(baseOpts, opts...)...)
}

func (c *RedisCache) k(key string) string {
	return c.keyPrefix + key
}

func (c *RedisCache) ttl(d time.Duration) time.Duration {
	if d == 0 {
		return c.defaultTTL
	}
	return d
}

// Set 设置缓存（自动 JSON 序列化）
func (c *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache: marshal value: %w", err)
	}
	return c.client.Set(ctx, c.k(key), data, c.ttl(ttl)).Err()
}

// Get 获取缓存，返回原始 JSON 字节
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.client.Get(ctx, c.k(key)).Bytes()
	if err == redis.Nil {
		return nil, &ErrCacheMiss{Key: key}
	}
	return data, err
}

// GetJSON 获取并反序列化到 dest
func (c *RedisCache) GetJSON(ctx context.Context, key string, dest any) error {
	data, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete 删除缓存
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.k(key)).Err()
}

// DeleteMulti 批量删除
func (c *RedisCache) DeleteMulti(ctx context.Context, keys ...string) error {
	prefixed := make([]string, len(keys))
	for i, k := range keys {
		prefixed[i] = c.k(k)
	}
	return c.client.Del(ctx, prefixed...).Err()
}

// Exists 判断 key 是否存在
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, c.k(key)).Result()
	return n > 0, err
}

// SetNX 仅在 key 不存在时设置（用于分布式锁）
func (c *RedisCache) SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("cache: marshal value: %w", err)
	}
	return c.client.SetNX(ctx, c.k(key), data, c.ttl(ttl)).Result()
}

// Expire 重置过期时间
func (c *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, c.k(key), ttl).Err()
}

// Incr key 自增 1
func (c *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, c.k(key)).Result()
}

// IncrBy key 自增 delta
func (c *RedisCache) IncrBy(ctx context.Context, key string, delta int64) (int64, error) {
	return c.client.IncrBy(ctx, c.k(key), delta).Result()
}

// Close 关闭 Redis 连接
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// Ping 检测 Redis 连通性
func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Client 返回底层 *redis.Client（特殊场景使用）
func (c *RedisCache) Client() *redis.Client {
	return c.client
}

// 批量删除带前缀的缓存数据，Scan 批量扫描删除
// Scan 操作谨慎使用
func (c *RedisCache) DeleteByPrefix(ctx context.Context, match string) error {
	iter := c.client.Scan(ctx, 0, match, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
		if len(keys) >= 100 { // 每次扫描 100 个 key 进行批量删除
			c.client.Del(ctx, keys...)
			keys = []string{} // 重置切片
		}
	}
	// 删除剩余的 key
	if len(keys) > 0 {
		c.client.Del(ctx, keys...)
	}
	return iter.Err()
}

// ─── fx Module ───────────────────────────────────────────────────────────────

// Module fx 模块，依赖 *config.Config，提供 *RedisCache 和 ICache
var Module = fx.Options(
	fx.Provide(func(cfg *config.Config) (*RedisCache, error) {
		return NewRedisFromConfig(cfg.Redis)
	}),
	fx.Provide(func(r *RedisCache) ICache {
		return r
	}),
)
