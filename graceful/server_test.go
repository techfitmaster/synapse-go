package graceful

import (
	"testing"
	"time"
)

func TestWithShutdownTimeout(t *testing.T) {
	cfg := &serverConfig{shutdownTimeout: 10 * time.Second}
	WithShutdownTimeout(30 * time.Second)(cfg)

	if cfg.shutdownTimeout != 30*time.Second {
		t.Errorf("shutdownTimeout = %v, want 30s", cfg.shutdownTimeout)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := &serverConfig{shutdownTimeout: 10 * time.Second}
	if cfg.shutdownTimeout != 10*time.Second {
		t.Errorf("default shutdownTimeout = %v, want 10s", cfg.shutdownTimeout)
	}
}
