package circuitbreaker

import (
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	// StateClosed is the normal state where requests are allowed.
	StateClosed State = iota
	// StateOpen is the tripped state where requests are blocked.
	StateOpen
	// StateHalfOpen allows a limited number of probe requests to test recovery.
	StateHalfOpen
)

// Config holds circuit breaker configuration.
type Config struct {
	FailureThreshold    int           // consecutive failures before opening (default 5)
	OpenDuration        time.Duration // how long to stay open before half-open (default 30s)
	HalfOpenMaxRequests int           // successful probes needed to close (default 3)
}

type circuitState struct {
	State               State
	ConsecutiveFailures int
	ConsecutiveSuccess  int
	OpenUntil           time.Time
}

// CircuitBreaker tracks circuit state per key (e.g. provider:model).
type CircuitBreaker struct {
	mu     sync.Mutex
	states map[string]*circuitState
	config Config
}

// New creates a CircuitBreaker with the given configuration.
func New(cfg Config) *CircuitBreaker {
	if cfg.FailureThreshold == 0 {
		cfg.FailureThreshold = 5
	}
	if cfg.OpenDuration == 0 {
		cfg.OpenDuration = 30 * time.Second
	}
	if cfg.HalfOpenMaxRequests == 0 {
		cfg.HalfOpenMaxRequests = 3
	}
	return &CircuitBreaker{
		states: make(map[string]*circuitState),
		config: cfg,
	}
}

// IsAvailable checks if the circuit for the given key allows requests.
// Uses a write lock to atomically transition Open → HalfOpen, preventing
// multiple goroutines from simultaneously probing.
func (cb *CircuitBreaker) IsAvailable(key string) bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	s, ok := cb.states[key]
	if !ok {
		return true
	}

	switch s.State {
	case StateClosed:
		return true
	case StateOpen:
		if time.Now().After(s.OpenUntil) {
			s.State = StateHalfOpen
			s.ConsecutiveSuccess = 0
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return true
}

// RecordSuccess records a successful request for the given key.
func (cb *CircuitBreaker) RecordSuccess(key string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	s := cb.getOrCreate(key)
	s.ConsecutiveFailures = 0
	s.ConsecutiveSuccess++

	if s.State == StateHalfOpen && s.ConsecutiveSuccess >= cb.config.HalfOpenMaxRequests {
		s.State = StateClosed
		s.ConsecutiveSuccess = 0
	}
}

// RecordFailure records a failed request for the given key.
func (cb *CircuitBreaker) RecordFailure(key string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	s := cb.getOrCreate(key)
	s.ConsecutiveFailures++
	s.ConsecutiveSuccess = 0

	switch s.State {
	case StateClosed:
		if s.ConsecutiveFailures >= cb.config.FailureThreshold {
			s.State = StateOpen
			s.OpenUntil = time.Now().Add(cb.config.OpenDuration)
		}
	case StateHalfOpen:
		s.State = StateOpen
		s.OpenUntil = time.Now().Add(cb.config.OpenDuration)
	case StateOpen:
		s.OpenUntil = time.Now().Add(cb.config.OpenDuration)
	}
}

func (cb *CircuitBreaker) getOrCreate(key string) *circuitState {
	s, ok := cb.states[key]
	if !ok {
		s = &circuitState{State: StateClosed}
		cb.states[key] = s
	}
	return s
}
