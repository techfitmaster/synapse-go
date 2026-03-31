package circuitbreaker

import (
	"testing"
	"time"
)

func defaultConfig() Config {
	return Config{
		FailureThreshold:    3,
		OpenDuration:        100 * time.Millisecond,
		HalfOpenMaxRequests: 2,
	}
}

func TestClosed_AllowsRequests(t *testing.T) {
	cb := New(defaultConfig())
	if !cb.IsAvailable("svc") {
		t.Error("closed circuit should be available")
	}
}

func TestClosed_OpensAfterThreshold(t *testing.T) {
	cb := New(defaultConfig())

	for i := 0; i < 3; i++ {
		cb.RecordFailure("svc")
	}

	if cb.IsAvailable("svc") {
		t.Error("circuit should be open after 3 failures")
	}
}

func TestOpen_BlocksRequests(t *testing.T) {
	cb := New(defaultConfig())

	for i := 0; i < 3; i++ {
		cb.RecordFailure("svc")
	}

	if cb.IsAvailable("svc") {
		t.Error("open circuit should block")
	}
}

func TestOpen_TransitionsToHalfOpen(t *testing.T) {
	cb := New(defaultConfig())

	for i := 0; i < 3; i++ {
		cb.RecordFailure("svc")
	}

	// Wait for open duration to expire
	time.Sleep(150 * time.Millisecond)

	if !cb.IsAvailable("svc") {
		t.Error("circuit should transition to half-open after timeout")
	}
}

func TestHalfOpen_ClosesAfterSuccesses(t *testing.T) {
	cb := New(defaultConfig())

	for i := 0; i < 3; i++ {
		cb.RecordFailure("svc")
	}

	time.Sleep(150 * time.Millisecond)
	cb.IsAvailable("svc") // Trigger half-open

	cb.RecordSuccess("svc")
	cb.RecordSuccess("svc") // HalfOpenMaxRequests = 2

	// Should be closed now
	cb.RecordFailure("svc") // 1 failure, below threshold
	if !cb.IsAvailable("svc") {
		t.Error("circuit should be closed after successful probes")
	}
}

func TestHalfOpen_ReopensOnFailure(t *testing.T) {
	cb := New(defaultConfig())

	for i := 0; i < 3; i++ {
		cb.RecordFailure("svc")
	}

	time.Sleep(150 * time.Millisecond)
	cb.IsAvailable("svc") // Trigger half-open

	cb.RecordFailure("svc") // Probe failed

	if cb.IsAvailable("svc") {
		t.Error("circuit should reopen on half-open failure")
	}
}

func TestSuccess_ResetsFailureCount(t *testing.T) {
	cb := New(defaultConfig())

	cb.RecordFailure("svc")
	cb.RecordFailure("svc")
	cb.RecordSuccess("svc") // Resets counter
	cb.RecordFailure("svc") // Only 1 failure now

	if !cb.IsAvailable("svc") {
		t.Error("circuit should still be closed after reset")
	}
}

func TestIndependentKeys(t *testing.T) {
	cb := New(defaultConfig())

	for i := 0; i < 3; i++ {
		cb.RecordFailure("svcA")
	}

	if cb.IsAvailable("svcA") {
		t.Error("svcA should be open")
	}
	if !cb.IsAvailable("svcB") {
		t.Error("svcB should be unaffected")
	}
}
