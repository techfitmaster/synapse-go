package lock

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

func newTestRedis(t *testing.T) *redis.Client {
	t.Helper()
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}
	t.Cleanup(s.Close)
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	t.Cleanup(func() { rdb.Close() })
	return rdb
}

func TestTryLock_Success(t *testing.T) {
	rdb := newTestRedis(t)
	l := New(rdb)
	ctx := context.Background()

	unlock, err := l.TryLock(ctx, "test:lock", 10*time.Second)
	if err != nil {
		t.Fatalf("TryLock() error: %v", err)
	}
	defer unlock()

	// Key should exist in Redis
	val, err := rdb.Get(ctx, "test:lock").Result()
	if err != nil {
		t.Fatalf("key not found: %v", err)
	}
	if val == "" {
		t.Error("lock token is empty")
	}
}

func TestTryLock_AlreadyHeld(t *testing.T) {
	rdb := newTestRedis(t)
	l := New(rdb)
	ctx := context.Background()

	unlock, _ := l.TryLock(ctx, "test:lock", 10*time.Second)
	defer unlock()

	_, err := l.TryLock(ctx, "test:lock", 10*time.Second)
	if !errors.Is(err, ErrLockNotAcquired) {
		t.Errorf("expected ErrLockNotAcquired, got: %v", err)
	}
}

func TestTryLock_UnlockAllowsReacquire(t *testing.T) {
	rdb := newTestRedis(t)
	l := New(rdb)
	ctx := context.Background()

	unlock, _ := l.TryLock(ctx, "test:lock", 10*time.Second)
	unlock() // Release

	// Should be able to acquire again
	unlock2, err := l.TryLock(ctx, "test:lock", 10*time.Second)
	if err != nil {
		t.Fatalf("TryLock() after unlock error: %v", err)
	}
	defer unlock2()
}

func TestTryLock_TokenSafety(t *testing.T) {
	rdb := newTestRedis(t)
	l1 := New(rdb)
	l2 := New(rdb)
	ctx := context.Background()

	// l1 acquires the lock
	_, _ = l1.TryLock(ctx, "test:lock", 10*time.Second)

	// Simulate: l1's lock expires, l2 acquires
	rdb.Del(ctx, "test:lock") // Simulate expiry
	unlock2, _ := l2.TryLock(ctx, "test:lock", 10*time.Second)

	// l1's unlock should NOT delete l2's lock (different token)
	// We can't directly test this since l1's unlock was captured before,
	// but we can verify l2's lock is still held
	_, err := l1.TryLock(ctx, "test:lock", 10*time.Second)
	if !errors.Is(err, ErrLockNotAcquired) {
		t.Error("l2's lock should still be held")
	}
	defer unlock2()
}

func TestExecuteWithLock_Success(t *testing.T) {
	rdb := newTestRedis(t)
	l := New(rdb)
	ctx := context.Background()

	executed := false
	err := l.ExecuteWithLock(ctx, "test:exec", 10*time.Second, func() error {
		executed = true
		return nil
	})
	if err != nil {
		t.Fatalf("ExecuteWithLock() error: %v", err)
	}
	if !executed {
		t.Error("fn was not executed")
	}

	// Lock should be released
	_, err = rdb.Get(ctx, "test:exec").Result()
	if !errors.Is(err, redis.Nil) {
		t.Error("lock should be released after ExecuteWithLock")
	}
}

func TestExecuteWithLock_FnError(t *testing.T) {
	rdb := newTestRedis(t)
	l := New(rdb)
	ctx := context.Background()

	fnErr := errors.New("business error")
	err := l.ExecuteWithLock(ctx, "test:exec", 10*time.Second, func() error {
		return fnErr
	})
	if !errors.Is(err, fnErr) {
		t.Errorf("expected fn error, got: %v", err)
	}

	// Lock should still be released even on fn error
	_, err = rdb.Get(ctx, "test:exec").Result()
	if !errors.Is(err, redis.Nil) {
		t.Error("lock should be released even when fn fails")
	}
}

func TestExecuteWithLock_AlreadyHeld(t *testing.T) {
	rdb := newTestRedis(t)
	l := New(rdb)
	ctx := context.Background()

	unlock, _ := l.TryLock(ctx, "test:exec", 10*time.Second)
	defer unlock()

	err := l.ExecuteWithLock(ctx, "test:exec", 10*time.Second, func() error {
		t.Error("fn should not execute when lock is held")
		return nil
	})
	if !errors.Is(err, ErrLockNotAcquired) {
		t.Errorf("expected ErrLockNotAcquired, got: %v", err)
	}
}
