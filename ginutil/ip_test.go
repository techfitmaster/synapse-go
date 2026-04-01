package ginutil

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestClientIP_XForwardedFor(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("X-Forwarded-For", "203.0.113.1, 10.0.0.1, 172.16.0.1")

	ip := ClientIP(c)
	if ip != "203.0.113.1" {
		t.Errorf("ClientIP() = %q, want %q", ip, "203.0.113.1")
	}
}

func TestClientIP_XRealIP(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("X-Real-IP", "198.51.100.1")

	ip := ClientIP(c)
	if ip != "198.51.100.1" {
		t.Errorf("ClientIP() = %q, want %q", ip, "198.51.100.1")
	}
}

func TestClientIP_FallbackRemoteAddr(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.RemoteAddr = "192.0.2.1:12345"

	ip := ClientIP(c)
	if ip != "192.0.2.1" {
		t.Errorf("ClientIP() = %q, want %q", ip, "192.0.2.1")
	}
}

func TestClientIP_SkipsUnknown(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("X-Forwarded-For", "unknown, 203.0.113.5")

	ip := ClientIP(c)
	if ip != "203.0.113.5" {
		t.Errorf("ClientIP() = %q, want %q", ip, "203.0.113.5")
	}
}

func TestClientIP_EmptyHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.RemoteAddr = "10.0.0.1:8080"

	ip := ClientIP(c)
	if ip != "10.0.0.1" {
		t.Errorf("ClientIP() = %q, want %q", ip, "10.0.0.1")
	}
}
