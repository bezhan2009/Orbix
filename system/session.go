package system

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

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
}

// AppState представляет структуру для хранения данных системы.
type AppState struct {
	Path      string
	User      string
	IsAdmin   bool
	GitBranch string
	Session   map[string]Session
	mu        sync.Mutex
}

// NewSystemData инициализирует новую структуру AppState.
func NewSystemData() *AppState {
	return &AppState{
		Session: make(map[string]Session),
	}
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

var GitExists = CheckPackageExists("git")

func GetCurrentGitBranch() (string, error) {
	if !GitExists {
		ErrGitNotInstalled := errors.New("ErrGitNotInstalled")
		return "", ErrGitNotInstalled
	}

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	branch := out.String()
	return branch[:len(branch)-1], nil
}

func Getwd() (wd string) {
	var err error
	wd, err = os.Getwd()
	if err != nil {
		fmt.Println(Yellow("WARNING: Some commands may not work because the Getwd function failed with an error"))
		fmt.Println(Red(err))
	}

	wd = strings.TrimSpace(wd)

	return
}

func SetGitBranch(sd *Session) {
	var errGitBranch error
	sd.GitBranch, errGitBranch = GetCurrentGitBranch()
	if errGitBranch != nil {
		sd.GitBranch = ""
	}
}

func SetPath(sd *Session) {
	sd.Path = Getwd()
}
