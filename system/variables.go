package system

type Command struct {
	Name        string
	Description string
}

var (
	Red         func(a ...interface{}) string
	Green       func(a ...interface{}) string
	Yellow      func(a ...interface{}) string
	Blue        func(a ...interface{}) string
	Magenta     func(a ...interface{}) string
	Cyan        func(a ...interface{}) string
	RedBold     func(a ...interface{}) string
	GreenBold   func(a ...interface{}) string
	YellowBold  func(a ...interface{}) string
	BlueBold    func(a ...interface{}) string
	MagentaBold func(a ...interface{}) string
	CyanBold    func(a ...interface{}) string
)

func InitSession(username string,
	session *Session) {
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

	// Set username in system var
	UserName = username

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
	Blue = colors["blue"]
	RedBold = colors["redBold"]
	YellowBold = colors["yellowBold"]
	CyanBold = colors["cyanBold"]
	MagentaBold = colors["magentaBold"]
	BlueBold = colors["blueBold"]
}

// BuildCommandMap создаёт карту команд для быстрого поиска.
func BuildCommandMap(commands []Command) map[string]struct{} {
	m := make(map[string]struct{}, len(commands))
	for _, cmd := range commands {
		m[cmd.Name] = struct{}{}
	}
	return m
}

// BuildStringCommandMap создаёт карту строковых команд для быстрого поиска.
func BuildStringCommandMap(commands []string) map[string]struct{} {
	m := make(map[string]struct{}, len(commands))
	for _, cmd := range commands {
		m[cmd] = struct{}{}
	}
	return m
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
	"debug":    &Debug,
}

var CmdMap map[string]struct{} = BuildCommandMap(Commands)                     // for system.Command
var AdditionalCmdMap map[string]struct{} = BuildCommandMap(AdditionalCommands) // for system.Command

var AvailableEditableVars = []string{"location", "prompt", "user", "empty", "user_dir", "debug"}
var CustomEditableVars []string

// R Previously entered command
var R string
