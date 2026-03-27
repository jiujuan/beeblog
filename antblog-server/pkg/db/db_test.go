package db

import (
	"testing"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

func TestNewDialectorUnsupported(t *testing.T) {
	t.Parallel()

	_, err := newDialector("sqlite", "test")
	if err == nil {
		t.Fatal("expected error for unsupported driver")
	}
}

func TestMustNewPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	_ = MustNew(WithDriver("unsupported"))
}

func TestNewGormLoggerLevel(t *testing.T) {
	t.Parallel()

	cases := []struct {
		logLevel string
		want     gormlogger.LogLevel
	}{
		{logLevel: "silent", want: gormlogger.Silent},
		{logLevel: "error", want: gormlogger.Error},
		{logLevel: "info", want: gormlogger.Info},
		{logLevel: "warn", want: gormlogger.Warn},
	}

	for _, c := range cases {
		o := defaultOptions()
		o.logLevel = c.logLevel
		l, ok := newGormLogger(o).(*gormZapLogger)
		if !ok {
			t.Fatal("unexpected logger type")
		}
		if l.level != c.want {
			t.Fatalf("log level mismatch for %s: got=%v want=%v", c.logLevel, l.level, c.want)
		}
	}
}

func TestOptionsSetters(t *testing.T) {
	t.Parallel()

	o := defaultOptions()
	WithDriver("postgres")(o)
	WithDSN("dsn")(o)
	WithMaxOpenConns(20)(o)
	WithMaxIdleConns(5)(o)
	WithLogLevel("info")(o)
	WithColorfulLog(true)(o)

	if o.driver != "postgres" || o.dsn != "dsn" || o.maxOpenConns != 20 || o.maxIdleConns != 5 || o.logLevel != "info" || !o.colorfulLog {
		t.Fatal("options setters did not apply correctly")
	}
}

func TestTraceNoPanic(t *testing.T) {
	t.Parallel()

	l := &gormZapLogger{level: gormlogger.Info}
	l.Trace(nil, nowTime(), func() (string, int64) { return "SELECT 1", 1 }, nil)
}

func nowTime() (t time.Time) {
	return
}

func TestLogMode(t *testing.T) {
	t.Parallel()

	l := &gormZapLogger{level: gormlogger.Warn}
	n := l.LogMode(gormlogger.Error)
	nl, ok := n.(*gormZapLogger)
	if !ok {
		t.Fatal("unexpected logger type")
	}
	if nl.level != gormlogger.Error {
		t.Fatalf("unexpected level: %v", nl.level)
	}
}
