package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithSignaiture/Create"
	"goCmd/internal/logger"
	"path/filepath"
)

func CreateFileUtil(commandArgs []string, command, user, dir string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	name, err := Create.File(commandArgs)
	if err != nil {
		fmt.Println(err)
		err := logger.Commands(command, false, commandArgs, user, dir)
		if err != nil {
			fmt.Println(red("Error: " + err.Error()))
		}
	} else if name != "" {

		fmt.Printf(green(fmt.Sprintf("File %s successfully created!\n", name)))
		fmt.Printf(green("Directory of the new file: %s\n", filepath.Join(dir, name)))
		err := logger.Commands(command, true, commandArgs, user, dir)
		if err != nil {
			fmt.Println(red("Error: " + err.Error()))
		}
	}
}
