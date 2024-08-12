package template

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithSignaiture/Edit"
	"os"
)

func Make() {
	name := ""
	red := color.New(color.FgRed).SprintFunc()
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	fmt.Println(magenta("Template name:"))
	_, err := fmt.Scan(&name)
	if err != nil {
		fmt.Println(red(err))
	}

	_, err = os.Create(name)
	if err != nil {
		fmt.Println(red(err))
	}
	err = Edit.File(name)
	if err != nil {
		fmt.Println(red(err))
	}
}
