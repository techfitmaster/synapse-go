package tenant

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

func TestInjectAndFromContext(t *testing.T) {
	ctx := InjectToContext(context.Background(), "tenant-123")
	got := FromContext(ctx)
	if got != "tenant-123" {
		t.Errorf("FromContext() = %q, want %q", got, "tenant-123")
	}
}

func TestFromContext_Empty(t *testing.T) {
	got := FromContext(context.Background())
	if got != "" {
		t.Errorf("FromContext() = %q, want empty", got)
	}
}

func TestMiddleware_InjectsContext(t *testing.T) {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", "t-456")
		c.Next()
	})
	r.Use(Middleware("tenant_id"))
	r.GET("/test", func(c *gin.Context) {
		tid := FromContext(c.Request.Context())
		if tid != "t-456" {
			t.Errorf("tenant from context = %q, want %q", tid, "t-456")
		}
		c.JSON(http.StatusOK, gin.H{"tenant": tid})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestMiddleware_MissingTenant(t *testing.T) {
	r := gin.New()
	r.Use(Middleware("tenant_id")) // No upstream middleware setting tenant_id
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
