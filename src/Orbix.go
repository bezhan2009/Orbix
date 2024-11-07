package src

import (
	"fmt"
	"goCmd/cmd/dirInfo"
	"goCmd/structs"
	"goCmd/system"
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
	SessionsStarted       = uint(0)
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

	if !usingForLT(commandInput) {
		execLtCommand(commandInput)

		return
	}

	RestartAfterInit := false

	sessionData := initOrbixFn(&RestartAfterInit, echo, commandInput, rebooted, SD)

	isWorking := true
	isPermission := true
	if commandInput != "" {
		isPermission = false
	}

	username, err := defineUser(commandInput, rebooted, sessionData)
	if err != nil {
		return
	}

	var prompt string
	var prefix string

	var colorsMap map[string]func(...interface{}) string

	colorsMap = system.GetColorsMap()

	system.UserName = username

	setLocation()

	// Signal handling setup (outside the loop)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	var signalReceived bool

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if ignoreSI(signalChan, &signalReceived, sessionData, prompt, commandInput, username) {
			return
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
			restartAfterInit(
				SD,
				sessionData,
				rebooted,
				prefix,
				username,
				echo,
			)
			return
		}

		RecoverAndRestore(&rebooted)

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

		ExecCommandPromptLogic(
			firstCharIs,
			lastCharIs,
			isComHasFlag,
			echoTime,
			runOnNewThread,
			dir,
			&commandArgs,
			&prompt,
			&command,
			&commandLine,
			&commandInput,
			&commandLower,
			session,
		)

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

		err = ExecLoopCommand(
			commandLower,
			prefix,
			echoTime,
			runOnNewThread,
			execCommand,
		)

		UnknownCommandsCounter = 0

		if err != nil {
			continue
		}
	}

	EndOfSessions(originalStdout, originalStderr, session, sessionData, prefix)
}
