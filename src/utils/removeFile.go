package utils

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/Remove"
	"goCmd/internal/debug"
)

func RemoveFileUtil(commandArgs []string, command string, user, dir string) {
	_, err := Remove.File(commandArgs, command)
	if err != nil {
		debug.Commands(command, false, commandArgs, user, dir)
		fmt.Println(err)
	} else {
		debug.Commands(command, true, commandArgs, user, dir)
	}
}
