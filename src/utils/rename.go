package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithSignaiture/Rename"
	"goCmd/internal/logger"
)

func RenameFileUtil(commandArgs []string, command string, user, dir string) {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	if err := Rename.Rename(commandArgs, command); err != nil {
		err := logger.Commands(command, false, commandArgs, user, dir)
		if err != nil {
			fmt.Println(red(err))
		}
		fmt.Println(yellow(err))
	} else {
		err := logger.Commands(command, true, commandArgs, user, dir)
		if err != nil {
			fmt.Println(red(err))
		}
	}

	printSuccess := fmt.Sprintf("Successfully renamed file %s -> %s", commandArgs[0], commandArgs[1])
	fmt.Println(green(printSuccess))
}
