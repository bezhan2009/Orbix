package commands

import (
	"fmt"
	"goCmd/internal/OS"
	"goCmd/system"
)

func GetFileInfo(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: fileinfo <file path>")
		return
	}

	fileInfo, err := OS.GetFileInfo(commandArgs[0])
	if err != nil {
		fmt.Println(system.Red(fmt.Sprintf("%s\n", "Error getting file info: "+err.Error())))
		return
	}

	fmt.Println(system.GreenBold(fileInfo))
}
