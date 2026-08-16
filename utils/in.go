package utils

import (
	"goCmd/system"
)

func InStr(value string, values []string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}

func InTurtle(value string, values []system.Turtle) (float64, bool) {
	for _, v := range values {
		if v.Command == value {
			return v.ValueFloat, true
		}
	}

	return 0, false
}
