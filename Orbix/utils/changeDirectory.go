package utils

import (
	"fmt"
	"goCmd/commands/commandsWithoutSignature/CD"
	"os"
)

func ChangeDirectoryUtil(commandArgs []string) {
	if len(commandArgs) == 0 {
		dir, _ := os.Getwd()
		fmt.Println(dir)
		return
	}
	if err := CD.ChangeDirectory(commandArgs[0]); err != nil {
		fmt.Println(err)
	}
}
