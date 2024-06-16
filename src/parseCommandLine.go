package src

import "strings"

func parseCommandLine(commandLine string) []string {
	var parts []string
	var currentPart strings.Builder
	var inQuotes bool

	for _, char := range commandLine {
		switch char {
		case '"':
			inQuotes = !inQuotes
		case ' ':
			if inQuotes {
				currentPart.WriteRune(char)
			} else {
				parts = append(parts, currentPart.String())
				currentPart.Reset()
			}
		default:
			currentPart.WriteRune(char)
		}
	}

	if currentPart.Len() > 0 {
		parts = append(parts, currentPart.String())
	}

	return parts
}
