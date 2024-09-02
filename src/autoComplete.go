package src

import (
	"github.com/c-bata/go-prompt"
	"os"
	"strings"
)

func autoComplete(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()
	if len(text) == 0 {
		return []prompt.Suggest{}
	}

	parts := strings.Split(text, " ")

	lastChar := ""
	if d.Text != "" {
		lastChar = string(d.Text[len(d.Text)-1])
	}

	switch lastChar {
	case "(":
		return []prompt.Suggest{{Text: "()", Description: "Close parenthesis"}}
	case "[":
		return []prompt.Suggest{{Text: "[]", Description: "Close square bracket"}}
	case "{":
		return []prompt.Suggest{{Text: "{}", Description: "Close curly brace"}}
	case "\"":
		if strings.Count(d.Text, "\"")%2 == 1 {
			return []prompt.Suggest{{Text: "\"\"", Description: "Close double quote"}}
		}
	case "'":
		if strings.Count(d.Text, "''")%2 == 1 {
			return []prompt.Suggest{{Text: "'", Description: "Close single quote"}}
		}
	}

	// Подсказки для первого слова в команде (команды и файлы в текущем каталоге)
	if len(parts) == 1 {
		commandSuggestions := prompt.FilterHasPrefix(createUniqueCommandSuggestions(), text, true)
		fileSuggestions := prompt.FilterHasPrefix(createFileSuggestions("."), text, true)
		return append(commandSuggestions, fileSuggestions...)
	}

	// Подсказки для всех последующих слов (только файлы и ранее введенные команды)
	dir := "."
	if len(parts) > 2 {
		dir = strings.Join(parts[:len(parts)-1], " ")
	}
	fileSuggestions := prompt.FilterHasPrefix(createFileSuggestions(dir), parts[len(parts)-1], true)
	commandHistorySuggestions := prompt.FilterHasPrefix(createCommandHistorySuggestions(), parts[len(parts)-1], true)
	return append(fileSuggestions, commandHistorySuggestions...)
}

func createUniqueCommandSuggestions() []prompt.Suggest {
	uniqueCommands := make(map[string]struct{})
	var suggestions []prompt.Suggest

	// Assuming AdditionalCommands is a predefined list of commands
	for _, cmd := range AdditionalCommands {
		if _, exists := uniqueCommands[cmd.Name]; !exists {
			uniqueCommands[cmd.Name] = struct{}{}
			suggestions = append(suggestions, prompt.Suggest{Text: cmd.Name, Description: cmd.Description})
		}
	}

	return suggestions
}

func createCommandHistorySuggestions() []prompt.Suggest {
	uniqueCommands := make(map[string]struct{})
	var suggestions []prompt.Suggest

	// Assuming CommandHistory is a predefined slice of strings (previous commands)
	for _, cmd := range CommandHistory {
		if _, exists := uniqueCommands[cmd]; !exists {
			uniqueCommands[cmd] = struct{}{}
			suggestions = append(suggestions, prompt.Suggest{Text: cmd})
		}
	}

	return suggestions
}

func createFileSuggestions(dir string) []prompt.Suggest {
	files, err := os.ReadDir(dir)
	if err != nil {
		return []prompt.Suggest{}
	}

	var suggestions []prompt.Suggest
	for _, file := range files {
		suggestions = append(suggestions, prompt.Suggest{Text: file.Name()})
	}
	return suggestions
}

func suggestCommand(input string) string {
	for _, cmd := range AdditionalCommands {
		if strings.HasPrefix(cmd.Name, input) {
			return cmd.Name
		}
	}

	return ""
}
