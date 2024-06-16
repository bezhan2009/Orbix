package utils

import (
	"fmt"
	"goCmd/commands/commandsWithSignaiture/Remove"
	"goCmd/debug"
)

func RemoveFileUtil(commandArgs []string, command string, user, dir string) {
	name, err := Remove.File(commandArgs)
	if err != nil {
		debug.Commands(command, false, commandArgs, user, dir)
		fmt.Println(err)
	} else {
		debug.Commands(command, true, commandArgs, user, dir)
		fmt.Printf("File %s successfully deleted!\n", name)
	}
}
