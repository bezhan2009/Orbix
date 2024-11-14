package system

import (
	"time"
)

const (
	// MaxRetryAttempts Maximum number of restart attempts
	MaxRetryAttempts = 5
	// RetryDelay Delay before restart
	RetryDelay = 1 * time.Second
)

var (
	Port                = "6060"
	ErrorStartingServer = false
	UserName            = ""
	OrbixWorking        = false
	Localhost           = ""
	UserDir             = ""
	GitHubURL           = "https://github.com/bezhan2009/Orbix"
)
