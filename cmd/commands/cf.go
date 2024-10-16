package commands

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func CreateFolder(commandArgs []string) (bool, error) {
	yellow := color.New(color.FgYellow).SprintFunc()

	if len(commandArgs) < 1 {
		fmt.Println(yellow("Usage: cf <folder_name>"))
		return false, nil
	}
	err := os.Mkdir(commandArgs[0], 0755) // 0755 - права доступа к директории
	if err != nil {
		return false, err
	}
	return true, nil
}
