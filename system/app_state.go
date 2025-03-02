package system

import "sync"

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
