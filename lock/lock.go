package lock

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// ErrLockNotAcquired is returned when the lock is already held by another process.
var ErrLockNotAcquired = errors.New("lock: not acquired")

// Locker provides Redis-based distributed locking with token ownership.
// Each lock is identified by a key and protected by a random token to prevent
// accidental release by a different process.
type Locker struct {
	rdb *redis.Client
}

// New creates a Locker backed by the given Redis client.
func New(rdb *redis.Client) *Locker {
	return &Locker{rdb: rdb}
}

// TryLock attempts to acquire a lock with the given key and TTL.
// On success, returns an unlock function that must be called to release the lock.
// Returns ErrLockNotAcquired if the lock is already held.
func (l *Locker) TryLock(ctx context.Context, key string, ttl time.Duration) (unlock func(), err error) {
	token := uuid.New().String()

	ok, err := l.rdb.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrLockNotAcquired
	}

	unlock = func() {
		// Lua script: atomic check-and-delete to prevent deleting another process's lock
		l.rdb.Eval(context.Background(), luaRelease, []string{key}, token)
	}
	return unlock, nil
}

// ExecuteWithLock acquires the lock, executes fn, then releases the lock.
// Returns ErrLockNotAcquired if the lock is already held.
// The lock is always released after fn completes, even if fn returns an error.
func (l *Locker) ExecuteWithLock(ctx context.Context, key string, ttl time.Duration, fn func() error) error {
	unlock, err := l.TryLock(ctx, key, ttl)
	if err != nil {
		return err
	}
	defer unlock()
	return fn()
}

// luaRelease atomically checks the token and deletes the key only if it matches.
const luaRelease = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`
