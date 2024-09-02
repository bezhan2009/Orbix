package src

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"goCmd/system"
	"goCmd/utils"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// New helper functions
func initializeRunningFile(username string) {
	// Check and initialize running.txt if not exists
	if _, err := os.Stat("running.txt"); os.IsNotExist(err) {
		if _, err := os.Create("running.txt"); err != nil {
			panic(err)
		}
	}

	// Check for username in running.txt and add if missing
	runningPath := filepath.Join(Absdir, "running.txt")
	if sourceRunning, err := os.ReadFile(runningPath); err == nil {
		if !strings.Contains(string(sourceRunning), username) {
			if file, err := os.OpenFile("running.txt", os.O_APPEND|os.O_WRONLY, 0644); err == nil {
				defer file.Close()
				if _, err := file.WriteString(username + "\n"); err != nil {
					fmt.Println("Error writing to running.txt:", err)
				}
			}
		}
	}
}

func checkUserInRunningFile(username string) bool {
	runningPath := filepath.Join(Absdir, "running.txt")
	sourceRunning, err := os.ReadFile(runningPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(sourceRunning), username)
}

func printPromptInfo(location, user, dirC string, green, cyan, yellow, magenta func(...interface{}) string, commandInput string) {
	var optimize error
	fmt.Println(optimize)
	fmt.Printf("\n%s%s%s%s%s%s%s%s %s%s%s%s%s%s%s%s%s%s%s\n",
		yellow("┌"), yellow("─"), yellow("("), cyan("Orbix@"+user), yellow(")"), yellow("─"), yellow("["),
		yellow(location), magenta(time.Now().Format("15:04")), yellow("]"), yellow("─"), yellow("["),
		cyan("~"), cyan(dirC), yellow("]"), yellow(" git:"), green("["), green(system.GitBranch), green("]"))
	fmt.Printf("%s%s%s %s",
		yellow("└"), yellow("─"), green("$"), green(commandInput))
}

func readCommandLine(commandInput string) (string, string, []string, string) {
	var commandLine string
	if commandInput != "" {
		commandLine = strings.TrimSpace(commandInput)
	} else {
		commandLine = strings.TrimSpace(prompt.Input("", autoComplete))
	}

	commandParts := utils.SplitCommandLine(commandLine)
	if len(commandParts) == 0 {
		return "", "", nil, ""
	}

	command := commandParts[:1]

	return commandLine, command[0], commandParts[1:], strings.ToLower(commandParts[0])
}

func processCommand(commandLower string, commandArgs []string, dir string) error {
	if commandLower == "cd" && len(commandArgs) < 1 {
		fmt.Println(dir)
		SetGitBranch()
		return nil
	}

	if commandLower == "git" && len(commandArgs) > 2 {
		if commandArgs[0] == "switch" {
			SetGitBranch()
		}
	}

	if commandLower == "help" {
		displayHelp()
		return nil
	}

	if commandLower == "signout" {
		return fmt.Errorf("signout")
	}

	return nil
}