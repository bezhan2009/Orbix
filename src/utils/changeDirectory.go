package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithoutSignature/CD"
	"os"
)

func ChangeDirectoryUtil(commandArgs []string) {
	red := color.New(color.FgRed).SprintFunc()
	if len(commandArgs) == 0 {
		dir, _ := os.Getwd()
		fmt.Println(dir)
		return
	}
	if err := CD.ChangeDirectory(commandArgs[0]); err != nil {
		fmt.Println(red(err))
	}
}
