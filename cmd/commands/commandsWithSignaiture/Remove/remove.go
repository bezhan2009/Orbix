package Remove

import (
	"fmt"
	"os"
	"strings"
)

func File(commandArgs []string) (string, error) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: remove <file>")
		return "", nil
	}

	name := commandArgs[0]

	if err := os.Remove(name); err != nil {
		if strings.Contains(err.Error(), "being used by another process") {
			fmt.Printf("Error: Cannot remove '%s' because it is being used by another process.\n", name)
		} else {
			fmt.Printf("Error removing file: %v\n", err)
		}
		return name, err
	}

	if name != "" {
		fmt.Printf("File '%s' successfully removed.\n", name)
	}

	return name, nil
}
