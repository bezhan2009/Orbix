package service

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	_chan "goCmd/chan"
	"goCmd/system"
	"os"
	"sort"
	"strings"
)

func AutoComplete(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()

	text = strings.TrimSpace(strings.ToLower(text))

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

	if len(words) > 0 {
		// Получаем последний элемент
		lastElement := words[len(words)-1]

		// Разбиваем строку на части, разделенные оператором "+"
		parts := strings.Split(lastElement, "+")

		part := parts[len(parts)-1]
		if strings.HasPrefix(part, "$") {
			// Убираем "$" и передаем остаток для подсказок
			variableName := part[1:] // Убираем первый символ "$"
			commandVarsSuggestions := createUniqueVariableCommandSuggestions(lastElement, variableName)
			return commandVarsSuggestions
		}
	}

	// Если это первое слово, предлагать команды, файлы и историю команд
	if len(words) == 1 {
		commandSuggestions := createUniqueCommandSuggestions(words[0])
		fileSuggestions := createFileSuggestions(".", words[0])
		commandHistorySuggestions := createCommandHistorySuggestions(words[0])
		commandShortcuts := createUniqueShortcutsSuggestions(words[0])
		return removeDuplicateSuggestions(append(append(commandSuggestions, fileSuggestions...), append(commandShortcuts, commandHistorySuggestions...)...))
	}

	// После первого пробела предлагать только файлы и историю команд
	lastWord := words[len(words)-1]
	fileSuggestions := createFileSuggestions(".", lastWord)
	commandHistorySuggestions := createCommandHistorySuggestions(lastWord)
	commandShortcuts := createUniqueShortcutsSuggestions(words[0])

	return removeDuplicateSuggestions(append(append(fileSuggestions, commandHistorySuggestions...), commandShortcuts...))
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

	for _, cmd := range system.AdditionalCommands {
		if _, exists := uniqueCommands[strings.ToLower(cmd.Name)]; !exists && strings.HasPrefix(strings.ToLower(cmd.Name), prefix) {
			uniqueCommands[cmd.Name] = struct{}{}
			suggestions = append(suggestions, prompt.Suggest{Text: cmd.Name, Description: cmd.Description})
		}
	}

	return suggestions
}

func createUniqueShortcutsSuggestions(prefix string) []prompt.Suggest {
	uniqueCommands := make(map[string]struct{})
	var suggestions []prompt.Suggest

	for _, shortcut := range system.AvailableShortcuts {
		if _, exists := uniqueCommands[strings.ToLower(shortcut)]; !exists && strings.HasPrefix(strings.ToLower(shortcut), prefix) {
			uniqueCommands[shortcut] = struct{}{}
			suggestions = append(suggestions, prompt.Suggest{Text: shortcut, Description: fmt.Sprintf("Shortcut: %s", system.Shortcuts[shortcut])})
		}
	}

	return suggestions
}

func createUniqueVariableCommandSuggestions(word, prefix string) []prompt.Suggest {
	uniqueCommands := make(map[string]struct{})
	var suggestions []prompt.Suggest

	for _, vars := range system.AvailableEditableVars {
		if _, exists := uniqueCommands[strings.ToLower(vars)]; !exists && strings.HasPrefix(strings.ToLower(vars), prefix) {
			res := word + vars
			var sumStr string
			for _, variable := range strings.Split(res, "+") {
				value, _ := _chan.SetVarFn(variable[1:])
				sumStr = fmt.Sprintf("%s%s", sumStr, value)
			}

			if len(sumStr) > 10 && strings.Contains(res, "+") {
				sumStr = fmt.Sprintf("%s%s", "...", sumStr[10:])
			}

			uniqueCommands[vars] = struct{}{}
			suggestions = append(suggestions, prompt.Suggest{Text: word + vars, Description: sumStr})
		}
	}

	return suggestions
}

func createCommandHistorySuggestions(prefix string) []prompt.Suggest {
	uniqueCommands := make(map[string]struct{})
	var suggestions []prompt.Suggest

	// Проверяем, что GlobalSession и CommandHistory не равны nil
	if system.GlobalSession.CommandHistory == nil {
		system.GlobalSession.CommandHistory = []string{}
		system.InitSession(system.UserName, &system.GlobalSession)
		return suggestions
	}

	for _, cmd := range system.GlobalSession.CommandHistory {
		if system.GlobalSession.CommandHistory == nil {
			break
		}

		if cmd != "" && !strings.HasPrefix(strings.ToLower(cmd), prefix) {
			continue
		}
		if _, exists := uniqueCommands[strings.ToLower(cmd)]; !exists {
			uniqueCommands[cmd] = struct{}{}
			suggestions = append(suggestions, prompt.Suggest{
				Text:        cmd,
				Description: "previously entered command",
			})
		}
	}

	return suggestions
}

func createFileSuggestions(dir string, prefix string) []prompt.Suggest {
	files, err := os.ReadDir(dir)
	fileSuggestDescription := fmt.Sprintf("Finded in %s", system.UserDir)
	if err != nil {
		return []prompt.Suggest{}
	}

	var suggestions []prompt.Suggest
	for _, file := range files {
		if strings.HasPrefix(strings.ToLower(file.Name()), prefix) {
			suggestions = append(suggestions, prompt.Suggest{Text: file.Name(), Description: fileSuggestDescription})
		}
	}
	return suggestions
}

// SuggestCommand выполняет бинарный поиск по отсортированному списку команд
func SuggestCommand(input string) string {
	idx := sort.Search(len(system.AdditionalCommands), func(i int) bool {
		return system.AdditionalCommands[i].Name >= input
	})
	if idx < len(system.AdditionalCommands) && strings.HasPrefix(system.AdditionalCommands[idx].Name, input) {
		return system.AdditionalCommands[idx].Name
	}
	return ""
}
