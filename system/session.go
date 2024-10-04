package system

type Session struct {
	// Path the value for this variable is given after the program is started
	Path string
	// PreviousPath previous sessions path
	PreviousPath string
	// User the value for this variable is given also after the program is started
	User string
	// IsAdmin It was added for the sake of optimization
	IsAdmin bool
	// GitBranch It is stored here, current git branch
	GitBranch string
	// CommandHistory Users command history
	CommandHistory []string
}
