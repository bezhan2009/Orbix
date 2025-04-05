package formats

import (
	"goCmd/system/formats/shortcut"
	"goCmd/system/formats/variable"
)

func IsValidFormatSV(input string) (bool, bool) {
	return shortcut.IsValidFormat(input), variable.IsValidFormat(input)
}

func ConvertSV(command string) string {
	shortcutCmd, variableCmd := IsValidFormatSV(command)
	if !(shortcutCmd || variableCmd) {
		return ""
	} else if shortcutCmd {
		return shortcut.ConvertShortcut(command)
	} else {
		return variable.СonvertSetVar(command)
	}
}

func UnconvertSV(command string) string {
	shortcutCmd := shortcut.UnconvertShortcut(command)
	if shortcutCmd != command {
		return shortcutCmd
	}

	return variable.UnconvertSetVar(command)
}
