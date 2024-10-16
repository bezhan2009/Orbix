package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/Rename"
)

func RenameFileUtil(commandArgs []string, command string, yellow func(a ...interface{}) string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	if len(commandArgs) < 2 {
		if len(commandArgs) == 1 {
			fmt.Println(yellow(fmt.Sprintf("Usage: %s %s <new name for file>", command, commandArgs[0])))
			return
		}
		fmt.Println(yellow(fmt.Sprintf("Usage: %s <file> <new name for file>", command)))
		return
	}

	err := Rename.Rename(commandArgs, command)
	if err != nil {
		fmt.Printf("%s %s\n", red("Error:"), err)
	}

	PrintSuccess := fmt.Sprintf("Successfully renamed file %s -> %s", commandArgs[0], commandArgs[1])
	fmt.Println(green(PrintSuccess))
}
