// Package db 封装 GORM，支持 MySQL 和 PostgreSQL，
// 集成 Zap 日志、连接池配置，通过函数选项模式灵活配置。
package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"antblog/pkg/config"
)

// ─── 构建函数 ────────────────────────────────────────────────────────────────

// New 根据选项构建 *gorm.DB
func New(opts ...Option) (*gorm.DB, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return build(o)
}

// MustNew 构建失败时 panic
func MustNew(opts ...Option) *gorm.DB {
	db, err := New(opts...)
	if err != nil {
		panic(err)
	}
	return db
}

// NewFromConfig 从 config.DatabaseConfig 构建
func NewFromConfig(cfg config.DatabaseConfig, zapLogger *zap.Logger, appOpts ...Option) (*gorm.DB, error) {
	baseOpts := []Option{
		WithDriver(cfg.Driver),
		WithDSN(cfg.DSN),
		WithMaxOpenConns(cfg.MaxOpenConns),
		WithMaxIdleConns(cfg.MaxIdleConns),
		WithConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second),
		WithLogLevel(cfg.LogLevel),
	}
	return New(append(baseOpts, appOpts...)...)
}

// ─── 内部构建 ────────────────────────────────────────────────────────────────

func build(o *options) (*gorm.DB, error) {
	dialector, err := newDialector(o.driver, o.dsn)
	if err != nil {
		return nil, err
	}

	gormCfg := &gorm.Config{
		Logger:                                   newGormLogger(o),
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: false,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, fmt.Errorf("db: open connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("db: get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(o.maxOpenConns)
	sqlDB.SetMaxIdleConns(o.maxIdleConns)
	sqlDB.SetConnMaxLifetime(o.connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(o.connMaxIdleTime)

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("db: ping database: %w", err)
	}

	return db, nil
}

func newDialector(driver, dsn string) (gorm.Dialector, error) {
	switch strings.ToLower(driver) {
	case "mysql":
		return mysql.Open(dsn), nil
	case "postgres", "postgresql":
		return postgres.Open(dsn), nil
	default:
		return nil, fmt.Errorf("db: unsupported driver: %s", driver)
	}
}

// ─── GORM 日志适配器 ─────────────────────────────────────────────────────────

type gormZapLogger struct {
	level         gormlogger.LogLevel
	slowThreshold time.Duration
}

func newGormLogger(o *options) gormlogger.Interface {
	level := gormlogger.Warn
	switch strings.ToLower(o.logLevel) {
	case "silent":
		level = gormlogger.Silent
	case "error":
		level = gormlogger.Error
	case "info":
		level = gormlogger.Info
	}
	return &gormZapLogger{level: level, slowThreshold: o.slowThreshold}
}

func (l *gormZapLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return &gormZapLogger{level: level, slowThreshold: l.slowThreshold}
}

func (l *gormZapLogger) Info(_ context.Context, msg string, args ...any) {
	if l.level >= gormlogger.Info {
		zap.L().Sugar().Infof(msg, args...)
	}
}

func (l *gormZapLogger) Warn(_ context.Context, msg string, args ...any) {
	if l.level >= gormlogger.Warn {
		zap.L().Sugar().Warnf(msg, args...)
	}
}

func (l *gormZapLogger) Error(_ context.Context, msg string, args ...any) {
	if l.level >= gormlogger.Error {
		zap.L().Sugar().Errorf(msg, args...)
	}
}

func (l *gormZapLogger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level <= gormlogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sqlStr, rows := fc()

	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sqlStr),
	}

	switch {
	case err != nil && l.level >= gormlogger.Error:
		zap.L().Error("gorm error", append(fields, zap.Error(err))...)
	case l.slowThreshold > 0 && elapsed > l.slowThreshold && l.level >= gormlogger.Warn:
		zap.L().Warn("gorm slow query", fields...)
	case l.level >= gormlogger.Info:
		zap.L().Debug("gorm query", fields...)
	}
}

// ─── 健康检查 ────────────────────────────────────────────────────────────────

// Ping 检测数据库连接是否正常
func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Close 关闭数据库连接池
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Stats 返回连接池统计信息
func Stats(db *gorm.DB) (sql.DBStats, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return sql.DBStats{}, err
	}
	return sqlDB.Stats(), nil
}

// ─── fx Module ───────────────────────────────────────────────────────────────

// Module fx 模块，依赖 *config.Config 和 *zap.Logger，提供 *gorm.DB
var Module = fx.Options(
	fx.Provide(func(cfg *config.Config, zapLogger *zap.Logger) (*gorm.DB, error) {
		return NewFromConfig(cfg.Database, zapLogger)
	}),
)
