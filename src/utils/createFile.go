package utils

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/Create"
	"goCmd/debug"
	"path/filepath"
)

func CreateFileUtil(commandArgs []string, command, user, dir string) {
	name, err := Create.File(commandArgs)
	if err != nil {
		fmt.Println(err)
		debug.Commands(command, false, commandArgs, user, dir)
	} else if name != "" {
		fmt.Printf("File %s successfully created!\n", name)
		fmt.Printf("Directory of the new file: %s\n", filepath.Join(dir, name))
		debug.Commands(command, true, commandArgs, user, dir)
	}
}
