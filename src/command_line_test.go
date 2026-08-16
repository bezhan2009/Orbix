package src

import (
	"reflect"
	"testing"
)

func TestReadCommandLineSplitsFontEqualsForPrint(t *testing.T) {
	commandLine, command, args, commandLower := ReadCommandLine(
		"print cyan:ORBIX font=3d",
	)

	if commandLine != "print cyan:ORBIX font=3d" {
		t.Fatalf("unexpected command line: %q", commandLine)
	}

	if command != "print" {
		t.Fatalf("unexpected command: %q", command)
	}

	if commandLower != "print" {
		t.Fatalf("unexpected command lower: %q", commandLower)
	}

	expectedArgs := []string{"cyan:ORBIX", "font", "3d"}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Fatalf("unexpected args: got %#v, want %#v", args, expectedArgs)
	}
}
