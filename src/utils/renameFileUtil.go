package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithSignaiture/Read"
)

func ReadFileUtil(commandArgs []string) {
	err := Read.File(commandArgs)
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Println(red(err))
	}
}
