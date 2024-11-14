package Remove

import (
	"fmt"
	"goCmd/system"
	"os"
	"strings"
)

func File(command string, commandArgs []string) (string, error) {
	colorsMap := system.GetColorsMap()
	yellow := colorsMap["yellow"]
	green := colorsMap["green"]
	red := colorsMap["red"]

	if len(commandArgs) < 1 {
		fmt.Println(yellow(fmt.Sprintf("Usage: %s <file>",
			strings.TrimSpace(strings.ToLower(command)))))
		return "", nil
	}

	name := commandArgs[0]

	if strings.TrimSpace(name) == "running.env" {
		fmt.Println(red("You can not delete file \"running.env\""))
		return "", nil
	}

	if err := os.Remove(name); err != nil {
		if strings.Contains(err.Error(), "being used by another process") {
			fmt.Println(red(fmt.Sprintf("Error: Cannot remove '%s' because it is being used by another process.\n", name)))
		} else {
			fmt.Println(red(fmt.Sprintf("Error removing file: %v\n", err)))
		}
		return name, err
	}

	if name != "" {
		fmt.Println(green(fmt.Sprintf("File '%s' successfully deleted.\n", name)))
	}

	return name, nil
}
