package Rename

import (
	"fmt"
	"os"
	"strings"
)

func Rename(commandArgs []string, command string) error {
	if len(commandArgs) < 2 {
		if len(commandArgs) == 1 {
			fmt.Printf("Usage: %s %s <new name for file>\n", command, commandArgs[0])
			return nil
		}
		fmt.Printf("Usage: %s <file> <new name for file>\n", command)
		return nil
	}

	oldName := commandArgs[0]
	newName := commandArgs[1]

	if _, err := os.Stat(newName); err == nil {
		fmt.Printf("Error: A file named '%s' already exists.\n", newName)
		return fmt.Errorf("file '%s' already exists", newName)
	}

	if err := os.Rename(oldName, newName); err != nil {
		if strings.Contains(err.Error(), "being used by another process") {
			fmt.Printf("Error: Cannot rename '%s' because it is being used by another process.\n", oldName)
		} else {
			fmt.Printf("Error renaming file: %v\n", err)
		}
		return err
	}

	fmt.Printf("File '%s' successfully renamed to '%s'.\n", oldName, newName)
	return nil
}
