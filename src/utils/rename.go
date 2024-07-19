package utils

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/Rename"
	"goCmd/internal/debug"
)

func RenameFileUtil(commandArgs []string, command string, user, dir string) {
	if err := Rename.Rename(commandArgs); err != nil {
		debug.Commands(command, false, commandArgs, user, dir)
		fmt.Println(err)
	} else {
		debug.Commands(command, true, commandArgs, user, dir)
	}
}
