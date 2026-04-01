package healthcheck

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// CheckFunc performs a single health check and returns an error if unhealthy.
type CheckFunc func(ctx context.Context) error

// Status represents the aggregated health status.
type Status struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}

// Checker aggregates multiple health checks.
type Checker struct {
	mu     sync.RWMutex
	checks map[string]CheckFunc
}

// New creates a Checker with no registered checks.
func New() *Checker {
	return &Checker{checks: make(map[string]CheckFunc)}
}

// Register adds a named health check function.
func (c *Checker) Register(name string, fn CheckFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checks[name] = fn
}

// Handler returns a Gin handler that runs all registered checks and responds with the result.
// Returns HTTP 200 if all checks pass, HTTP 503 if any check fails.
func (c *Checker) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c.mu.RLock()
		checks := make(map[string]CheckFunc, len(c.checks))
		for k, v := range c.checks {
			checks[k] = v
		}
		c.mu.RUnlock()

		checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
		defer cancel()

		result := Status{
			Status: "ok",
			Checks: make(map[string]string, len(checks)),
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		for name, fn := range checks {
			wg.Add(1)
			go func(name string, fn CheckFunc) {
				defer wg.Done()
				if err := fn(checkCtx); err != nil {
					mu.Lock()
					result.Status = "degraded"
					result.Checks[name] = err.Error()
					mu.Unlock()
				} else {
					mu.Lock()
					result.Checks[name] = "ok"
					mu.Unlock()
				}
			}(name, fn)
		}
		wg.Wait()

		code := http.StatusOK
		if result.Status != "ok" {
			code = http.StatusServiceUnavailable
		}
		ctx.JSON(code, result)
	}
}

// DBCheck returns a CheckFunc that verifies database connectivity.
func DBCheck(db *gorm.DB) CheckFunc {
	return func(ctx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.PingContext(ctx)
	}
}

// RedisCheck returns a CheckFunc that verifies Redis connectivity.
func RedisCheck(rdb *redis.Client) CheckFunc {
	return func(ctx context.Context) error {
		return rdb.Ping(ctx).Err()
	}
}
