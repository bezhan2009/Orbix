package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands"
	"goCmd/system"
	"os"
	"strings"
)

func ChangeDirectoryUtil(commandArgs []string, session *system.Session) {
	red := color.New(color.FgRed).SprintFunc()
	dir, _ := os.Getwd()
	if len(commandArgs) == 0 {
		fmt.Println(dir)
		return
	}

	changeDir := strings.Join(commandArgs, " ")

	if err := commands.ChangeDirectory(changeDir); err != nil {
		fmt.Println(red(err))
	}

	dir, _ = os.Getwd()
	session.Path = dir
}
