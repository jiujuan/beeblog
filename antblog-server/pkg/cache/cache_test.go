package cache

import (
	"context"
	"testing"
	"time"
)

func TestBuildKeyAndIsMiss(t *testing.T) {
	t.Parallel()

	if got := BuildKey("ant:", "k1"); got != "ant:k1" {
		t.Fatalf("unexpected key: %s", got)
	}
	if got := BuildKey("", "k1"); got != "k1" {
		t.Fatalf("unexpected key without prefix: %s", got)
	}
	if !IsMiss(&ErrCacheMiss{Key: "x"}) {
		t.Fatal("expected miss")
	}
	if IsMiss(nil) {
		t.Fatal("nil should not be miss")
	}
}

func TestOptionsDefaultsAndSetters(t *testing.T) {
	t.Parallel()

	o := defaultOptions()
	if o.redisAddr == "" || o.defaultTTL <= 0 || o.keyPrefix == "" {
		t.Fatal("unexpected defaults")
	}
	WithRedisAddr("127.0.0.1:6380")(o)
	WithRedisPassword("pwd")(o)
	WithRedisDB(2)(o)
	WithRistrettoMaxCost(1024)(o)
	WithRistrettoNumCounters(100)(o)
	WithDefaultTTL(time.Second)(o)
	WithKeyPrefix("x:")(o)
	if o.redisAddr != "127.0.0.1:6380" || o.redisPassword != "pwd" || o.redisDB != 2 || o.ristrettoMaxCost != 1024 || o.ristrettoNumCounters != 100 || o.defaultTTL != time.Second || o.keyPrefix != "x:" {
		t.Fatal("setters not applied")
	}
}

func TestRistrettoCacheFlow(t *testing.T) {
	t.Parallel()

	c, err := NewRistretto(WithDefaultTTL(100 * time.Millisecond), WithKeyPrefix("t:"))
	if err != nil {
		t.Fatalf("new ristretto failed: %v", err)
	}
	defer func() { _ = c.Close() }()

	ctx := context.Background()
	if err = c.Set(ctx, "k1", map[string]any{"a": 1}, 0); err != nil {
		t.Fatalf("set failed: %v", err)
	}

	got, err := c.Get(ctx, "k1")
	if err != nil || len(got) == 0 {
		t.Fatalf("get failed: err=%v", err)
	}

	exists, err := c.Exists(ctx, "k1")
	if err != nil || !exists {
		t.Fatalf("exists failed: exists=%v err=%v", exists, err)
	}

	ok, err := c.SetNX(ctx, "k1", "x", time.Second)
	if err != nil || ok {
		t.Fatalf("setnx on existing key should be false, ok=%v err=%v", ok, err)
	}

	ok, err = c.SetNX(ctx, "k2", "x", time.Second)
	if err != nil || !ok {
		t.Fatalf("setnx on new key should be true, ok=%v err=%v", ok, err)
	}
}

func TestRistrettoExpireAndIncr(t *testing.T) {
	t.Parallel()

	c, err := NewRistretto(WithKeyPrefix("t:"))
	if err != nil {
		t.Fatalf("new ristretto failed: %v", err)
	}
	defer func() { _ = c.Close() }()

	ctx := context.Background()
	if err = c.Set(ctx, "k", "v", time.Second); err != nil {
		t.Fatalf("set failed: %v", err)
	}
	if err = c.Expire(ctx, "k", 20*time.Millisecond); err != nil {
		t.Fatalf("expire failed: %v", err)
	}
	time.Sleep(30 * time.Millisecond)
	if _, err = c.Get(ctx, "k"); !IsMiss(err) {
		t.Fatalf("expected cache miss after expire, got err=%v", err)
	}

	v, err := c.Incr(ctx, "counter")
	if err != nil || v != 1 {
		t.Fatalf("incr failed: v=%d err=%v", v, err)
	}
	v, err = c.IncrBy(ctx, "counter", 3)
	if err != nil || v != 4 {
		t.Fatalf("incrby failed: v=%d err=%v", v, err)
	}
	if err = c.Ping(ctx); err != nil {
		t.Fatalf("ping failed: %v", err)
	}
}
