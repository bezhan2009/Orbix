package utils

import "strings"

func SplitCommandLine(input string) []string {
	var result []string
	var current strings.Builder
	inQuotes := false
	for _, r := range input {
		switch {
		case r == ' ' && !inQuotes:
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		case r == '"':
			inQuotes = !inQuotes
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		result = append(result, current.String())
	}
	return result
}
