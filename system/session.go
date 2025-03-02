package system

import (
	"fmt"
	"github.com/google/uuid"
	"time"
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

var Attempts uint8
var Path string

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
	// R previously entered command
	R string
}

// NewSessionData создает новую сессию и добавляет её в карту сессий.
func (s *AppState) NewSessionData(path, user, gitBranch string,
	isAdmin bool) (prefix string) {
	s.mu.Lock() // Защищаем доступ к карте сессий.
	defer s.mu.Unlock()

	// Проверяем пустые значения и заменяем их на дефолтные
	if user == "" {
		user = "unknown_user"
	}

	// Используем текущую временную метку для уникальности
	timeStamp := time.Now().Format("20060102_150405") // Формат: YYYYMMDD_HHMMSS

	// Генерация уникального идентификатора (UUID)
	sessionID := uuid.New().String()

	// Формируем уникальный префикс
	prefix = fmt.Sprintf("%s.%s.%s.%s", timeStamp, gitBranch, user, sessionID)

	// Добавляем сессию в карту
	s.Session[prefix] = Session{
		Path:      path,
		User:      user,
		GitBranch: gitBranch,
		IsAdmin:   isAdmin,
	}

	return prefix
}

// GetSession возвращает данные сессии по префиксу.
func (s *AppState) GetSession(prefix string) (*Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.Session[prefix]
	return &session, exists
}

// DeleteSession удаляет сессию по префиксу.
func (s *AppState) DeleteSession(prefix string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.Session, prefix)
}
