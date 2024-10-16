package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/Remove"
)

func RemoveFileUtil(commandArgs []string, command string) {
	red := color.New(color.FgRed).SprintFunc()

	_, err := Remove.File(command, commandArgs)
	if err != nil {
		fmt.Println(red(err))
	}
}
