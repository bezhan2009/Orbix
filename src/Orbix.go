package src

import (
	"fmt"
	"goCmd/structs"
	"goCmd/system"
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
	RunningPath           = filepath.Join(Absdir, system.OrbixRunningUsersFileName)
	GlobalSession         = system.Session{}
	Location              = ""
	User                  = ""
	Empty                 = ""
	PreviousSessionPath   = ""
	PreviousSessionPrefix = ""
	Prompt                = ""
	Prefix                = ""
	ExecutingCommand      = false
	GitCheck              = CheckGit()
	Unauthorized          = true
	RebootAttempts        = uint(0)
	SessionsStarted       = uint(0)
)

func Orbix(commandInput string,
	echo bool,
	rebooted structs.RebootedData,
	SD *system.AppState) {
	defer func() {
		if strings.TrimSpace(commandInput) == "" {
			SaveVars()
		}

		if r := recover(); r != nil {
			RecoverFromThePanic(commandInput,
				r,
				echo,
				SD,
			)
		}
	}()

	if !usingForLT(commandInput) {
		Prompt = "_>"
		execLtCommand(commandInput)

		return
	}

	RestartAfterInit := false

	sessionData := initOrbixFn(&RestartAfterInit,
		echo,
		commandInput,
		rebooted,
		SD)

	isWorking := true
	isPermission := true
	if commandInput != "" {
		isPermission = false
	}

	username, err := defineUser(commandInput,
		rebooted,
		sessionData)
	if err != nil {
		return
	}

	var prompt string
	var prefix string

	var colorsMap map[string]func(...interface{}) string

	colorsMap = system.GetColorsMap()

	system.UserName = username

	// Signal handling setup (outside the loop)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if ignoreSI(signalChan,
			sessionData,
			prompt, commandInput, username) {
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

	session := InitSession(&prefix,
		rebooted,
		sessionData,
	)

	// Redirect output based on the echo setting
	if echo {
		os.Stdout, os.Stderr = originalStdout, originalStderr
	} else {
		os.Stdout, os.Stderr = devNull, devNull
	}

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

	var (
		execCommand          structs.ExecuteCommandFuncParams
		processCommandParams structs.ProcessCommandParams

		TEXCOMARGS string

		commandLine  string
		command      string
		commandLower string
		commandArgs  []string

		runOnNewThread  bool
		echoTime        bool
		firstCharIs     bool
		lastCharIs      bool
		isComHasFlag    bool
		continueLoop    bool
		gitBranchUpdate bool

		startTimePRCOMARGS time.Time
	)

	if len(session.CommandHistory) < 10 {
		go Init(session)
	}

	if RebootAttempts != 0 {
		RecoverAndRestore(&rebooted)
		RebootAttempts = 0
	}

	for isWorking {
		OrbixPrompt(session,
			prompt,
			system.UserDir,
			username,
			commandInput,
			isWorking,
			isPermission,
			colorsMap,
		)

		// Command processing
		commandLine, command, commandArgs, commandLower = readCommandLine(commandInput) // Refactored input handling
		if commandLine == "" {
			continue
		}

		execCommand = structs.ExecuteCommandFuncParams{
			Command:       command,
			CommandLower:  commandLower,
			CommandArgs:   commandArgs,
			Dir:           system.UserDir,
			IsWorking:     &isWorking,
			IsPermission:  &isPermission,
			Username:      username,
			SD:            sessionData,
			SessionPrefix: prefix,
			Session:       session,
			GlobalSession: &GlobalSession,
		}

		processCommandParams = structs.ProcessCommandParams{
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

		startTimePRCOMARGS = time.Now()
		continueLoop = processCommandArgs(processCommandParams)

		if continueLoop {
			if echoTime {
				TEXCOMARGS = fmt.Sprintf("Command executed in: %s\n", time.Since(startTimePRCOMARGS))
				fmt.Println(green(TEXCOMARGS))
				continue
			}
			continue
		}

		ExecCommandPromptLogic(
			&firstCharIs,
			&lastCharIs,
			&isComHasFlag,
			&echoTime,
			&runOnNewThread,
			&commandArgs, &prompt, &command, &commandLine, &commandInput, &commandLower,
			session,
		)

		// Process command
		go func() {
			gitBranchUpdate, err = processCommand(commandLower)
			if err != nil {
				fmt.Println(red(err.Error()))
				DeleteUserFromRunningFile(username)
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

	EndOfSessions(originalStdout, originalStderr,
		session,
		sessionData,
		prefix)
}
