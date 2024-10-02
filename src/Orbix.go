package src

import (
	"fmt"
	"goCmd/cmd/dirInfo"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

var (
	Absdir, _       = filepath.Abs("")
	DirUser, _      = filepath.Abs("")
	UsernameFromDir = dirInfo.CmdUser(DirUser)
)

func Orbix(commandInput string, echo bool, rebooted structs.RebootedData) {
	// Initialize Global Vars
	Init()

	if !echo && commandInput == "" {
		fmt.Println(red("You cannot enable echo with an empty Input command!"))
		return
	}

	if echo && rebooted.Username == "" && commandInput == "" {
		utils.SystemInformation()
	}

	isWorking := true
	isPermission := true
	if commandInput != "" {
		isPermission = false
	}

	// Check if password directory is empty once and handle errors here
	isEmpty, err := isPasswordDirectoryEmpty()
	if err != nil {
		animatedPrint("Error checking password directory: " + err.Error() + "\n")
		return
	}

	username := ""

	if strings.TrimSpace(rebooted.Username) != "" {
		username = strings.TrimSpace(rebooted.Username)
	} else if !isEmpty && commandInput == "" {
		dir, _ := os.Getwd()
		user := dirInfo.CmdUser(dir)
		nameUser, isSuccess := CheckUser(user)
		if !isSuccess {
			return
		}

		username = nameUser
		initializeRunningFile(username) // New helper function
	}

	location := os.Getenv("CITY")
	if location == "" {
		location = "Unknown City"
	}

	// Signal handling setup (outside the loop)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var signalReceived bool

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			<-signalChan
			signalReceived = true
			fmt.Print(red("^C"))
			if !system.IsAdmin {
				dir, _ := os.Getwd()

				dirC := dirInfo.CmdDir(dir)
				user := system.User
				if user == "" {
					user = dirInfo.CmdUser(dir)
				}

				var printUserDir string

				if username != "" {
					user = username
					printUserDir = user
				}
				fmt.Println()
				printPromptInfo(location, printUserDir, dirC, green, cyan, yellow, magenta, commandInput)
			} else {
				dir, _ := os.Getwd()
				fmt.Println()
				fmt.Printf("\nORB %s>%s", dir, green(commandInput))
			}
		}
	}()

	originalStdout, originalStderr := os.Stdout, os.Stderr
	devNull, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0666)
	defer func() {
		err = devNull.Close()
		if err != nil {
			return
		}
	}()

	var prompt string

	for isWorking {
		// Check if signal was received and reset flag after handling it
		if signalReceived {
			signalReceived = false
			continue // Continue the loop after signal
		}

		// Redirect output based on the echo setting
		if echo {
			os.Stdout, os.Stderr = originalStdout, originalStderr
		} else {
			os.Stdout, os.Stderr = devNull, devNull
		}

		// Directory and user context setup (execute only when necessary)
		dir, _ := os.Getwd()
		printUserDir := UsernameFromDir // Use cached username for printing

		if echo && system.IsAdmin {
			if prompt == "" {
				fmt.Printf("ORB %s>%s", dir, green(commandInput))
			} else {
				fmt.Print(green(prompt))
			}
		}

		if !system.IsAdmin {
			dirC := dirInfo.CmdDir(dir)
			user := system.User
			if user == "" {
				user = dirInfo.CmdUser(dir)
			}

			if username != "" {
				user = username
				printUserDir = user
			}

			// Single user check outside repeated prompt formatting
			if !checkUserInRunningFile(username) {
				fmt.Println(red("User not authorized."))
				isWorking = false
				isPermission = false
				continue
			}

			if echo {
				if prompt == "" {
					printPromptInfo(location, printUserDir, dirC, green, cyan, yellow, magenta, commandInput) // New helper function
				} else {
					fmt.Print(green(prompt))
				}
			}
		}

		// Command processing
		commandLine, command, commandArgs, commandLower := readCommandLine(commandInput) // Refactored input handling
		if commandLine == "" {
			continue
		}

		if commandInt, err := strconv.Atoi(command); err == nil && len(commandArgs) == 0 {
			fmt.Println(magenta(commandInt))
			continue
		}

		CommandHistory = append(CommandHistory, commandLine)

		isValid := utils.ValidCommand(commandLower, Commands)

		if !isValid {
			fullCommand := append([]string{command}, commandArgs...)
			err = utils.ExternalCommand(fullCommand)
			if err != nil {
				fullPath := filepath.Join(dir, command)
				fullCommand[0] = fullPath
				err = utils.ExternalCommand(fullCommand)
				if err != nil {
					isValid = utils.ValidCommand(commandLower, AdditionalCommands)
					if !isValid {
						HandleUnknownCommandUtil(commandLower, commandLine, Commands)
						continue
					}
				}
			}
			continue
		}

		if strings.TrimSpace(commandLower) == "prompt" {
			handlePromptCommand(commandArgs, &prompt)
			continue
		}

		// Process command
		if err := processCommand(commandLower, commandArgs); err != nil {
			fmt.Println(red(err.Error()))
			isWorking = false
		}

		execCommand := structs.ExecuteCommandFuncParams{
			Command:      command,
			CommandLower: commandLower,
			CommandArgs:  commandArgs,
			Dir:          dir,
			IsWorking:    &isWorking,
			IsPermission: isPermission,
			Username:     username,
		}

		ExecuteCommand(execCommand)
	}

	// Restore original outputs
	os.Stdout, os.Stderr = originalStdout, originalStderr
}
