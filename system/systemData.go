package system

var (
	// Path the value for this variable is given after the program is started
	Path = ""
	// User the value for this variable is given also after the program is started
	User = ""
	// IsAdmin It was added for the sake of optimization
	IsAdmin = true
	// GitBranch It is stored here, current git branch
	GitBranch = ""
	// Attempts The number of attempts to recover from errors
	Attempts = 0
)
