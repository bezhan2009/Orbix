package src

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/system"
)

func ExecuteTemplateUtil(commandArgs []string, SD *system.AppState) {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	if len(commandArgs) < 2 {
		fmt.Println(yellow("Usage: template <template_name> echo=on"))
		fmt.Println(yellow("Or: template <template_name> echo=off if you want without outputting the result"))
		return
	}

	if commandArgs[1] != "echo=on" && commandArgs[1] != "echo=off" {
		commandArgs[1] = "true"
	} else if commandArgs[1] == "echo=on" {
		commandArgs[1] = "true"
	} else if commandArgs[1] == "echo=off" {
		commandArgs[1] = "false"
	}

	if err := Start(commandArgs[0], commandArgs[1], SD); err != nil {
		fmt.Println(red(err))
	}
}
