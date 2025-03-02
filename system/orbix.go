package system

import (
	"path/filepath"
	"time"
)

const (
	// MaxRetryAttempts Maximum number of restart attempts
	MaxRetryAttempts = 5
	// RetryDelay Delay before restart
	RetryDelay = 1 * time.Second
)

// Soft dynamic variables
var (
	Absdir, _   = filepath.Abs("")
	RunningPath = filepath.Join(Absdir, OrbixRunningUsersFileName)

	PreviousSessionPath   = ""
	PreviousSessionPrefix = ""

	Prefix = ""

	ExecutingCommand = false
	Unauthorized     = true

	RebootAttempts  = uint(0)
	SessionsStarted = uint(0)

	Neofetch = ""
)

// Orbix dynamic variable utils
var (
	Port                = "6060"
	Localhost           = ""
	ErrorStartingServer = false

	UserName = ""

	OrbixWorking = false

	UserDir = ""
)

const (
	// GitHubURL Orbix github repository
	GitHubURL = "https://github.com/bezhan2009/Orbix"
)
