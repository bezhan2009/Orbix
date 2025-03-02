package system

type Command struct {
	Name        string
	Description string
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

// Flags — список флагов, которые нужно удалить
var Flags = []string{"--timing", "-t", "--run-in-new-thread"}

// R Previously entered command
var R string
