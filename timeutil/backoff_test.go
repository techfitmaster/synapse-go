package timeutil

import (
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	tests := []struct {
		attempt  int
		initial  time.Duration
		max      time.Duration
		expected time.Duration
	}{
		{0, 1 * time.Second, 60 * time.Second, 1 * time.Second},
		{1, 1 * time.Second, 60 * time.Second, 2 * time.Second},
		{2, 1 * time.Second, 60 * time.Second, 4 * time.Second},
		{3, 1 * time.Second, 60 * time.Second, 8 * time.Second},
		{10, 1 * time.Second, 60 * time.Second, 60 * time.Second}, // capped
		{0, 10 * time.Minute, 3 * time.Hour, 10 * time.Minute},
		{1, 10 * time.Minute, 3 * time.Hour, 20 * time.Minute},
	}
	for _, tt := range tests {
		d := Backoff(tt.attempt, tt.initial, tt.max)
		if d != tt.expected {
			t.Errorf("Backoff(%d, %v, %v) = %v, want %v", tt.attempt, tt.initial, tt.max, d, tt.expected)
		}
	}
}

func TestShouldRetry(t *testing.T) {
	if !ShouldRetry(0, 3) {
		t.Error("attempt 0 should retry")
	}
	if !ShouldRetry(2, 3) {
		t.Error("attempt 2 should retry (max=3)")
	}
	if ShouldRetry(3, 3) {
		t.Error("attempt 3 should NOT retry (max=3)")
	}
}
