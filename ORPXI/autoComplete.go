package ORPXI

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
	if len(parts) == 1 {
		return prompt.FilterHasPrefix(createUniqueCommandSuggestions(), text, true)
	} else {
		dir := "."
		if len(parts) > 2 {
			dir = strings.Join(parts[:len(parts)-1], " ")
		}
		return prompt.FilterHasPrefix(createFileSuggestions(dir), parts[len(parts)-1], true)
	}
}

func createUniqueCommandSuggestions() []prompt.Suggest {
	uniqueCommands := make(map[string]struct{})
	var suggestions []prompt.Suggest

	for _, cmd := range commands {
		if _, exists := uniqueCommands[cmd.Name]; !exists {
			uniqueCommands[cmd.Name] = struct{}{}
			suggestions = append(suggestions, prompt.Suggest{Text: cmd.Name, Description: cmd.Description})
		}
	}

	for _, cmd := range commandHistory {
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
	for _, cmd := range commands {
		if strings.HasPrefix(cmd.Name, input) {
			return cmd.Name
		}
	}
	return ""
}
