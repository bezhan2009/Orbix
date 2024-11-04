package src

import (
	"fmt"
	"goCmd/cmd/commands"
	"goCmd/cmd/dirInfo"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	Absdir, _             = filepath.Abs("")
	RunningPath           = filepath.Join(Absdir, "running.txt")
	GlobalSession         = system.Session{}
	Location              = ""
	User                  = ""
	PreviousSessionPath   = ""
	PreviousSessionPrefix = ""
	Prompt                = ""
	Prefix                = ""
	ExecutingCommand      = false
	SignalReceived        = false
	GitCheck              = CheckGit()
	Unauthorized          = true
	RebootAttempts        = uint(0)
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
	func() {
		Prompt = string(strings.TrimSpace(os.Getenv("PROMPT")))
	}()

	// Initialize colors
	InitColors()

	usingForLT := func() bool {
		if strings.TrimSpace(commandInput) != "" && strings.TrimSpace(commandInput) != "restart" {
			return false
		}

		return true
	}

	if !usingForLT() {
		user := "OneCom"

		if strings.TrimSpace(Location) == "" {
			Location = os.Getenv("CITY")
			if strings.TrimSpace(Location) == "" {
				Location = string(strings.TrimSpace(os.Getenv("USERS_LOCATION")))
			}
		}

		dir, _ := os.Getwd()
		dirC := dirInfo.CmdDir(dir)

		printPromptInfoWithoutGit(Location, user, dirC, commandInput)

		commandLine, command, commandArgs, commandLower := readCommandLine(commandInput) // Refactored input handling
		if commandLine == "" {
			return
		}

		var (
			runOnNewThread bool
			echoTime       bool
			firstCharIs    bool
			lastCharIs     bool
			isComHasFlag   bool
		)

		isWorking := true
		isPermission := false
		sessionData := system.AppState{}
		session := system.Session{Path: dir, PreviousPath: dir, User: user, IsAdmin: false, GitBranch: "main", CommandHistory: []string{}}
		GlobalSession = session

		execCommand := structs.ExecuteCommandFuncParams{
			Command:       command,
			CommandLower:  commandLower,
			CommandArgs:   commandArgs,
			Dir:           dir,
			IsWorking:     &isWorking,
			IsPermission:  &isPermission,
			Username:      user,
			SD:            &sessionData,
			SessionPrefix: "",
			Session:       &session,
			GlobalSession: &GlobalSession,
		}

		processCommandParams := structs.ProcessCommandParams{
			Command:        command,
			CommandInput:   commandInput,
			CommandLower:   commandLower,
			CommandLine:    commandLine,
			CommandArgs:    commandArgs,
			RunOnNewThread: &runOnNewThread,
			EchoTime:       &echoTime,
			FirstCharIs:    &firstCharIs,
			LastCharIs:     &lastCharIs,
			IsWorking:      &isWorking,
			IsComHasFlag:   &isComHasFlag,
			Session:        &session,
			ExecCommand:    execCommand,
		}

		startTimePRCOMARGS := time.Now()
		continueLoop := processCommandArgs(processCommandParams)

		if continueLoop {
			if echoTime {
				TEXCOMARGS := fmt.Sprintf("Command executed in: %s\n", time.Since(startTimePRCOMARGS))
				fmt.Println(green(TEXCOMARGS))
				return
			}
			return
		}

		if isComHasFlag && (echoTime || runOnNewThread) {
			commandLine = removeFlags(commandLine)
			commandInput = removeFlags(commandInput)
			commandLine, command, commandArgs, commandLower = readCommandLine(commandLine) // Refactored input handling
		}

		if firstCharIs && lastCharIs {
			commandLower = "print"
			commandLine = fmt.Sprintf("print %s", commandLine)
			commandLine, command, commandArgs, commandLower = readCommandLine(commandLine) // Refactored input handling
		}

		isValid := utils.ValidCommand(commandLower, Commands)

		if len(strings.TrimSpace(commandLower)) != len(strings.TrimSpace(commandLine)) && isValid {
			session.CommandHistory = append(session.CommandHistory, commandLine)
			GlobalSession.CommandHistory = session.CommandHistory
		}

		if !isValid {
			session.CommandHistory = append(session.CommandHistory, commandLine)
			GlobalSession.CommandHistory = session.CommandHistory

			if commandFile(strings.TrimSpace(commandLower)) {
				fullFileName(&commandArgs)
			}

			fullCommand := append([]string{command}, commandArgs...)

			// Логика выполнения команды, которую можно запускать в новом потоке
			executeCommandOrbix := func() {
				err := utils.ExternalCommand(fullCommand)
				if err != nil {
					fullPath := filepath.Join(dir, command)
					fullCommand[0] = fullPath
					err = utils.ExternalCommand(fullCommand)
					if err != nil {
						isValid = utils.ValidCommand(commandLower, AdditionalCommands)
						if !isValid {
							HandleUnknownCommandUtil(commandLower, Commands)
							return
						}
					}
				}
			}

			if runOnNewThread {
				go executeCommandOrbix()

				if strings.TrimSpace(commandInput) != "" {
					return
				}
			} else {
				if echoTime {
					// Запоминаем время начала
					startTime := time.Now()
					executeCommandOrbix()
					// Выводим время выполнения
					TEXCOM := fmt.Sprintf("Command executed in: %s\n", time.Since(startTime))
					fmt.Println(green(TEXCOM))

					if strings.TrimSpace(commandInput) != "" {
						return
					}
				} else {
					executeCommandOrbix()

					if strings.TrimSpace(commandInput) != "" {
						return
					}
				}
			}

			if strings.TrimSpace(commandInput) != "" {
				return
			}
			return
		}

		execCommand = structs.ExecuteCommandFuncParams{
			Command:       command,
			CommandLower:  commandLower,
			CommandArgs:   commandArgs,
			CommandInput:  commandInput,
			Dir:           dir,
			IsWorking:     &isWorking,
			IsPermission:  &isPermission,
			Username:      "OneCom",
			SD:            &sessionData,
			SessionPrefix: "",
			Session:       &session,
			GlobalSession: &GlobalSession,
		}

		execCommandCatchErrs := structs.ExecuteCommandCatchErrs{
			EchoTime:       &echoTime,
			RunOnNewThread: &runOnNewThread,
		}

		if catchSyntaxErrs(execCommandCatchErrs) {
			return
		}

		if runOnNewThread {
			go ExecuteCommand(execCommand)
		} else {
			if echoTime {
				// Запоминаем время начала
				startTime := time.Now()
				ExecuteCommand(execCommand)
				// Выводим время выполнения
				TEXCOM := fmt.Sprintf("Command executed in: %s\n", time.Since(startTime))
				fmt.Println(green(TEXCOM))
			} else {
				ExecuteCommand(execCommand)
			}
		}

		return
	}

	if strings.TrimSpace(strings.ToLower(system.OperationSystem)) == "windows" {
		Commands = append(Commands, structs.Command{Name: "neofetch", Description: "Displays information about the system"})
		AdditionalCommands = append(AdditionalCommands, structs.Command{Name: "neofetch", Description: "Displays information about the system"})
	}

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
		Unauthorized = false
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
		if strings.TrimSpace(Location) == "" {
			Location = string(strings.TrimSpace(os.Getenv("USERS_LOCATION")))
		}
	}

	// Signal handling setup (outside the loop)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	var signalReceived bool

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			sig := <-signalChan
			signalReceived = true
			SignalReceived = signalReceived

			if sig == syscall.SIGHUP {
				RemoveUserFromRunningFile(system.UserName)
				os.Exit(1)
			}

			if !ExecutingCommand {
				fmt.Println(red("^C"))
				if !GlobalSession.IsAdmin {
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
						if GitCheck {
							gitBranch, _ := GetCurrentGitBranch()
							printPromptInfo(Location, user, dirC, commandInput, &system.Session{Path: dir, GitBranch: gitBranch})
						} else {
							printPromptInfoWithoutGit(Location, user, dirC, commandInput)
						}
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
	go Init(session)

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
			go Init(session)
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
			if len(os.Args) > 1 {
				return
			}

			Orbix("", echo, rebooted, SD)
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
			if !Unauthorized {
				go func() {
					watchFile(RunningPath, user, &isWorking, &isPermission)
				}()
			}

			if echo {
				if prompt == "" {
					if GitCheck {
						printPromptInfo(Location, user, dirC, commandInput, session) // New helper function
					} else {
						printPromptInfoWithoutGit(Location, user, dirC, commandInput) // New helper function
					}
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

		var (
			runOnNewThread bool
			echoTime       bool
			firstCharIs    bool
			lastCharIs     bool
			isComHasFlag   bool
		)

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
			Session:       session,
			GlobalSession: &GlobalSession,
		}

		processCommandParams := structs.ProcessCommandParams{
			Command:        command,
			CommandInput:   commandInput,
			CommandLower:   commandLower,
			CommandLine:    commandLine,
			CommandArgs:    commandArgs,
			RunOnNewThread: &runOnNewThread,
			EchoTime:       &echoTime,
			FirstCharIs:    &firstCharIs,
			LastCharIs:     &lastCharIs,
			IsWorking:      &isWorking,
			IsComHasFlag:   &isComHasFlag,
			Session:        session,
			ExecCommand:    execCommand,
		}

		startTimePRCOMARGS := time.Now()
		continueLoop := processCommandArgs(processCommandParams)

		if continueLoop {
			if echoTime {
				TEXCOMARGS := fmt.Sprintf("Command executed in: %s\n", time.Since(startTimePRCOMARGS))
				fmt.Println(green(TEXCOMARGS))
				continue
			}
			continue
		}

		if isComHasFlag && (echoTime || runOnNewThread) {
			commandLine = removeFlags(commandLine)
			commandInput = removeFlags(commandInput)
			commandLine, command, commandArgs, commandLower = readCommandLine(commandLine) // Refactored input handling
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

			if commandFile(strings.TrimSpace(commandLower)) {
				fullFileName(&commandArgs)
			}

			fullCommand := append([]string{command}, commandArgs...)

			// Логика выполнения команды, которую можно запускать в новом потоке
			executeCommandOrbix := func() {
				err = utils.ExternalCommand(fullCommand)
				if err != nil {
					fullPath := filepath.Join(dir, command)
					fullCommand[0] = fullPath
					err = utils.ExternalCommand(fullCommand)
					if err != nil {
						isValid = utils.ValidCommand(commandLower, AdditionalCommands)
						if !isValid {
							HandleUnknownCommandUtil(commandLower, Commands)
							return
						}
					}
				}
			}

			if runOnNewThread {
				go executeCommandOrbix()

				if strings.TrimSpace(commandInput) != "" {
					return
				}
			} else {
				if echoTime {
					// Запоминаем время начала
					startTime := time.Now()
					executeCommandOrbix()
					// Выводим время выполнения
					TEXCOM := fmt.Sprintf("Command executed in: %s\n", time.Since(startTime))
					fmt.Println(green(TEXCOM))

					if strings.TrimSpace(commandInput) != "" {
						return
					}
				} else {
					executeCommandOrbix()

					if strings.TrimSpace(commandInput) != "" {
						return
					}
				}
			}

			if strings.TrimSpace(commandInput) != "" {
				return
			}
			continue
		}

		if strings.TrimSpace(commandLower) == "prompt" {
			handlePromptCommand(commandArgs, &prompt)
			continue
		}

		// Process command
		go func() {
			gitBranchUpdate, err := processCommand(commandLower)
			if err != nil {
				fmt.Println(red(err.Error()))
				RemoveUserFromRunningFile(username)
				return
			}

			if gitBranchUpdate {
				session.GitBranch, err = GetCurrentGitBranch()
				if err != nil {
					fmt.Println("Error Updating Git Branch", red(err.Error()))
				}
			}
		}()

		execCommand = structs.ExecuteCommandFuncParams{
			Command:       command,
			CommandLower:  commandLower,
			CommandArgs:   commandArgs,
			CommandInput:  commandInput,
			Dir:           dir,
			IsWorking:     &isWorking,
			IsPermission:  &isPermission,
			Username:      username,
			SD:            sessionData,
			SessionPrefix: prefix,
			Session:       session,
			GlobalSession: &GlobalSession,
		}

		execCommandCatchErrs := structs.ExecuteCommandCatchErrs{
			EchoTime:       &echoTime,
			RunOnNewThread: &runOnNewThread,
		}

		if strings.TrimSpace(commandLower) == "orbix" && isWorking {
			PreviousSessionPrefix = prefix
		}

		if catchSyntaxErrs(execCommandCatchErrs) {
			continue
		}

		if runOnNewThread {
			go ExecuteCommand(execCommand)
		} else {
			if echoTime {
				// Запоминаем время начала
				startTime := time.Now()
				ExecuteCommand(execCommand)
				// Выводим время выполнения
				TEXCOM := fmt.Sprintf("Command executed in: %s\n", time.Since(startTime))
				fmt.Println(green(TEXCOM))
			} else {
				ExecuteCommand(execCommand)
			}
		}
	}

	// Restore original outputs
	os.Stdout, os.Stderr = originalStdout, originalStderr

	PreviousSessionPath = session.Path
	session, _ = sessionData.GetSession(PreviousSessionPrefix)

	if err = commands.ChangeDirectory(session.Path); err != nil {
		fmt.Println(red("Error changing directory:", err))
	}

	sessionData.DeleteSession(prefix)

	system.OrbixWorking = false
}
