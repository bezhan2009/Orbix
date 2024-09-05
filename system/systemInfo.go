package system

import (
	"runtime"
)

const (
	Version         = "1.1.8"
	License         = "MIT"
	SystemName      = "Orbix"
	OperationSystem = runtime.GOOS
)

var (
	// Path the value for this variable is given after the program is started
	Path = ""
	// User the value for this variable is given also after the program is started
	User = ""
	// IsAdmin It was added for the sake of optimization
	IsAdmin = true
	// GitBranch It is stored here, current git branch
	GitBranch = ""
)
