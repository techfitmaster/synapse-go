package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// init() is in auth_test.go (same package) — no need to duplicate here.

// ── RequireRole Tests ───────────────────────────────────────────

func TestRequireRole_Allowed(t *testing.T) {
	tests := []struct {
		name         string
		allowedRoles []string
		currentRole  string
	}{
		{"single_match", []string{"admin"}, "admin"},
		{"multi_first", []string{"admin", "editor"}, "admin"},
		{"multi_second", []string{"admin", "editor"}, "editor"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(func(c *gin.Context) {
				c.Set("role", tt.currentRole)
				c.Next()
			})
			r.Use(RequireRole(tt.allowedRoles...))
			r.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"ok": true})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
			}
		})
	}
}

func TestRequireRole_Denied(t *testing.T) {
	tests := []struct {
		name        string
		allowed     []string
		currentRole string
	}{
		{"wrong_role", []string{"admin"}, "viewer"},
		{"empty_role", []string{"admin"}, ""},
		{"unknown_role", []string{"admin", "editor"}, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(func(c *gin.Context) {
				if tt.currentRole != "" {
					c.Set("role", tt.currentRole)
				}
				c.Next()
			})
			r.Use(RequireRole(tt.allowed...))
			r.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"ok": true})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusForbidden {
				t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
			}
		})
	}
}

func TestRequireRole_NoRoleInContext(t *testing.T) {
	r := gin.New()
	r.Use(RequireRole("admin"))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

// ── RequireHeaderSecret Tests ───────────────────────────────────

func TestRequireHeaderSecret(t *testing.T) {
	const secret = "test-secret-123"

	tests := []struct {
		name       string
		header     string
		wantStatus int
	}{
		{"missing_header", "", http.StatusForbidden},
		{"wrong_secret", "wrong-secret", http.StatusForbidden},
		{"correct_secret", secret, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(RequireHeaderSecret("X-Admin-Secret", secret))
			r.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"ok": true})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.header != "" {
				req.Header.Set("X-Admin-Secret", tt.header)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
