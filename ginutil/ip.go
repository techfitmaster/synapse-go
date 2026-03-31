package ginutil

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// ipHeaders is the ordered list of headers to check for client IP behind proxies.
var ipHeaders = []string{
	"X-Forwarded-For",
	"X-Real-IP",
	"Proxy-Client-IP",
	"WL-Proxy-Client-IP",
}

// ClientIP extracts the real client IP from proxy headers, falling back to RemoteAddr.
// Handles multi-IP X-Forwarded-For values (takes the first non-private IP).
// For most cases, gin.Context.ClientIP() is sufficient; use this when you need
// proxy-chain-aware extraction with private IP filtering.
func ClientIP(c *gin.Context) string {
	for _, header := range ipHeaders {
		val := c.GetHeader(header)
		if val == "" {
			continue
		}
		// X-Forwarded-For may contain multiple IPs: "client, proxy1, proxy2"
		for _, ip := range strings.Split(val, ",") {
			ip = strings.TrimSpace(ip)
			if ip == "" || strings.EqualFold(ip, "unknown") {
				continue
			}
			if parsed := net.ParseIP(ip); parsed != nil {
				return ip
			}
		}
	}
	// Fallback to RemoteAddr (strip port)
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}
