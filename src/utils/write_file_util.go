package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithSignaiture/Write"
)

func WriteFileUtil(commandArgs []string) {
	err := Write.File(commandArgs)
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Println(red(err))
	}
}
