package utils

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

func Getwd() (wd string) {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	var err error
	wd, err = os.Getwd()
	if err != nil {
		fmt.Println(yellow("WARNING: Some commands may not work because the Getwd function failed with an error"))
		fmt.Println(red(err))
	}

	wd = strings.TrimSpace(wd)

	return
}
