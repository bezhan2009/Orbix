package AddOwnCommand

import (
	"encoding/json"
	"fmt"
	"goCmd/commands/commandsWithSignaiture/AddOwnCommand/utils"
	utils2 "goCmd/utils"
	"os"
	"os/exec"
	"path/filepath"
)

var customCommands map[string]string

func init() {
	// Initialize customCommands map
	customCommands = make(map[string]string)

	// Load custom commands from file on application startup
	err := loadCustomCommandsFromFile("custom_commands.json")
	if err != nil {
		fmt.Printf("Error loading custom commands: %v\n", err)
	}
}

func Start() {
	utils.PrintAddCommand()

	var name string
	fmt.Println("Command Name:")
	fmt.Scan(&name)

	// Check if command already exists
	if _, exists := customCommands[name]; exists {
		fmt.Printf("Command '%s' already exists.\n", name)
		return
	}

	// Validate the command name
	if !utils2.ValidCommand(name, utils.Commands) {
		fmt.Println("Invalid Command")
		return
	}

	fmt.Println("Select the executable file for the command:")
	var exeFile string
	fmt.Scan(&exeFile)

	// Validate executable file path
	exeFileFullPath, err := validateExeFilePath(exeFile)
	if err != nil {
		fmt.Printf("Error validating executable file path: %v\n", err)
		return
	}

	// Store the command and executable file path in customCommands
	customCommands[name] = exeFileFullPath
	fmt.Printf("Command '%s' added successfully with executable path '%s'.\n", name, exeFileFullPath)

	// Save custom commands to file after adding new command
	err = saveCustomCommandsToFile("custom_commands.json")
	if err != nil {
		fmt.Printf("Error saving custom commands: %v\n", err)
	}
}

func validateExeFilePath(exeFile string) (string, error) {
	// Check if the file exists
	if _, err := os.Stat(exeFile); os.IsNotExist(err) {
		// Try to find the file in the current directory or along the PATH
		fullPath, err := exec.LookPath(exeFile)
		if err != nil {
			return "", fmt.Errorf("executable file '%s' not found: %v", exeFile, err)
		}
		return fullPath, nil
	}

	// If the file exists, return its absolute path
	return filepath.Abs(exeFile)
}

func saveCustomCommandsToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(customCommands)
	if err != nil {
		return fmt.Errorf("error encoding custom commands to JSON: %v", err)
	}

	return nil
}

func loadCustomCommandsFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&customCommands)
	if err != nil {
		return fmt.Errorf("error decoding custom commands from JSON: %v", err)
	}

	return nil
}

// GetCustomCommands Function to retrieve custom commands and their executable paths
func GetCustomCommands() map[string]string {
	return customCommands
}
