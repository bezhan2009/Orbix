package template

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands"
	"goCmd/system"
	"os"
)

func Make(commandArgs []string) {
	name := ""
	red := color.New(color.FgRed).SprintFunc()
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	if len(commandArgs) == 0 {
		fmt.Print(magenta("Template name:"))
		_, err := fmt.Scan(&name)
		if err != nil {
			fmt.Println(red(err))
			return
		}
	} else {
		name = commandArgs[0]
	}

	_, err := os.Create(fmt.Sprintf("%s.%s", name, system.OrbixTemplatesExtension))
	if err != nil {
		fmt.Println(red(err))
	}
	err = commands.EditFile(fmt.Sprintf("%s.%s", name, system.OrbixTemplatesExtension))
	if err != nil {
		fmt.Println(red(err))
	}
}
