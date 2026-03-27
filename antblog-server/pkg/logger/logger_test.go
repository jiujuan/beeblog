package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestParseLevel(t *testing.T) {
	cases := []struct {
		in   string
		want zapcore.Level
	}{
		{in: "debug", want: zapcore.DebugLevel},
		{in: "warn", want: zapcore.WarnLevel},
		{in: "warning", want: zapcore.WarnLevel},
		{in: "error", want: zapcore.ErrorLevel},
		{in: "panic", want: zapcore.PanicLevel},
		{in: "unknown", want: zapcore.InfoLevel},
	}

	for _, c := range cases {
		if got := parseLevel(c.in); got != c.want {
			t.Fatalf("parseLevel(%s)=%v, want=%v", c.in, got, c.want)
		}
	}
}

func TestOptionsDefaultsAndSetters(t *testing.T) {
	o := defaultOptions()
	if o.level != "info" || o.format != "console" || o.output != "stdout" {
		t.Fatal("unexpected default options")
	}

	WithLevel("debug")(o)
	WithFormat("json")(o)
	WithOutput("stderr")(o)
	WithCaller(false)(o)
	WithStacktrace(false)(o)
	if o.level != "debug" || o.format != "json" || o.output != "stderr" || o.caller || o.stacktrace {
		t.Fatal("setters not applied")
	}
}

func TestNewSetsGlobalAndWritesFile(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "app.log")
	l, err := New(
		WithFormat("json"),
		WithOutput(logFile),
		WithCaller(false),
		WithStacktrace(false),
	)
	if err != nil {
		t.Fatalf("new logger failed: %v", err)
	}
	defer func() { _ = l.Sync() }()

	if Raw() != l {
		t.Fatal("global logger was not replaced")
	}

	l.Info("unit-test-log")
	_ = l.Sync()

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("read log file failed: %v", err)
	}
	if !strings.Contains(string(data), "unit-test-log") {
		t.Fatal("expected log message in file")
	}
}
