package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithSignaiture/copySource"
)

func CommandCopySourceUtil(commandArgs []string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	if len(commandArgs) < 2 {
		fmt.Println(yellow("Usage: copysource <srcDir> <dstDir>"))
		fmt.Println(yellow("Example: copysource example.txt destination_directory.txt:"))
		fmt.Println(yellow("copysource example.txt bufer"))
		fmt.Println(yellow("Arguments:"))
		fmt.Println(yellow("copysource example.txt bufer"))
		return
	}
	err := copySource.File(commandArgs[0], commandArgs[1])
	if err != nil {
		fmt.Println(red(err))
	}
}
