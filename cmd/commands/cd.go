package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

func ChangeDirectory(path string) error {
	// Special case: If path is empty, change to the user's home directory
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("не удалось получить домашнюю директорию: %v", err)
		}
		path = homeDir
	}

	// Special case: If path is ".", stay in the current directory
	if path == "." {
		return nil
	}

	// Special case: If path is "..", move to the parent directory
	if path == ".." {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("не удалось получить текущую директорию: %v", err)
		}
		parentDir := filepath.Dir(currentDir)
		path = parentDir
	}

	// Attempt to change directory
	err := os.Chdir(path)
	if err != nil {
		return fmt.Errorf("не удалось сменить директорию: %v", err)
	}

	// Print the new current directory, similar to how Windows `cd` works
	_, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("не удалось получить текущую директорию после смены: %v", err)
	}

	return nil
}
