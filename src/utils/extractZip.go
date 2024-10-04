package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands"
)

func ExtractZipUtil(commandArgs []string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	if len(commandArgs) < 2 {
		fmt.Println(yellow("Usage: extractzip <zip_file> <destination>"))
		return
	}
	if err := commands.ExtractZip(commandArgs[0], commandArgs[1]); err != nil {
		fmt.Println(red("Error extracting ZIP file:", err))
	}
}
