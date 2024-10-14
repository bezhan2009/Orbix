package src

import (
	"github.com/c-bata/go-prompt"
	"os"
	"strings"
)

func autoComplete(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()

	text = strings.TrimSpace(text)

	// Если ничего не введено, не показывать подсказки
	if len(text) == 0 {
		return []prompt.Suggest{}
	}

	lastChar := ""
	if strings.TrimSpace(d.Text) != "" {
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

	// Разделяем строку на слова
	words := strings.Fields(text)

	// Если это первое слово, предлагать команды, файлы и историю команд
	if len(words) == 1 {
		commandSuggestions := createUniqueCommandSuggestions(words[0])
		fileSuggestions := createFileSuggestions(".", words[0])
		commandHistorySuggestions := createCommandHistorySuggestions(words[0])
		return removeDuplicateSuggestions(append(append(commandSuggestions, fileSuggestions...), commandHistorySuggestions...))
	}

	// После первого пробела предлагать только файлы и историю команд
	lastWord := words[len(words)-1]
	fileSuggestions := createFileSuggestions(".", lastWord)
	commandHistorySuggestions := createCommandHistorySuggestions(lastWord)

	return removeDuplicateSuggestions(append(fileSuggestions, commandHistorySuggestions...))
}

// Функция для удаления дубликатов из списка подсказок
func removeDuplicateSuggestions(suggestions []prompt.Suggest) []prompt.Suggest {
	unique := make(map[string]struct{})
	var result []prompt.Suggest

	for _, suggestion := range suggestions {
		// Используем текст подсказки как ключ для проверки уникальности
		if _, exists := unique[suggestion.Text]; !exists {
			unique[suggestion.Text] = struct{}{}
			result = append(result, suggestion)
		}
	}

	return result
}

func createUniqueCommandSuggestions(prefix string) []prompt.Suggest {
	uniqueCommands := make(map[string]struct{})
	var suggestions []prompt.Suggest

	// Предполагается, что AdditionalCommands - это список доступных команд
	for _, cmd := range AdditionalCommands {
		if _, exists := uniqueCommands[cmd.Name]; !exists && strings.HasPrefix(cmd.Name, prefix) {
			uniqueCommands[cmd.Name] = struct{}{}
			suggestions = append(suggestions, prompt.Suggest{Text: cmd.Name, Description: cmd.Description})
		}
	}

	return suggestions
}

func createCommandHistorySuggestions(prefix string) []prompt.Suggest {
	uniqueCommands := make(map[string]struct{})
	var suggestions []prompt.Suggest

	// Предполагается, что CommandHistory - это слайс строк с историей команд
	for _, cmd := range GlobalSession.CommandHistory {
		if _, exists := uniqueCommands[cmd]; !exists && strings.HasPrefix(cmd, prefix) {
			uniqueCommands[cmd] = struct{}{}
			suggestions = append(suggestions, prompt.Suggest{Text: cmd})
		}
	}

	return suggestions
}

func createFileSuggestions(dir string, prefix string) []prompt.Suggest {
	files, err := os.ReadDir(dir)
	if err != nil {
		return []prompt.Suggest{}
	}

	var suggestions []prompt.Suggest
	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) {
			suggestions = append(suggestions, prompt.Suggest{Text: file.Name()})
		}
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