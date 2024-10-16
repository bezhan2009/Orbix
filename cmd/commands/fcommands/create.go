package fcommands

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands"
)

func CFile(commandArgs []string) (string, error) {
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	if len(commandArgs) < 1 {
		fmt.Println(yellow("Usage: create <file>"))
		return "", nil
	}

	var name string

	name = commandArgs[0]

	if name == "debug.txt" {
		fmt.Println(red("PermissionDenied: You cannot write, delete or create a debug.txt file"))
		return name, nil
	}

	errExisting := commands.IsExists(name)

	if errExisting == nil {
		return name, errExisting
	}

	name, err := CreateFile(name)

	if err != nil {
		return name, err
	} else {
		return name, err
	}
}
