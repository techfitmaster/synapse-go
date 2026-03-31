package audit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func init() { gin.SetMode(gin.TestMode) }

// mockStore captures audit entries for testing.
type mockStore struct {
	mu      sync.Mutex
	entries []*Entry
}

func (m *mockStore) Save(_ context.Context, entry *Entry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, entry)
	return nil
}

func (m *mockStore) SaveInTx(_ *gorm.DB, entry *Entry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, entry)
	return nil
}

func (m *mockStore) getEntries() []*Entry {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make([]*Entry, len(m.entries))
	copy(cp, m.entries)
	return cp
}

func TestMiddleware_AuditsWriteOperations(t *testing.T) {
	store := &mockStore{}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", int64(42))
		c.Set("username", "alice")
		c.Next()
	})
	r.Use(Middleware(store))
	r.POST("/api/orders", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/orders", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	r.ServeHTTP(w, req)

	// Wait for async goroutine
	time.Sleep(100 * time.Millisecond)

	entries := store.getEntries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	e := entries[0]
	if e.UserID != 42 {
		t.Errorf("UserID = %d, want 42", e.UserID)
	}
	if e.Username != "alice" {
		t.Errorf("Username = %q, want %q", e.Username, "alice")
	}
	if e.Action != "POST" {
		t.Errorf("Action = %q, want %q", e.Action, "POST")
	}
	if e.Resource != "/api/orders" {
		t.Errorf("Resource = %q, want %q", e.Resource, "/api/orders")
	}
	if e.Detail != "status:200" {
		t.Errorf("Detail = %q, want %q", e.Detail, "status:200")
	}
}

func TestMiddleware_SkipsGET(t *testing.T) {
	store := &mockStore{}

	r := gin.New()
	r.Use(Middleware(store))
	r.GET("/api/orders", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/orders", nil)
	r.ServeHTTP(w, req)

	time.Sleep(100 * time.Millisecond)

	if len(store.getEntries()) != 0 {
		t.Error("GET should not be audited")
	}
}

func TestMiddleware_AuditsPUT(t *testing.T) {
	store := &mockStore{}

	r := gin.New()
	r.Use(Middleware(store))
	r.PUT("/api/orders/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/api/orders/1", nil)
	r.ServeHTTP(w, req)

	time.Sleep(100 * time.Millisecond)

	entries := store.getEntries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Action != "PUT" {
		t.Errorf("Action = %q, want PUT", entries[0].Action)
	}
}

func TestMiddleware_AuditsDELETE(t *testing.T) {
	store := &mockStore{}

	r := gin.New()
	r.Use(Middleware(store))
	r.DELETE("/api/orders/:id", func(c *gin.Context) {
		c.JSON(http.StatusNoContent, nil)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/api/orders/1", nil)
	r.ServeHTTP(w, req)

	time.Sleep(100 * time.Millisecond)

	entries := store.getEntries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Detail != "status:204" {
		t.Errorf("Detail = %q, want %q", entries[0].Detail, "status:204")
	}
}

func TestGORMStore_SaveInTx(t *testing.T) {
	// Test that SaveInTx accepts a *gorm.DB and returns no error with mock
	store := &mockStore{}
	entry := &Entry{
		UserID:   1,
		Username: "test",
		Action:   "CREATE",
		Resource: "order",
	}
	if err := store.SaveInTx(nil, entry); err != nil {
		t.Errorf("SaveInTx() error: %v", err)
	}
	if len(store.getEntries()) != 1 {
		t.Error("entry not saved")
	}
}
