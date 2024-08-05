package src

import (
	"fmt"
)

func ExecuteShablonUtil(commandArgs []string) {
	if len(commandArgs) < 2 {
		fmt.Println("Usage: shablon <template_name> echo=on")
		fmt.Println("Or: shablon <template_name> echo=off if you want without outputting the result")
		return
	}

	if commandArgs[1] != "echo=on" && commandArgs[1] != "echo=off" {
		commandArgs[1] = "true"
	} else if commandArgs[1] == "echo=on" {
		commandArgs[1] = "true"
	} else if commandArgs[1] == "echo=off" {
		commandArgs[1] = "false"
	}

	if err := Start(commandArgs[0], commandArgs[1]); err != nil {
		fmt.Println(err)
	}
}
