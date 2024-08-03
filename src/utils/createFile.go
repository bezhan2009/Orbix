package utils

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/Create"
	"goCmd/internal/logger"
	"path/filepath"
)

func CreateFileUtil(commandArgs []string, command, user, dir string) {
	name, err := Create.File(commandArgs)
	if err != nil {
		fmt.Println(err)
		logger.Commands(command, false, commandArgs, user, dir)
	} else if name != "" {
		fmt.Printf("File %s successfully created!\n", name)
		fmt.Printf("Directory of the new file: %s\n", filepath.Join(dir, name))
		logger.Commands(command, true, commandArgs, user, dir)
	}
}
