package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithSignaiture/Write"
)

func WriteFileUtil(commandLower string, commandArgs []string, user string, dir string) {
	err := Write.File(commandLower, commandArgs, user, dir)
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Println(red(err))
	}
}
