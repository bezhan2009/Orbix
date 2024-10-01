package system

import "time"

const (
	// MaxRetryAttempts Maximum number of restart attempts
	MaxRetryAttempts = 5
	// RetryDelay Delay before restart
	RetryDelay = 1 * time.Second
)
