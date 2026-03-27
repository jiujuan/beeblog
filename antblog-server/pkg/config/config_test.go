package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSuccess(t *testing.T) {
	dir := t.TempDir()
	content := `
app:
  name: antblog
  env: dev
  version: 1.0.0
  debug: true
server:
  host: 0.0.0.0
  port: 8080
database:
  driver: mysql
  dsn: "root:pwd@tcp(127.0.0.1:3306)/antblog"
redis:
  addr: 127.0.0.1:6379
jwt:
  secret: test
log:
  level: info
upload:
  driver: local
`
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config failed: %v", err)
	}

	cfg, err := Load(
		WithConfigPath(dir),
		WithConfigName("config"),
		WithConfigType("yaml"),
		WithAutoEnv(false),
	)
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}
	if cfg.App.Name != "antblog" {
		t.Fatalf("unexpected app name: %s", cfg.App.Name)
	}
	if cfg.Server.Port != 8080 {
		t.Fatalf("unexpected port: %d", cfg.Server.Port)
	}
}

func TestLoadWithEnvOverride(t *testing.T) {
	dir := t.TempDir()
	content := `
app:
  name: antblog
server:
  host: 0.0.0.0
  port: 8080
database:
  driver: mysql
  dsn: "dsn"
redis:
  addr: 127.0.0.1:6379
jwt:
  secret: test
log:
  level: info
upload:
  driver: local
`
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config failed: %v", err)
	}

	t.Setenv("ANTBLOG_SERVER_PORT", "9090")

	cfg, err := Load(
		WithConfigPath(dir),
		WithConfigName("config"),
		WithConfigType("yaml"),
		WithAutoEnv(true),
	)
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}
	if cfg.Server.Port != 9090 {
		t.Fatalf("env override not applied, got: %d", cfg.Server.Port)
	}
}

func TestMustLoadPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	_ = MustLoad(
		WithConfigPath(filepath.Join(t.TempDir(), "not-exist")),
		WithConfigName("config"),
		WithConfigType("yaml"),
	)
}
