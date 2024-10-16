package Remove

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

func File(command string, commandArgs []string) (string, error) {
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	if len(commandArgs) < 1 {
		fmt.Println(yellow(fmt.Sprintf("Usage: %s <file>", command)))
		return "", nil
	}

	name := commandArgs[0]

	if name == "running.txt" {
		fmt.Println(red("You can not delete file \"running.txt\""))
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
