package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ExternalCommand executes an external command
func ExternalCommand(command []string) error {
	cmdName := command[0]
	cmdArgs := command[1:]

	cmd := exec.Command(cmdName, cmdArgs...)

	// Check if the command is an absolute path
	if !filepath.IsAbs(cmdName) {
		// First, try to find the command in the current directory
		if _, err := os.Stat(cmdName); err == nil {
			cmd.Path = cmdName
		} else {
			// Then try PATH environment variable
			if path, err := exec.LookPath(cmdName); err == nil {
				cmd.Path = path
			} else {
				return err
			}
		}
	} else {
		// If command is an absolute path, just set it
		cmd.Path = cmdName
	}

	// Special handling for Windows to ensure the correct executable is found
	if runtime.GOOS == "windows" && !strings.HasSuffix(cmd.Path, ".exe") && !strings.HasSuffix(cmd.Path, ".") {
		cmd.Path += ".exe"
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
