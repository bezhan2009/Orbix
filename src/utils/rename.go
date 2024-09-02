package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithSignaiture/Rename"
)

func RenameFileUtil(commandArgs []string, command string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	err := Rename.Rename(commandArgs, command)
	if err != nil {
		fmt.Printf("%s %s\n", red("Error:"), err)
	}

	PrintSuccess := fmt.Sprintf("Successfully renamed file %s -> %s", commandArgs[0], commandArgs[1])
	fmt.Println(green(PrintSuccess))
}
