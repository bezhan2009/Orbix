package system

import (
	"os"
	"runtime"
)

const (
	Version    = "1.11.2"
	License    = "MIT"
	SystemName = "Orbix"

	OperationSystem = runtime.GOOS

	MaxInt              = int(^uint64(0) >> 1)
	MaxUserAuthAttempts = uint(3)

	OrbixRunningUsersFileName = "running.env"
	OrbixTemplatesExtension   = "tmpl"
)

var (
	Beta        = false
	BetaVersion = ""

	colors = SetColorsMap()

	GlobalSession = Session{}

	SourcePath, _ = os.Getwd()

	LaunchedOrbixes    = make(map[string]string)
	CntLaunchedOrbixes = uint(0)

	Debug = false
)

var (
	OrbixFileNames = map[string]uint{
		OrbixRunningUsersFileName: 1,
		"user.json":               1,
		".env":                    1,
		"commands.json":           1,
	}
	OrbixUser       = &User
	OrbixRecovering = false
)
