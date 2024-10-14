package src

import (
	"fmt"
	"goCmd/cmd/commands"
	"goCmd/cmd/dirInfo"
	ExCommUtils "goCmd/src/utils"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

var (
	Absdir, _             = filepath.Abs("")
	DirUser, _            = filepath.Abs("")
	ExecutingCommand      = false
	GlobalSession         = system.Session{}
	Location              = ""
	User                  = ""
	PreviousSessionPath   = ""
	PreviousSessionPrefix = ""
	Prompt                = "$"
	Prefix                = ""
	RunningPath           = filepath.Join(Absdir, "running.txt")
	RebootAttempts        = uint(0)
	SignalReceived        = false
)

func Orbix(commandInput string, echo bool, rebooted structs.RebootedData, SD *system.AppState) {
	defer func() {
		if r := recover(); r != nil {
			PanicText := fmt.Sprintf("Panic recovered: %v", r)
			fmt.Printf("\n%s\n", red(PanicText))

			RebootAttempts += 1

			fmt.Println(yellow("Recovering from panic"))

			log.Printf("Panic recovered: %v", r)

			var reboot = structs.RebootedData{
				Username: system.UserName,
				Recover:  r,
				Prefix:   Prefix,
			}

			Orbix(commandInput, echo, reboot, SD)
		}
	}()

	if RebootAttempts > 5 {
		system.OrbixWorking = false
		fmt.Println(red("Max retry attempts reached. Exiting..."))
		log.Println("Max retry attempts reached. Exiting...")
		return
	}

	system.OrbixWorking = true
	RestartAfterInit := false

	if strings.TrimSpace(commandInput) == "restart" {
		RestartAfterInit = true
	}

	if SD == nil {
		fmt.Println(red("Fatal: App State is nil!"))
		os.Exit(1)
	}

	if err := commands.ChangeDirectory(Absdir); err != nil {
		fmt.Println(red(err))
	}

	sessionData := SD

	// Initialize colors
	InitColors()

	if !echo && commandInput == "" {
		fmt.Println(red("You cannot enable echo with an empty Input command!"))
		return
	}

	if echo && rebooted.Username == "" && commandInput == "" {
		SystemInformation()
	}

	isWorking := true
	isPermission := true
	if commandInput != "" {
		isPermission = false
	}

	// Check if password directory is empty once and handle errors here
	isEmpty, err := isPasswordDirectoryEmpty()
	if err != nil {
		animatedPrint(fmt.Sprintf("Error checking password directory: %s\n", err.Error()), "red")
		return
	}

	var username string

	if strings.TrimSpace(rebooted.Username) != "" {
		username = strings.TrimSpace(rebooted.Username)
	} else if !isEmpty && commandInput == "" {
		dir, _ := os.Getwd()
		user := dirInfo.CmdUser(dir)

		nameUser, isSuccess := CheckUser(user, sessionData)
		if !isSuccess {
			return
		}

		username = nameUser
		if username != user {
			initializeRunningFile(username)
		}

		if user == username {
			sessionData.IsAdmin = true
			sessionData.User = user
		} else {
			sessionData.IsAdmin = false
			sessionData.User = username
		}
	}

	var prompt string
	var prefix string

	var colorsMap map[string]func(...interface{}) string

	colorsMap = system.GetColorsMap()

	system.UserName = username

	if strings.TrimSpace(Location) == "" {
		Location = os.Getenv("CITY")
		if Location == "" {
			Location = "Unknown City"
		}
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
			SignalReceived = signalReceived

			if !ExecutingCommand {
				fmt.Print(red("^C"))
				if !sessionData.IsAdmin {
					dir, _ := os.Getwd()

					dirC := dirInfo.CmdDir(dir)
					user := sessionData.User
					if user == "" {
						user = dirInfo.CmdUser(dir)
					}

					if username != "" {
						user = username
					}

					fmt.Println()
					if prompt == "" {
						gitBranch, _ := GetCurrentGitBranch()
						printPromptInfo(Location, user, dirC, green, cyan, yellow, magenta, &system.Session{Path: dir, GitBranch: gitBranch}, commandInput)
					} else {
						splitPrompt := strings.Split(prompt, ", ")
						fmt.Print(colorsMap[splitPrompt[1]](splitPrompt[0]))
					}
				} else {
					dir, _ := os.Getwd()
					if prompt == "" {
						fmt.Printf("ORB %s>%s", dir, green(commandInput))
					} else {
						splitPrompt := strings.Split(prompt, ", ")
						fmt.Print(colorsMap[splitPrompt[1]](splitPrompt[0]))
					}
				}
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

	if rebooted.Prefix != "" {
		prefix = rebooted.Prefix
	} else {
		prefix = sessionData.NewSessionData(sessionData.Path, sessionData.User, sessionData.GitBranch, sessionData.IsAdmin)
	}

	session, exists := sessionData.GetSession(prefix)
	if !exists {
		fmt.Println(red("Session does not exist!"))
		return
	}

	if session == nil {
		fmt.Println(red("Session is nil!"))
		return
	}

	Prefix = fmt.Sprintf(prefix)

	// Initialize Global Vars
	Init(session)

	session.PreviousPath = PreviousSessionPath
	fmt.Println(green(session.PreviousPath))
	if PreviousSessionPrefix != "" {
		session, _ = sessionData.GetSession(PreviousSessionPrefix)
	}

	GlobalSession = *session

	dir, _ := os.Getwd()
	system.Path = dir

	for isWorking {
		if len(session.CommandHistory) < 10 {
			Init(session)
		}

		// Check if signal was received and reset flag after handling it
		if SignalReceived {
			SignalReceived = false
			continue // Continue the loop after signal
		}

		// Redirect output based on the echo setting
		if echo {
			os.Stdout, os.Stderr = originalStdout, originalStderr
		} else {
			os.Stdout, os.Stderr = devNull, devNull
		}

		// Directory and user context setup (execute only when necessary)
		dir, _ = os.Getwd()

		if RestartAfterInit {
			SD.User = username
			SD.IsAdmin = sessionData.IsAdmin
			rebooted.Prefix = prefix
			Orbix(commandInput, echo, rebooted, SD)
			return
		}

		func(rebooted *structs.RebootedData) {
			if rebooted.Recover != nil {
				RecoverText := fmt.Sprintf("Successfully recovered from the panic: %v", rebooted.Recover)
				fmt.Printf("\n%s\n", green(RecoverText))
				rebooted.Recover = nil
			}
		}(&rebooted)

		if RebootAttempts != 0 {
			RebootAttempts = 0
		}

		if echo && session.IsAdmin {
			if prompt == "" {
				fmt.Printf("ORB %s>%s", dir, green(commandInput))
			} else {
				splitPrompt := strings.Split(prompt, ", ")
				fmt.Print(colorsMap[splitPrompt[1]](splitPrompt[0]))
			}
		}

		dirC := dirInfo.CmdDir(dir)
		user := session.User
		if user == "" {
			user = dirInfo.CmdUser(dir)
		}

		if username != "" {
			user = username
		}

		if echo && !session.IsAdmin {
			// Single user check outside repeated prompt formatting
			go func() {
				watchFile(RunningPath, user, &isWorking, &isPermission)
			}()

			if echo {
				if prompt == "" {
					printPromptInfo(Location, user, dirC, green, cyan, yellow, magenta, session, commandInput) // New helper function
				} else {
					splitPrompt := strings.Split(prompt, ", ")
					fmt.Print(colorsMap[splitPrompt[1]](splitPrompt[0]))
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

		var firstCharIs, lastCharIs bool
		for index, commandLetter := range commandLine {
			if (string(commandLetter) == string('"') || string(commandLetter) == string("'")) && index == 0 {
				firstCharIs = true
			} else if (string(commandLetter) == string('"') || string(commandLetter) == string("'")) && index == len(commandLine)-1 {
				lastCharIs = true
			}
		}

		if firstCharIs && lastCharIs {
			commandLower = "print"
			commandLine = fmt.Sprintf("print %s", commandLine)
			commandLine, command, commandArgs, commandLower = readCommandLine(commandLine) // Refactored input handling
		}

		session.Path = dir

		isValid := utils.ValidCommand(commandLower, Commands)

		if len(strings.TrimSpace(commandLower)) != len(strings.TrimSpace(commandLine)) && isValid {
			session.CommandHistory = append(session.CommandHistory, commandLine)
			GlobalSession.CommandHistory = session.CommandHistory
		}

		if !isValid {
			session.CommandHistory = append(session.CommandHistory, commandLine)
			GlobalSession.CommandHistory = session.CommandHistory

			fullCommand := append([]string{command}, commandArgs...)
			err = utils.ExternalCommand(fullCommand)
			if err != nil {
				fullPath := filepath.Join(dir, command)
				fullCommand[0] = fullPath
				err = utils.ExternalCommand(fullCommand)
				if err != nil {
					isValid = utils.ValidCommand(commandLower, AdditionalCommands)
					if !isValid {
						HandleUnknownCommandUtil(commandLower, Commands)
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
		gitBranchUpdate, err := processCommand(commandLower)
		if err != nil {
			fmt.Println(red(err.Error()))
		}

		if gitBranchUpdate {
			SetGitBranch(session)
		}

		execCommand := structs.ExecuteCommandFuncParams{
			Command:       command,
			CommandLower:  commandLower,
			CommandArgs:   commandArgs,
			Dir:           dir,
			IsWorking:     &isWorking,
			IsPermission:  &isPermission,
			Username:      username,
			SD:            sessionData,
			SessionPrefix: prefix,
		}

		if strings.TrimSpace(commandLower) == "orbix" && isWorking {
			PreviousSessionPrefix = prefix
		}

		if strings.TrimSpace(commandLower) == "neofetch" && isWorking {
			ExCommUtils.NeofetchUtil(execCommand, session, Commands)
			continue
		}

		if strings.TrimSpace(commandInput) != "" && len(os.Args) > 0 {
			fmt.Println()
		}

		ExecuteCommand(execCommand)
		ExecutingCommand = false

		if strings.TrimSpace(commandInput) != "" {
			isWorking = false
		}
	}

	// Restore original outputs
	os.Stdout, os.Stderr = originalStdout, originalStderr

	PreviousSessionPath = session.Path
	session, _ = sessionData.GetSession(PreviousSessionPrefix)

	if err = commands.ChangeDirectory(session.Path); err != nil {
		fmt.Println(red(err))
	}

	sessionData.DeleteSession(prefix)

	system.OrbixWorking = false
}
