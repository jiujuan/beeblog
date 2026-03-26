// Package logger 基于 Zap 封装统一日志库，支持控制台/JSON 格式、
// 文件轮转（lumberjack）、结构化字段，通过函数选项模式灵活配置。
package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"antblog/pkg/config"
)

// ─── 全局 Logger ─────────────────────────────────────────────────────────────

var global *zap.Logger = zap.NewNop()

// L 返回全局 SugaredLogger（简便调用，如 logger.L().Infow）
func L() *zap.SugaredLogger {
	return global.Sugar()
}

// Raw 返回全局 *zap.Logger（高性能调用）
func Raw() *zap.Logger {
	return global
}

// ─── 构建函数 ────────────────────────────────────────────────────────────────

// New 根据选项构建 *zap.Logger 并设置为全局
func New(opts ...Option) (*zap.Logger, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return build(o)
}

// MustNew 构建失败时 panic
func MustNew(opts ...Option) *zap.Logger {
	l, err := New(opts...)
	if err != nil {
		panic(err)
	}
	return l
}

// NewFromConfig 从 config.LogConfig 构建 Logger
func NewFromConfig(cfg config.LogConfig, opts ...Option) (*zap.Logger, error) {
	baseOpts := []Option{
		WithLevel(cfg.Level),
		WithFormat(cfg.Format),
		WithOutput(cfg.Output),
		WithMaxSize(cfg.MaxSize),
		WithMaxBackups(cfg.MaxBackups),
		WithMaxAge(cfg.MaxAge),
		WithCompress(cfg.Compress),
	}
	return New(append(baseOpts, opts...)...)
}

// ─── 内部构建逻辑 ────────────────────────────────────────────────────────────

func build(o *options) (*zap.Logger, error) {
	level := parseLevel(o.level)

	// 编码器配置
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if strings.ToLower(o.format) == "json" {
		encoder = zapcore.NewJSONEncoder(encCfg)
	} else {
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encCfg)
	}

	// 输出目标
	writer := buildWriter(o)

	// Core
	core := zapcore.NewCore(encoder, writer, level)

	// 构建选项
	zapOpts := []zap.Option{}
	if o.caller {
		zapOpts = append(zapOpts, zap.AddCaller(), zap.AddCallerSkip(0))
	}
	if o.stacktrace {
		zapOpts = append(zapOpts, zap.AddStacktrace(zap.ErrorLevel))
	}

	l := zap.New(core, zapOpts...)
	global = l
	return l, nil
}

func buildWriter(o *options) zapcore.WriteSyncer {
	if o.output == "" || o.output == "stdout" {
		return zapcore.AddSync(os.Stdout)
	}
	if o.output == "stderr" {
		return zapcore.AddSync(os.Stderr)
	}

	// 文件输出 + 轮转
	rotation := &lumberjack.Logger{
		Filename:   o.output,
		MaxSize:    o.maxSize,
		MaxBackups: o.maxBackups,
		MaxAge:     o.maxAge,
		Compress:   o.compress,
	}

	// 同时写 stdout 和文件
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(io.Writer(rotation)),
	)
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// ─── fx Module ───────────────────────────────────────────────────────────────

// Module fx 模块，依赖 *config.Config，提供 *zap.Logger
var Module = fx.Options(
	fx.Provide(func(cfg *config.Config) (*zap.Logger, error) {
		return NewFromConfig(cfg.Log)
	}),
)
