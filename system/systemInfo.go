package system

import (
	"runtime"
)

const (
	Version         = "1.0.11"
	License         = "MIT"
	SystemName      = "Orbix"
	OperationSystem = runtime.GOOS
)

var (
	// Path the value for this variable is given after the program is started
	Path = ""
)
