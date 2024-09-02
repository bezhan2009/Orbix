package src

import (
	"fmt"
	"goCmd/cmd/dirInfo"
	"goCmd/system"
	"goCmd/utils"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

var (
	Absdir, _       = filepath.Abs("")
	DirUser, _      = filepath.Abs("")
	UsernameFromDir = dirInfo.CmdUser(DirUser)
)

func Orbix(commandInput string, echo bool) {
	panic("GGGG")
	// Initialize Global Vars
	Init()

	if !echo && commandInput == "" {
		fmt.Println(red("You cannot enable echo with an empty Input command!"))
		return
	}

	if echo {
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

	// Directory initialization and user check outside the loop
	username := ""
	if !isEmpty && commandInput == "" {
		dir, _ := os.Getwd()
		user := dirInfo.CmdUser(dir)
		nameuser, isSuccess := CheckUser(user)
		if !isSuccess {
			return
		}
		username = nameuser
		initializeRunningFile(username) // New helper function
	}

	// Signal handling setup (outside the loop)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-signalChan
		fmt.Println("\nReceived interrupt signal. Exiting program...")
		isWorking = false
	}()

	originalStdout, originalStderr := os.Stdout, os.Stderr
	devNull, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0666)
	defer devNull.Close()

	location := os.Getenv("CITY")
	if location == "" {
		location = "Unknown City"
	}

	for isWorking {
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
			fmt.Printf("\nORB %s>%s", dir, green(commandInput))
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
				printPromptInfo(location, printUserDir, dirC, green, cyan, yellow, magenta, commandInput) // New helper function
			}
		}

		// Command processing
		commandLine, command, commandArgs, commandLower := readCommandLine(commandInput) // Refactored input handling
		if commandLine == "" {
			continue
		}

		// Process command
		if err := processCommand(commandLower, commandArgs, dir); err != nil {
			fmt.Println(red(err.Error()))
		}
		ExecuteCommand(commandLower, command, commandLine, dir, Commands, commandArgs, &isWorking, isPermission, username)
	}

	// Restore original outputs
	os.Stdout, os.Stderr = originalStdout, originalStderr
}
