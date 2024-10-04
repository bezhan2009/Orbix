package src

import (
	"bufio"
	"fmt"
	"goCmd/structs"
	"goCmd/system"
	"os"
	"strings"
)

// Start executes a series of Commands defined in the given template file.
func Start(templateName string, echo string, SD *system.AppState) error {
	echoExecute := false

	if echo == "true" {
		echoExecute = true
	} else {
		echoExecute = false
	}

	templateName = strings.TrimSpace(templateName)

	file, err := os.OpenFile(templateName, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // Ignore empty lines
		}
		// Execute each command from the template
		if err := executeCommand(line, echoExecute, SD); err != nil {
			fmt.Printf("Error executing command '%s': %v\n", line, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return nil
}

// executeCommand executes a single command using the Orbix function.
func executeCommand(command string, echo bool, SD *system.AppState) error {
	// Assuming Orbix function handles the command execution.
	Orbix(command, echo, structs.RebootedData{}, SD)
	return nil
}
