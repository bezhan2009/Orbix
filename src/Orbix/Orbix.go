package Orbix

import (
	"fmt"
	"goCmd/src"
	"goCmd/structs"
	"goCmd/system"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Orbix(commandInput string,
	echo bool,
	rebooted structs.RebootedData,
	SD *system.AppState) {
	defer func() {
		if r := recover(); r != nil {
			RecoverFromThePanic(commandInput,
				r,
				echo,
				SD,
			)
		}
	}()

	if !UsingForLT(commandInput) {
		system.Prompt = "_>"
		ExecLtCommand(commandInput)

		return
	}

	RestartAfterInit := false

	sessionData := src.InitOrbixFn(&RestartAfterInit,
		echo,
		commandInput,
		rebooted,
		SD)

	isWorking := true
	isPermission := true
	if commandInput != "" {
		isPermission = false
	}

	username, err := src.DefineUser(commandInput,
		rebooted,
		sessionData)
	if err != nil {
		return
	}

	// Load User Configs
	src.LoadConfigs()
	if username != "" {
		system.EditableVars["user"] = &username
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
		if src.IgnoreSI(signalChan,
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

	session := src.InitSession(&prefix,
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
		RestartAfterInitFn(
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
		go system.InitSession(session)
	}

	if system.RebootAttempts != 0 {
		RecoverAndRestore(&rebooted)
		system.RebootAttempts = 0
	}

	for isWorking {
		src.OrbixPrompt(session,
			prompt,
			system.UserDir,
			username,
			commandInput,
			isWorking,
			isPermission,
			colorsMap,
		)

		// Command processing
		commandLine, command, commandArgs, commandLower = src.ReadCommandLine(commandInput) // Refactored input handling
		if commandLine == "" {
			continue
		}

		execCommand = structs.ExecuteCommandFuncParams{
			Command:       command,
			CommandLower:  commandLower,
			CommandArgs:   commandArgs,
			IsWorking:     &isWorking,
			IsPermission:  &isPermission,
			Username:      username,
			SD:            sessionData,
			SessionPrefix: prefix,
			Session:       session,
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
		continueLoop = ProcessCommandArgs(processCommandParams)

		if continueLoop {
			if echoTime {
				TEXCOMARGS = fmt.Sprintf("Command executed in: %s\n", time.Since(startTimePRCOMARGS))
				fmt.Println(system.Green(TEXCOMARGS))
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
		}

		err = ExecLoopCommand(
			commandLower,
			prefix,
			echoTime,
			runOnNewThread,
			execCommand,
		)

		src.UnknownCommandsCounter = 0

		if err != nil {
			continue
		}

		// Process command
		go func() {
			gitBranchUpdate = src.ProcessCommand(commandLower)

			if gitBranchUpdate {
				session.GitBranch, _ = system.GetCurrentGitBranch()
			}
		}()
	}

	EndOfSessions(originalStdout, originalStderr,
		session,
		sessionData,
		prefix)
}
