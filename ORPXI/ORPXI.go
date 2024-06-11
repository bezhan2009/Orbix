package ORPXI

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"goCmd/cmdPress"
	"goCmd/debug"
	"goCmd/utils"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func animatedPrint(text string) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(1 * time.Millisecond)
	}
}

func CMD(commandInput string) {
	utils.SystemInformation()

	isWorking := true
	isPermission := true

	var promptText string

	// Check password directory
	isEmpty, err := isPasswordDirectoryEmpty()
	if err != nil {
		animatedPrint("Ошибка при проверке директории с паролями:" + err.Error() + "\n")
		return
	}
	username := ""
	if !isEmpty && commandInput == "" {
		dir, _ := os.Getwd()
		user := cmdPress.CmdUser(dir)
		nameuser, isSuccess := CheckUser(user)
		if !isSuccess {
			return
		}
		username = nameuser
	}

	for isWorking {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		dir, _ := os.Getwd()
		dirC := cmdPress.CmdDir(dir)
		user := cmdPress.CmdUser(dir)

		if username != "" {
			user = username
		}

		if promptText != "" {
			animatedPrint("\n" + promptText)
		} else {
			animatedPrint(fmt.Sprintf("\n┌─(%s)-[%s%s]\n", cyan("ORPXI "+user), cyan("~"), cyan(dirC)))
			animatedPrint(fmt.Sprintf("└─$ %s", green(commandInput)))
		}

		var commandLine string
		var commandParts []string
		var commandArgs []string
		var commandLower string
		var command string

		if commandInput != "" {
			isWorking = false
			isPermission = false
			commandLine = strings.TrimSpace(commandInput)
			commandParts = utils.SplitCommandLine(commandLine)
			if len(commandParts) == 0 {
				continue
			}

			command = commandParts[0]
			commandArgs = commandParts[1:]
			commandLower = strings.ToLower(command)
		} else {
			commandLine = prompt.Input("> ", autoComplete)
			commandLine = strings.TrimSpace(commandLine)
			commandParts = utils.SplitCommandLine(commandLine)

			if len(commandParts) == 0 {
				continue
			}

			command = commandParts[0]
			commandArgs = commandParts[1:]
			commandLower = strings.ToLower(command)
		}

		animatedPrint("\n")

		if commandLower == "prompt" {
			handlePromptCommand(commandArgs, &promptText)
			continue
		}

		if commandLower == "help" {
			displayHelp(commandArgs, user, dir)
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

		ExecuteCommand(commandLower, command, commandLine, dir, commands, commandArgs, &isWorking, isPermission)
	}
}

func handlePromptCommand(commandArgs []string, prompt *string) {
	if len(commandArgs) < 1 {
		animatedPrint("prompt <name_prompt>\n")
		animatedPrint("to delete prompt enter:\n")
		animatedPrint("prompt delete\n")
		return
	}

	namePrompt := commandArgs[0]

	if namePrompt != "delete" {
		namePrompt = strings.TrimSpace(namePrompt)
		*prompt = namePrompt
		animatedPrint(fmt.Sprintf("Prompt set to: %s\n", *prompt))
	} else {
		*prompt, _ = os.Getwd()
		animatedPrint(fmt.Sprintf("Prompt set to: %s\n", *prompt))
		*prompt = ""
	}
}

func displayHelp(commandArgs []string, user, dir string) {
	helpText := `
For command information, type HELP
CREATE             creates a new file
CLEAN              clears the screen
CD                 changes the current directory
COPUSOURCE         copies the source code of the file
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
	animatedPrint(helpText)
	errDebug := debug.Commands("help", true, commandArgs, user, dir)
	if errDebug != nil {
		animatedPrint(errDebug.Error() + "\n")
	}
}
