package Orbix

import (
	"fmt"
	"goCmd/src"
	"goCmd/structs"
	"goCmd/system"
	"os"
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

	OrbixLoopData := src.OrbixUser(commandInput,
		echo,
		&rebooted,
		SD,
		ExecLtCommand)

	var prompt string
	var prefix string

	var colorsMap map[string]func(...interface{}) string

	colorsMap = system.GetColorsMap()

	// Signal handling setup (outside the loop)
	src.IgnoreSiC(commandInput, prompt,
		OrbixLoopData)

	originalStdout, originalStderr := os.Stdout, os.Stderr
	devNull, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0666)
	defer func() {
		err := devNull.Close()
		if err != nil {
			return
		}
	}()

	session := src.InitSession(&prefix,
		rebooted,
		OrbixLoopData,
	)

	// Redirect output based on the echo setting
	if echo {
		os.Stdout, os.Stderr = originalStdout, originalStderr
	} else {
		os.Stdout, os.Stderr = devNull, devNull
	}

	if OrbixLoopData.RestartAfterInit {
		RestartAfterInitFn(
			SD,
			OrbixLoopData.SessionData,
			rebooted,
			prefix,
			OrbixLoopData.Username,
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

		err error

		startTimePRCOMARGS time.Time
	)

	src.EdgeCases(session,
		OrbixLoopData,
		rebooted,
		RecoverAndRestore)

	for OrbixLoopData.IsWorking {
		src.OrbixPrompt(session,
			prompt,
			system.UserDir,
			OrbixLoopData.Username,
			commandInput,
			OrbixLoopData.IsWorking,
			OrbixLoopData.IsPermission,
			colorsMap,
		)

		// Command processing
		commandLine, command, commandArgs, commandLower = src.ReadCommandLine(commandInput) // Refactored input handling
		if commandLine == "" {
			continue
		}

		execCommand = structs.ExecuteCommandFuncParams{
			Prompt:        &prompt,
			Command:       command,
			CommandLower:  commandLower,
			CommandArgs:   commandArgs,
			IsWorking:     &OrbixLoopData.IsWorking,
			IsPermission:  &OrbixLoopData.IsPermission,
			Username:      OrbixLoopData.Username,
			SD:            OrbixLoopData.SessionData,
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
			IsWorking:      &OrbixLoopData.IsWorking,
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

		if ExecCommandPromptLogic(
			&firstCharIs,
			&lastCharIs,
			&isComHasFlag,
			&echoTime,
			&runOnNewThread,
			&commandArgs, &command, &commandLine, &commandInput, &commandLower,
			session,
		) {
			updateGlobalCommVars()
			continue
		}

		execCommand = structs.ExecuteCommandFuncParams{
			Prompt:        &prompt,
			Command:       command,
			CommandLower:  commandLower,
			CommandArgs:   commandArgs,
			CommandInput:  commandInput,
			IsWorking:     &OrbixLoopData.IsWorking,
			IsPermission:  &OrbixLoopData.IsPermission,
			Username:      OrbixLoopData.Username,
			SD:            OrbixLoopData.SessionData,
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

		updateGlobalCommVars()
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
		OrbixLoopData.SessionData,
		prefix)
}
