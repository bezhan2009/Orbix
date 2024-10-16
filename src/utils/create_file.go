package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/fcommands"
	"path/filepath"
)

func CreateFileUtil(commandArgs []string, dir string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	name, err := fcommands.CFile(commandArgs)
	if err != nil {
		fmt.Println(red(err))
	} else if name != "" {
		fmt.Printf(green(fmt.Sprintf("CFile %s successfully created!\n", name)))
		fmt.Println(green(fmt.Sprintf("Directory of the new file: %s\n", filepath.Join(dir, name))))
	}
}
