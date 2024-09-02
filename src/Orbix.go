package src

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"goCmd/cmd/dirInfo"
	"goCmd/system"
	"goCmd/utils"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	Absdir, _       = filepath.Abs("")
	DirUser, _      = filepath.Abs("")
	UsernameFromDir = dirInfo.CmdUser(DirUser)
)

func Orbix(commandInput string, echo bool) {
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
