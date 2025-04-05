package shortcut

import (
	"fmt"
	"regexp"
	"strings"
)

func IsValidFormat(input string) bool {
	re := regexp.MustCompile(`^#[a-zA-Z_][a-zA-Z0-9_]*\s*=\s*.+$`)
	return re.MatchString(input)
}

func ConvertShortcut(command string) string {
	parts := strings.SplitN(command, " ", 3)
	if len(parts) != 3 || parts[0] != "shortcut" {
		return ""
	}
	shortcut, value := parts[1], parts[2]
	if !IsValidFormat(shortcut) {
		return ""
	}
	return fmt.Sprintf("#%s = %s", shortcut, value)
}

func UnconvertShortcut(command string) string {
	parts := strings.SplitN(command, " ", 3)
	if len(parts) != 3 || !strings.Contains(parts[0], "#") || len(parts[0]) < 2 {
		return command
	}

	shortcut, value := parts[0][1:], parts[2]
	return fmt.Sprintf("shortcut %s %s", shortcut, value)
}
