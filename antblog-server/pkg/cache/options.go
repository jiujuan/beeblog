package cache

import "time"

// Option cache 构建选项（函数选项模式）
type Option func(*options)

type options struct {
	// Redis 选项
	redisAddr         string
	redisPassword     string
	redisDB           int
	redisPoolSize     int
	redisMinIdleConns int
	redisDialTimeout  time.Duration
	redisReadTimeout  time.Duration
	redisWriteTimeout time.Duration

	// Ristretto（本地缓存）选项
	ristrettoNumCounters int64
	ristrettoMaxCost     int64
	ristrettoBufferItems int64

	// 通用
	defaultTTL time.Duration
	keyPrefix  string
}

func defaultOptions() *options {
	return &options{
		redisAddr:         "127.0.0.1:6379",
		redisDB:           0,
		redisPoolSize:     10,
		redisMinIdleConns: 2,
		redisDialTimeout:  5 * time.Second,
		redisReadTimeout:  3 * time.Second,
		redisWriteTimeout: 3 * time.Second,

		ristrettoNumCounters: 1e7,
		ristrettoMaxCost:     1 << 30, // 1GB
		ristrettoBufferItems: 64,

		defaultTTL: 10 * time.Minute,
		keyPrefix:  "antblog:",
	}
}

// ── Redis 选项 ────────────────────────────────────────────────────────────────

func WithRedisAddr(addr string) Option {
	return func(o *options) { o.redisAddr = addr }
}

func WithRedisPassword(password string) Option {
	return func(o *options) { o.redisPassword = password }
}

func WithRedisDB(db int) Option {
	return func(o *options) { o.redisDB = db }
}

func WithRedisPoolSize(size int) Option {
	return func(o *options) { o.redisPoolSize = size }
}

func WithRedisMinIdleConns(n int) Option {
	return func(o *options) { o.redisMinIdleConns = n }
}

func WithRedisDialTimeout(d time.Duration) Option {
	return func(o *options) { o.redisDialTimeout = d }
}

// ── Ristretto 选项 ───────────────────────────────────────────────────────────

// WithRistrettoMaxCost 设置本地缓存最大容量（字节）
func WithRistrettoMaxCost(cost int64) Option {
	return func(o *options) { o.ristrettoMaxCost = cost }
}

// WithRistrettoNumCounters 设置计数器数量（建议为最大 item 数的 10 倍）
func WithRistrettoNumCounters(n int64) Option {
	return func(o *options) { o.ristrettoNumCounters = n }
}

// ── 通用选项 ─────────────────────────────────────────────────────────────────

// WithDefaultTTL 设置缓存默认过期时间
func WithDefaultTTL(d time.Duration) Option {
	return func(o *options) { o.defaultTTL = d }
}

// WithKeyPrefix 设置缓存 key 前缀
func WithKeyPrefix(prefix string) Option {
	return func(o *options) { o.keyPrefix = prefix }
}
