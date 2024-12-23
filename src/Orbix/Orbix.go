package Orbix

import (
	"fmt"
	_chan "goCmd/chan"
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

	LoopData, LoadUserConfigsFn := src.OrbixUser(commandInput,
		echo,
		&rebooted,
		SD,
		ExecLtCommand)

	if !*LoopData.IsWorking {
		if LoadUserConfigsFn != nil {
			// Load User Configs
			_ = LoadUserConfigsFn(false)
		}

		return
	}

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
			LoopData.SessionData,
			rebooted,
			prefix,
			LoopData.Username,
			echo,
		)

		return
	}

	LoopData.Session = session

	if _chan.UseNewPrompt {
		LoopData.SessionData.IsAdmin = false
		LoopData.Session.IsAdmin = false
		_chan.UpdateChan("orbix__src_prompt")
	} else if _chan.UseOldPrompt {
		LoopData.SessionData.IsAdmin = true
		LoopData.Session.IsAdmin = true
		_chan.UpdateChan("orbix__src_prompt")
	}

	if _chan.EnableSecure {
		LoopData.SessionData.IsAdmin = false
		LoopData.Session.IsAdmin = false
		_chan.UpdateChan("orbix__src_prompt")
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

	src.EdgeCases(LoopData,
		rebooted,
		RecoverAndRestore)

	for *LoopData.IsWorking {
		_chan.LoopData = &LoopData

		src.OrbixPrompt(LoopData.Session,
			&prompt,
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

		session.R = commandLine

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
				LoopData.Session.GitBranch, _ = system.GetCurrentGitBranch()
			}
		}()
	}

	EndOfSessions(originalStdout, originalStderr,
		LoopData.Session,
		LoopData.SessionData,
		prefix)
}
