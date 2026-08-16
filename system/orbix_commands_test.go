package system

import "testing"

func TestTurtleRegisteredAsPrimaryCommand(t *testing.T) {
	if _, ok := CmdMap["turtle"]; !ok {
		t.Fatal("turtle must be registered in CmdMap for template execution")
	}

	if _, ok := AdditionalCmdMap["turtle"]; !ok {
		t.Fatal("turtle must remain registered in AdditionalCmdMap")
	}
}
