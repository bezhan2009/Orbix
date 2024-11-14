package commands

import (
	"fmt"
	"goCmd/system"
	"os"
	"path/filepath"
)

func ChangeDirectory(path string) error {
	// Special case: If path is empty, change to the user's home directory
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("couldn't get the home directory: %v", err)
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
			return fmt.Errorf("couldn't get the current directory: %v", err)
		}
		parentDir := filepath.Dir(currentDir)
		path = parentDir
	}

	// Attempt to change directory
	err := os.Chdir(path)
	if err != nil {
		return fmt.Errorf("failed to change directory: %v", err)
	}

	// Print the new current directory, similar to how Windows `cd` works
	system.UserDir, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("couldn't get the current directory after the change: %v", err)
	}

	return nil
}
