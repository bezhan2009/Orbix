package ORPXI

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"goCmd/cmdPress"
	"goCmd/utils"
	"os"
	"path/filepath"
	"strings"
)

func CMD(commandInput string) {
	utils.SystemInformation()

	isWorking := true
	isPermission := true

	var promptText string

	isEmpty, err := isPasswordDirectoryEmpty()
	if err != nil {
		fmt.Println("Error checking password directory:", err)
		return
	}

	if !isEmpty {
		if commandInput == "" {
			dir, _ := os.Getwd()
			user := cmdPress.CmdUser(dir)

			if !CheckUser(user) {
				return
			}
		}
	}

	for isWorking {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		dir, _ := os.Getwd()
		dirC := cmdPress.CmdDir(dir)
		user := cmdPress.CmdUser(dir)

		if promptText != "" {
			fmt.Print("\n" + promptText)
		} else {
			fmt.Printf("\n┌─(%s)-[%s%s]\n", cyan("ORPXI "+user), cyan("~"), cyan(dirC))
			fmt.Printf("└─$ %s", green(">", commandInput))
		}

		commandLine := prompt.Input("", autoComplete)
		commandLine = strings.TrimSpace(commandLine)
		commandParts := parseCommandLine(commandLine)

		if len(commandParts) == 0 {
			continue
		}

		command := commandParts[0]
		commandArgs := commandParts[1:]
		commandLower := strings.ToLower(command)

		commandHistory = append(commandHistory, commandLine)

		fmt.Println()

		if commandLower == "prompt" {
			if len(commandArgs) < 1 {
				fmt.Println("prompt <name_prompt>")
				fmt.Println("to delete prompt enter:")
				fmt.Println("prompt delete")
				continue
			}

			namePrompt := commandArgs[0]

			if namePrompt != "delete" {
				namePrompt = strings.TrimSpace(namePrompt)
				promptText = namePrompt
				fmt.Println("Prompt set to:", promptText)
			} else {
				promptText, _ = os.Getwd()
				fmt.Println("Prompt set to:", promptText)
				promptText = ""
			}

			continue
		}

		helpText := `
For command information, type HELP
CREATE             creates a new file
CLEAN              clears the screen
CD                 changes the current directory
LS                 lists directory contents
NEWSHABLON         creates a new command template for execution
REMOVE             deletes a file
READ               displays the contents of a file
PROMPT             changes ORPXI prompt.
PINGVIEW           shows ping.
PRIMES             finds large prime numbers
PICALC             calculates the value of π.
NEWUSER            creates a new user for ORPXI.
ORPXI              starts another ORPXI session
SHABLON            executes a specific command template
SYSTEMGOCMD        displays information about ORPXI
SYSTEMINFO         displays system information
SIGNOUT            user signs out of ORPXI
TREE               graphically displays directory structure
WRITE              writes data to a file
EDIT               edits a file
WIFIUTILS          launches a utility for working with WiFi
EXTRACTZIP         extracts .zip archives
SCANPORT           scans ports
WHOIS              domain information
DNSLOOKUP          DNS queries
FILEIO             intensive file operation test
IPINFO             IP address information
GEOIP              IP address geolocation
MATRIXMUL          multiplies large matrices
EXIT               exit
`

		if commandLower == "help" {
			fmt.Println(helpText)
			continue
		}

		isValid := utils.ValidCommand(commandLower, commands)

		if !isValid {
			fullCommand := append([]string{command}, commandArgs...)
			err := utils.ExternalCommand(fullCommand)
			if err != nil {
				fullPath := filepath.Join(dir, command)
				fullCommand[0] = fullPath
				err = utils.ExternalCommand(fullCommand)
				if err != nil {
					suggestedCommand := suggestCommand(commandLower)
					fmt.Printf("Error executing command '%s': %v\n", commandLine, err)
					if suggestedCommand != "" {
						fmt.Printf("Did you mean: %s?\n", suggestedCommand)
					}
				}
			}
			continue
		}

		if commandInput != "" {
			isPermission = false
		} else {
			isPermission = true
		}

		ExecuteCommand(commandLower, command, commandLine, dir, commands, commandArgs, &isWorking, isPermission)
	}
}
