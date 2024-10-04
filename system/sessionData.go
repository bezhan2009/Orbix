package system

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

var Attempts uint8
var Path string

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
func (s *AppState) NewSessionData(path, user, gitBranch string, isAdmin bool) (prefix string) {
	s.mu.Lock() // Защищаем доступ к карте сессий.
	defer s.mu.Unlock()

	// Проверяем пустые значения и заменяем их на дефолтные
	if user == "" {
		user = "unknown_user"
	}
	if gitBranch == "" {
		gitBranch = ""
	}

	// Используем текущую временную метку для уникальности
	timeStamp := time.Now().Format("20060102_150405") // Формат: YYYYMMDD_HHMMSS

	// Генерация уникального идентификатора (UUID)
	sessionID := uuid.New().String()

	// Формируем уникальный префикс
	prefix = fmt.Sprintf("%s_%s_%s_%s", timeStamp, gitBranch, user, sessionID)

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
