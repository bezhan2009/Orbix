package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithSignaiture/Read"
)

func ReadFileUtil(commandLower string, commandArgs []string, user, dir string) {
	err := Read.File(commandLower, commandArgs, user, dir)
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Println(red(err))
	}
}
