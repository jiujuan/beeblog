package logger

import (
	"os"
	"path/filepath"

	"beeblog/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func NewLogger(cfg config.LogConfig) (*zap.Logger, error) {
	// 确保日志目录存在
	logDir := filepath.Dir(cfg.Filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// 日志轮转配置
	hook := lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 文件输出
	fileWriter := zapcore.AddSync(&hook)
	// 控制台输出
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 核心
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileWriter, level),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleWriter, level),
	)

	l = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	Logger = l // // 为了兼容 pkg 内部的便捷方法，可以内部赋值一次（非强制，看团队习惯）

	return l, nil
}

var WireSet = NewLogger

// 以下是便捷方法，注意：如果严格遵循 DI，业务代码应该通过依赖注入获取 logger 实例
// 但为了方便中间件等地方使用，保留便捷方法也是常见的做法。
// 如果使用便捷方法，需要确保 NewLogger 被调用后，全局变量被赋值（如上注释部分）。

// 提供便捷方法
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
