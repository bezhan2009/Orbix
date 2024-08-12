package template

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithSignaiture/Edit"
	"os"
)

func Make(commandArgs []string) {
	name := ""
	red := color.New(color.FgRed).SprintFunc()
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	if len(commandArgs) < 0 {
		fmt.Println(magenta("Template name:"))
		_, err := fmt.Scan(&name)
		if err != nil {
			fmt.Println(red(err))
		}
	}

	if len(commandArgs) < 1 {
		name = commandArgs[0]
		_, err := os.Create(name)
		if err != nil {
			fmt.Println(red(err))
		}
		err = Edit.File(name)
		if err != nil {
			fmt.Println(red(err))
		}
	}
}
