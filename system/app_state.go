package system

import (
	"os"
	"strings"
	"sync"
)

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

func Init() *AppState {
	BetaVersion = string(strings.TrimSpace(strings.ToLower(os.Getenv("BETA"))))

	Beta = SetBetaVersion(colors)

	if UserDir == "" {
		UserDir, _ = os.Getwd()
	}

	// Initialization AppState
	return NewSystemData()
}
