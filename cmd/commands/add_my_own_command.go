package commands

import (
	"encoding/json"
	"fmt"
	utils2 "goCmd/utils"
	"os"
	"os/exec"
	"path/filepath"
)

var customCommands map[string]string

func initCustomCommands() {
	customCommands = make(map[string]string)

	err := loadCustomCommandsFromFile("custom_commands.json")
	if err != nil {
		fmt.Printf("Error loading custom commands: %v\n", err)
	}
}

func Start() {
	PrintAddCommand()

	var name string
	fmt.Println("Command Name:")
	fmt.Scan(&name)

	if _, exists := customCommands[name]; exists {
		fmt.Printf("Command '%s' already exists.\n", name)
		return
	}

	if !utils2.ValidCommand(name, Commands) {
		fmt.Println("Invalid Command")
		return
	}

	fmt.Println("Select the executable file for the command:")
	var exeFile string
	fmt.Scan(&exeFile)

	exeFileFullPath, err := validateExeFilePath(exeFile)
	if err != nil {
		fmt.Printf("Error validating executable file path: %v\n", err)
		return
	}

	customCommands[name] = exeFileFullPath
	fmt.Printf("Command '%s' added successfully with executable path '%s'.\n", name, exeFileFullPath)

	err = saveCustomCommandsToFile("custom_commands.json")
	if err != nil {
		fmt.Printf("Error saving custom commands: %v\n", err)
	}
}

func validateExeFilePath(exeFile string) (string, error) {
	if _, err := os.Stat(exeFile); os.IsNotExist(err) {
		fullPath, err := exec.LookPath(exeFile)
		if err != nil {
			return "", fmt.Errorf("executable file '%s' not found: %v", exeFile, err)
		}
		return fullPath, nil
	}

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
