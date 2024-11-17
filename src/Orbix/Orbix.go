package Orbix

import (
	"fmt"
	"goCmd/src"
	"goCmd/structs"
	"goCmd/system"
	"time"
)

func Orbix(commandInput string,
	echo bool,
	rebooted structs.RebootedData,
	SD *system.AppState) {
	defer func() {
		handlePanic(commandInput,
			echo,
			SD)
	}()

	LoopData := src.OrbixUser(commandInput,
		echo,
		&rebooted,
		SD,
		ExecLtCommand)

	var prompt string
	var prefix string

	colorsMap := system.GetColorsMap()

	// Signal handling setup (outside the loop)
	src.IgnoreSiC(&commandInput, &prompt,
		&LoopData)

	session := src.InitSession(&prefix,
		rebooted,
		LoopData,
	)

	originalStdout, originalStderr := setupOutputRedirect(echo)

	if *LoopData.RestartAfterInit {
		RestartAfterInitFn(
			SD,
			LoopData.SessionData,
			rebooted,
			prefix,
			LoopData.Username,
			echo,
		)
		return
	}

	LoopData.Session = session

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
		LoopData,
		rebooted,
		RecoverAndRestore)

	for *LoopData.IsWorking {
		src.OrbixPrompt(session,
			&prompt,
			&system.UserDir,
			&LoopData.Username,
			&commandInput,
			LoopData.IsWorking,
			LoopData.IsPermission,
			&colorsMap,
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
			SessionPrefix: prefix,
			LoopData:      LoopData,
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
			IsComHasFlag:   &isComHasFlag,
			ExecCommand:    execCommand,
			LoopData:       LoopData,
		}

		startTimePRCOMARGS = time.Now()
		continueLoop = ProcessCommandArgs(&processCommandParams)

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
			&commandArgs, &command, &commandLine, &commandInput, &commandLower,
			session,
		)

		execCommand = structs.ExecuteCommandFuncParams{
			Prompt:        &prompt,
			Command:       command,
			CommandLower:  commandLower,
			CommandArgs:   commandArgs,
			CommandInput:  commandInput,
			SessionPrefix: prefix,
			LoopData:      LoopData,
		}

		err = ExecLoopCommand(
			&commandLower,
			&prefix,
			&echoTime,
			&runOnNewThread,
			&execCommand,
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
		LoopData.SessionData,
		prefix)
}
