package timeutil

import (
	"math"
	"time"
)

// Backoff calculates exponential backoff duration for a given retry attempt.
// Formula: initialInterval * 2^attempt, capped at maxInterval.
// Attempt is zero-based (0 = first retry).
func Backoff(attempt int, initialInterval, maxInterval time.Duration) time.Duration {
	d := time.Duration(float64(initialInterval) * math.Pow(2, float64(attempt)))
	if d > maxInterval {
		return maxInterval
	}
	return d
}

// ShouldRetry returns true if the attempt number is within the max retry limit.
func ShouldRetry(attempt, maxRetries int) bool {
	return attempt < maxRetries
}
