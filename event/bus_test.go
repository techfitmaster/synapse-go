package event

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPublish_Sync(t *testing.T) {
	bus := New()
	var received string

	bus.Subscribe("order.created", func(ctx context.Context, payload any) {
		received = payload.(string)
	})

	bus.Publish(context.Background(), "order.created", "order-123")

	if received != "order-123" {
		t.Errorf("received = %q, want %q", received, "order-123")
	}
}

func TestPublish_MultipleHandlers(t *testing.T) {
	bus := New()
	var count int

	bus.Subscribe("user.login", func(ctx context.Context, payload any) { count++ })
	bus.Subscribe("user.login", func(ctx context.Context, payload any) { count++ })

	bus.Publish(context.Background(), "user.login", nil)

	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}

func TestPublish_NoHandlers(t *testing.T) {
	bus := New()
	// Should not panic
	bus.Publish(context.Background(), "nonexistent", nil)
}

func TestPublishAsync(t *testing.T) {
	bus := New()
	var count int64

	bus.Subscribe("async.event", func(ctx context.Context, payload any) {
		atomic.AddInt64(&count, 1)
	})

	bus.PublishAsync(context.Background(), "async.event", nil)

	time.Sleep(100 * time.Millisecond)

	if atomic.LoadInt64(&count) != 1 {
		t.Errorf("count = %d, want 1", atomic.LoadInt64(&count))
	}
}

func TestPublishAsync_PanicRecovery(t *testing.T) {
	bus := New()
	var wg sync.WaitGroup
	wg.Add(1)

	bus.Subscribe("panic.event", func(ctx context.Context, payload any) {
		defer wg.Done()
		panic("test panic")
	})

	// Should not crash the process
	bus.PublishAsync(context.Background(), "panic.event", nil)

	wg.Wait() // Wait for goroutine to finish (panic is recovered)
}

func TestSubscribe_DifferentEvents(t *testing.T) {
	bus := New()
	var a, b string

	bus.Subscribe("eventA", func(ctx context.Context, payload any) { a = payload.(string) })
	bus.Subscribe("eventB", func(ctx context.Context, payload any) { b = payload.(string) })

	bus.Publish(context.Background(), "eventA", "alpha")

	if a != "alpha" {
		t.Errorf("a = %q, want %q", a, "alpha")
	}
	if b != "" {
		t.Errorf("b = %q, want empty (eventB not published)", b)
	}
}
