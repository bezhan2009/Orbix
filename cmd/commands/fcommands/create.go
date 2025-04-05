package fcommands

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands"
)

func CFile(commandArgs []string) (string, error) {
	yellow := color.New(color.FgYellow).SprintFunc()
	if len(commandArgs) < 1 {
		fmt.Println(yellow("Usage: create <file>"))
		return "", nil
	}

	var name string

	name = commandArgs[0]

	errExisting := commands.IsExists(name)

	if errExisting == nil {
		return name, errExisting
	}

	nameFile, err := CreateFile(name)

	if err != nil {
		return string(nameFile), err
	}

	return string(nameFile), nil
}
