package system

import (
	"fmt"
	"github.com/google/uuid"
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
