package logger

import (
	"sync"
	"sync/atomic"

	"go.uber.org/zap"
)

var (
	global   atomic.Value
	initOnce sync.Once
)

// New creates a zap logger. Uses production config for "production" env, development otherwise.
func New(env string) *zap.Logger {
	var log *zap.Logger
	if env == "production" {
		log, _ = zap.NewProduction()
	} else {
		log, _ = zap.NewDevelopment()
	}
	global.Store(log)
	return log
}

func get() *zap.Logger {
	if v := global.Load(); v != nil {
		return v.(*zap.Logger)
	}
	initOnce.Do(func() {
		l, _ := zap.NewDevelopment()
		global.Store(l)
	})
	return global.Load().(*zap.Logger)
}

// Info logs a message at info level.
func Info(msg string, fields ...zap.Field) { get().Info(msg, fields...) }

// Error logs a message at error level.
func Error(msg string, fields ...zap.Field) { get().Error(msg, fields...) }

// Warn logs a message at warn level.
func Warn(msg string, fields ...zap.Field) { get().Warn(msg, fields...) }

// Debug logs a message at debug level.
func Debug(msg string, fields ...zap.Field) { get().Debug(msg, fields...) }

// String creates a string-typed zap field.
func String(key, val string) zap.Field { return zap.String(key, val) }

// Int creates an int-typed zap field.
func Int(key string, val int) zap.Field { return zap.Int(key, val) }

// Int64 creates an int64-typed zap field.
func Int64(key string, val int64) zap.Field { return zap.Int64(key, val) }

// Bool creates a bool-typed zap field.
func Bool(key string, val bool) zap.Field { return zap.Bool(key, val) }

// Any creates a zap field that accepts any value type.
func Any(key string, val interface{}) zap.Field { return zap.Any(key, val) }
