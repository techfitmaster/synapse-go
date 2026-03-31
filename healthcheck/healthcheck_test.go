package healthcheck

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

func TestHandler_AllHealthy(t *testing.T) {
	c := New()
	c.Register("db", func(ctx context.Context) error { return nil })
	c.Register("redis", func(ctx context.Context) error { return nil })

	r := gin.New()
	r.GET("/health", c.Handler())

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var s Status
	json.Unmarshal(w.Body.Bytes(), &s)
	if s.Status != "ok" {
		t.Errorf("status = %q, want ok", s.Status)
	}
	if s.Checks["db"] != "ok" || s.Checks["redis"] != "ok" {
		t.Errorf("checks = %v", s.Checks)
	}
}

func TestHandler_Degraded(t *testing.T) {
	c := New()
	c.Register("db", func(ctx context.Context) error { return nil })
	c.Register("redis", func(ctx context.Context) error { return errors.New("connection refused") })

	r := gin.New()
	r.GET("/health", c.Handler())

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", w.Code, http.StatusServiceUnavailable)
	}

	var s Status
	json.Unmarshal(w.Body.Bytes(), &s)
	if s.Status != "degraded" {
		t.Errorf("status = %q, want degraded", s.Status)
	}
	if s.Checks["redis"] != "connection refused" {
		t.Errorf("redis check = %q", s.Checks["redis"])
	}
}

func TestHandler_NoChecks(t *testing.T) {
	c := New()

	r := gin.New()
	r.GET("/health", c.Handler())

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}
