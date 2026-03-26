// Package config 使用 Viper 提供配置加载能力，支持 YAML/JSON/TOML，
// 支持环境变量覆盖和热重载，通过函数选项模式灵活配置。
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// ─── 应用配置结构 ────────────────────────────────────────────────────────────

// Config 完整应用配置
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Upload   UploadConfig   `mapstructure:"upload"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Env     string `mapstructure:"env"`  // dev | test | prod
	Version string `mapstructure:"version"`
	Debug   bool   `mapstructure:"debug"`
}

type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`  // 秒
	WriteTimeout int    `mapstructure:"write_timeout"` // 秒
}

type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`   // mysql | postgres
	DSN             string `mapstructure:"dsn"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // 秒
	LogLevel        string `mapstructure:"log_level"`         // silent|error|warn|info
}

type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
	DialTimeout  int    `mapstructure:"dial_timeout"`  // 秒
	ReadTimeout  int    `mapstructure:"read_timeout"`  // 秒
	WriteTimeout int    `mapstructure:"write_timeout"` // 秒
}

type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	AccessExpiry    int    `mapstructure:"access_expiry"`  // 分钟
	RefreshExpiry   int    `mapstructure:"refresh_expiry"` // 天
	Issuer         string `mapstructure:"issuer"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`       // debug|info|warn|error
	Format     string `mapstructure:"format"`      // console|json
	Output     string `mapstructure:"output"`      // stdout | /path/to/file
	MaxSize    int    `mapstructure:"max_size"`    // MB
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`     // 天
	Compress   bool   `mapstructure:"compress"`
}

type UploadConfig struct {
	Driver    string `mapstructure:"driver"`     // local | oss
	LocalPath string `mapstructure:"local_path"` // 本地存储路径
	BaseURL   string `mapstructure:"base_url"`   // 访问域名
	MaxSize   int64  `mapstructure:"max_size"`   // 字节
	AllowedTypes []string `mapstructure:"allowed_types"`
}

// ─── 加载函数 ────────────────────────────────────────────────────────────────

// Load 加载配置，返回 *Config
func Load(opts ...Option) (*Config, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	v := viper.New()
	v.AddConfigPath(o.configPath)
	v.SetConfigName(o.configName)
	v.SetConfigType(o.configType)

	if o.autoEnv {
		v.SetEnvPrefix(o.envPrefix)
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.AutomaticEnv()
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config: read config file: %w", err)
	}

	if o.watchConfig {
		v.WatchConfig()
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal config: %w", err)
	}

	return cfg, nil
}

// MustLoad 加载配置，失败时 panic
func MustLoad(opts ...Option) *Config {
	cfg, err := Load(opts...)
	if err != nil {
		panic(err)
	}
	return cfg
}

// ─── fx Module ───────────────────────────────────────────────────────────────

// Params fx 注入参数
type Params struct {
	fx.In
	Opts []Option `optional:"true"`
}

// NewConfig fx provider
func NewConfig(opts ...Option) (*Config, error) {
	return Load(opts...)
}

// Module fx 模块，提供 *Config
var Module = fx.Options(
	fx.Provide(func() (*Config, error) {
		return Load(
			WithConfigPath("./config"),
			WithConfigName("config"),
			WithConfigType("yaml"),
		)
	}),
)
