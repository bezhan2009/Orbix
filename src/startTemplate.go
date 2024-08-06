package src

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Start executes a series of commands defined in the given shablon file.
func Start(shablonName string, echo string) error {
	echoExecute := false

	if echo == "true" {
		echoExecute = true
	} else {
		echoExecute = false
	}

	shablonName = strings.TrimSpace(shablonName)

	file, err := os.OpenFile(shablonName, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // Ignore empty lines
		}
		// Execute each command from the template
		if err := executeCommand(line, echoExecute); err != nil {
			fmt.Printf("Error executing command '%s': %v\n", line, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return nil
}

// executeCommand executes a single command using the Orbix function.
func executeCommand(command string, echo bool) error {
	// Assuming Orbix function handles the command execution.
	Orbix(command, echo)
	return nil
}
