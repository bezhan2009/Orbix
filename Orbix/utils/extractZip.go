package utils

import (
	"fmt"
	"goCmd/commands/commandsWithSignaiture/ExtractZip"
)

func ExtractZipUtil(commandArgs []string) {
	if len(commandArgs) < 2 {
		fmt.Println("Usage: extractzip <zipfile> <destination>")
		return
	}
	if err := ExtractZip.ExtractZip(commandArgs[0], commandArgs[1]); err != nil {
		fmt.Println("Error extracting ZIP file:", err)
	}
}
