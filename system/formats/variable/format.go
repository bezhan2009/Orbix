package variable

import (
	"fmt"
	"regexp"
	"strings"
)

func IsValidFormat(input string) bool {
	re := regexp.MustCompile(`^\$[a-zA-Z_][a-zA-Z0-9_]*\s*=\s*.+$`)
	return re.MatchString(input)
}

func СonvertSetVar(command string) string {
	parts := strings.SplitN(command, " ", 3)
	if len(parts) != 3 || parts[0] != "setvar" {
		return ""
	}
	variable, value := parts[1], parts[2]
	if !IsValidFormat(variable) {
		return ""
	}
	return fmt.Sprintf("$%s = %s", variable, value)
}

func UnconvertSetVar(command string) string {
	parts := strings.SplitN(command, " ", 3)
	if len(parts) != 3 || !strings.Contains(parts[0], "$") || len(parts[0]) < 2 {
		return command
	}

	shortcut, value := parts[0][1:], parts[2]
	return fmt.Sprintf("setvar %s %s", shortcut, value)
}
