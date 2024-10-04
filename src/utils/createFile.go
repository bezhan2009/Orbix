package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands"
	"path/filepath"
)

func CreateFileUtil(commandArgs []string, dir string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	name, err := commands.File(commandArgs)
	if err != nil {
		fmt.Println(red(err))
	} else if name != "" {
		fmt.Printf(green(fmt.Sprintf("EditFile %s successfully created!\n", name)))
		fmt.Printf(green("Directory of the new file: %s\n", filepath.Join(dir, name)))
	}
}
