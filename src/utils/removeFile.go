package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithSignaiture/Remove"
	"goCmd/internal/logger"
)

func RemoveFileUtil(commandArgs []string, command string, user, dir string) {
	red := color.New(color.FgRed).SprintFunc()

	_, err := Remove.File(command, commandArgs)
	if err != nil {
		err := logger.Commands(command, false, commandArgs, user, dir)
		if err != nil {
			fmt.Println(red(err))
		}
		fmt.Println(red(err))
	} else {
		err := logger.Commands(command, true, commandArgs, user, dir)
		if err != nil {
			fmt.Println(red(err))
		}
	}
}
