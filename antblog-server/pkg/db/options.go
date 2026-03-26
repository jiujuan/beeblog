package db

import "time"

// Option db 构建选项（函数选项模式）
type Option func(*options)

type options struct {
	driver          string        // mysql | postgres
	dsn             string
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
	logLevel        string        // silent | error | warn | info
	slowThreshold   time.Duration // 慢查询阈值
	colorfulLog     bool
}

func defaultOptions() *options {
	return &options{
		driver:          "mysql",
		maxOpenConns:    100,
		maxIdleConns:    10,
		connMaxLifetime: time.Hour,
		connMaxIdleTime: 30 * time.Minute,
		logLevel:        "warn",
		slowThreshold:   200 * time.Millisecond,
		colorfulLog:     false,
	}
}

// WithDriver 设置数据库驱动（mysql | postgres）
func WithDriver(driver string) Option {
	return func(o *options) { o.driver = driver }
}

// WithDSN 设置数据库连接串
func WithDSN(dsn string) Option {
	return func(o *options) { o.dsn = dsn }
}

// WithMaxOpenConns 设置最大打开连接数
func WithMaxOpenConns(n int) Option {
	return func(o *options) { o.maxOpenConns = n }
}

// WithMaxIdleConns 设置最大空闲连接数
func WithMaxIdleConns(n int) Option {
	return func(o *options) { o.maxIdleConns = n }
}

// WithConnMaxLifetime 设置连接最大存活时间
func WithConnMaxLifetime(d time.Duration) Option {
	return func(o *options) { o.connMaxLifetime = d }
}

// WithConnMaxIdleTime 设置连接最大空闲时间
func WithConnMaxIdleTime(d time.Duration) Option {
	return func(o *options) { o.connMaxIdleTime = d }
}

// WithLogLevel 设置 GORM 日志级别（silent|error|warn|info）
func WithLogLevel(level string) Option {
	return func(o *options) { o.logLevel = level }
}

// WithSlowThreshold 设置慢查询阈值
func WithSlowThreshold(d time.Duration) Option {
	return func(o *options) { o.slowThreshold = d }
}

// WithColorfulLog 是否开启彩色日志（开发环境）
func WithColorfulLog(colorful bool) Option {
	return func(o *options) { o.colorfulLog = colorful }
}
