package system

var (
	Red     func(a ...interface{}) string
	Green   func(a ...interface{}) string
	Yellow  func(a ...interface{}) string
	blue    func(a ...interface{}) string
	Magenta func(a ...interface{}) string
	Cyan    func(a ...interface{}) string
)

func InitSession(session *Session) {
	// Initialize CommandHistory with package or tool names
	session.CommandHistory = append(session.CommandHistory, "help")
	session.CommandHistory = append(session.CommandHistory, "run")
	session.CommandHistory = append(session.CommandHistory, "push")
	session.CommandHistory = append(session.CommandHistory, "pull")
	session.CommandHistory = append(session.CommandHistory, "origin")
	session.CommandHistory = append(session.CommandHistory, "main")
	session.CommandHistory = append(session.CommandHistory, "master")
	session.CommandHistory = append(session.CommandHistory, "merge")
	session.CommandHistory = append(session.CommandHistory, "run")
	session.CommandHistory = append(session.CommandHistory, "start")
	session.CommandHistory = append(session.CommandHistory, ".")
	session.CommandHistory = append(session.CommandHistory, "remote")
	session.CommandHistory = append(session.CommandHistory, "newthread")
	session.CommandHistory = append(session.CommandHistory, "neofetch")
	session.CommandHistory = append(session.CommandHistory, "location")
	session.CommandHistory = append(session.CommandHistory, "diruser")
	session.CommandHistory = append(session.CommandHistory, "prompt")
	session.CommandHistory = append(session.CommandHistory, "remote -v")
	session.CommandHistory = append(session.CommandHistory, "add")
	session.CommandHistory = append(session.CommandHistory, "add .")
	session.CommandHistory = append(session.CommandHistory, "add README.md")
	session.CommandHistory = append(session.CommandHistory, "tasklist")
	session.CommandHistory = append(session.CommandHistory, "--version")
	session.CommandHistory = append(session.CommandHistory, "install")
	session.CommandHistory = append(session.CommandHistory, "django")
	session.CommandHistory = append(session.CommandHistory, "flask")
	session.CommandHistory = append(session.CommandHistory, "config")
	session.CommandHistory = append(session.CommandHistory, "--global")
	session.CommandHistory = append(session.CommandHistory, "--timing")
	session.CommandHistory = append(session.CommandHistory, "--run-in-new-thread")
	session.CommandHistory = append(session.CommandHistory, "-t")
	session.CommandHistory = append(session.CommandHistory, "-m")
	session.CommandHistory = append(session.CommandHistory, "-am")
	session.CommandHistory = append(session.CommandHistory, "--list")
	session.CommandHistory = append(session.CommandHistory, "getvar *")
	session.CommandHistory = append(session.CommandHistory, "\"Your name\"")
	session.CommandHistory = append(session.CommandHistory, "\"your_email@example.com\"")
	session.CommandHistory = append(session.CommandHistory, "config")
	session.CommandHistory = append(session.CommandHistory, "--global user.name")
	session.CommandHistory = append(session.CommandHistory, "--global user.email")
	session.CommandHistory = append(session.CommandHistory, "branch")
	session.CommandHistory = append(session.CommandHistory, "checkout")
	session.CommandHistory = append(session.CommandHistory, "status")
	session.CommandHistory = append(session.CommandHistory, "commit")
	session.CommandHistory = append(session.CommandHistory, "clone")
	session.CommandHistory = append(session.CommandHistory, "log")
	session.CommandHistory = append(session.CommandHistory, "rebase")
	session.CommandHistory = append(session.CommandHistory, "cherry-pick")
	session.CommandHistory = append(session.CommandHistory, "stash")
	session.CommandHistory = append(session.CommandHistory, "reset")
	session.CommandHistory = append(session.CommandHistory, "diff")
	session.CommandHistory = append(session.CommandHistory, "grep")
	session.CommandHistory = append(session.CommandHistory, "fetch")
	session.CommandHistory = append(session.CommandHistory, "remote add")
	session.CommandHistory = append(session.CommandHistory, "remote remove")
	session.CommandHistory = append(session.CommandHistory, "tag")
	session.CommandHistory = append(session.CommandHistory, "show")
	session.CommandHistory = append(session.CommandHistory, "revert")
	session.CommandHistory = append(session.CommandHistory, "rm")
	session.CommandHistory = append(session.CommandHistory, "mv")
	session.CommandHistory = append(session.CommandHistory, "apply")
	session.CommandHistory = append(session.CommandHistory, "3d")
	session.CommandHistory = append(session.CommandHistory, "2d")
	session.CommandHistory = append(session.CommandHistory, "font")
	session.CommandHistory = append(session.CommandHistory, "hello")
	session.CommandHistory = append(session.CommandHistory, "patch")
	session.CommandHistory = append(session.CommandHistory, "delete")
	session.CommandHistory = append(session.CommandHistory, "echo")
	session.CommandHistory = append(session.CommandHistory, "echo=on")
	session.CommandHistory = append(session.CommandHistory, "echo=off")
	session.CommandHistory = append(session.CommandHistory, "changelog")
	session.CommandHistory = append(session.CommandHistory, "beta")
	session.CommandHistory = append(session.CommandHistory, Localhost)
	session.CommandHistory = append(session.CommandHistory, GitHubURL)
	session.CommandHistory = append(session.CommandHistory, "upgrade")
	session.CommandHistory = append(session.CommandHistory, "export")
	session.CommandHistory = append(session.CommandHistory, "import")
	session.CommandHistory = append(session.CommandHistory, "tar compress")
	session.CommandHistory = append(session.CommandHistory, "tar decompress")
	session.CommandHistory = append(session.CommandHistory, "convert")
	session.CommandHistory = append(session.CommandHistory, "nmap monitor")

	// Initialize session data
	SetGitBranch(session)
	SetPath(session)

	Attempts = 0
}

func InitColors() {
	colors = GetColorsMap()

	Red = colors["red"]
	Yellow = colors["yellow"]
	Cyan = colors["cyan"]
	Green = colors["green"]
	Magenta = colors["magenta"]
	blue = colors["blue"]
}

var (
	Location = ""
	User     = ""
	Empty    = ""
	Prompt   = ""
)

// Flags — список флагов, которые нужно удалить
var Flags = []string{"--timing", "-t", "--run-in-new-thread"}

// EditableVars Map, ограничивающий изменяемые переменные
var EditableVars = map[string]interface{}{
	"location": &Location,
	"prompt":   &Prompt,
	"user":     &User,
	"empty":    &Empty,
	"user_dir": &UserDir,
}

var AvailableEditableVars = []string{"location", "prompt", "user", "empty", "user_dir"}
var CustomEditableVars []string
