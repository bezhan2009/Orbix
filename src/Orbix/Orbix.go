package Orbix

import (
	"fmt"
	_chan "goCmd/chan"
	"goCmd/src"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
	"strings"
	"time"
)

func Orbix(commandInput string,
	echo bool,
	rebooted structs.RebootedData,
	SD *system.AppState) {
	LoadUserConfigsFn := func(echo bool) error { return nil }

	var LoopData structs.OrbixLoopData
	originalStdout, originalStderr := setupOutputRedirect(echo)

	if rebooted.LoopData != LoopData {
		LoopData, LoadUserConfigsFn = rebooted.LoopData, rebooted.LoadUserConfigsFn
	} else {
		// Initialize colors
		system.InitColors()

		LoopData, LoadUserConfigsFn = src.OrbixUser(commandInput,
			echo,
			&rebooted,
			SD,
			ExecLtCommand)
	}

	defer func() {
		if system.Debug {
			return
		}

		r := recover()
		handlePanic(commandInput,
			echo,
			SD,
			LoopData,
			LoadUserConfigsFn,
			r)
	}()

	if !*LoopData.IsWorking {
		_ = LoadUserConfigsFn(false)
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

		commandLine         string
		command             string
		commandLower        string
		commandArgs         []string
		SplittedCommandLine []string

		runOnNewThread  bool
		echoTime        bool
		firstCharIs     bool
		lastCharIs      bool
		isComHasFlag    bool
		gitBranchUpdate bool

		err error

		startTimePRCOMARGS time.Time
	)

	_chan.SetVarFn = getCustomVar

	src.EdgeCases(&LoopData,
		session,
		rebooted,
		RecoverAndRestore)

	LoopCommand := func(commandLine, command, commandLower string,
		commandArgs []string) (continueLoop bool) {

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

			return true
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
			return true
		}

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
			return true
		}

		// Process command
		go func() {
			updGitBranch := func() {
				gitBranchUpdate = src.ProcessCommand(commandLower)

				if gitBranchUpdate {
					LoopData.Session.GitBranch, _ = system.GetCurrentGitBranch()
				}
			}

			cleanHistoryDuplicates := func() {
				if system.GlobalSession.CommandHistory != nil {
					utils.CleanSliceDuplicates(system.GlobalSession.CommandHistory)
				}

				if session.CommandHistory != nil {
					utils.CleanSliceDuplicates(session.CommandHistory)
				}
			}

			updGitBranch()
			for {
				time.Sleep(10 * time.Second)

				updGitBranch()
				cleanHistoryDuplicates()
			}
		}()

		return false
	}

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

		SplittedCommandLine = strings.Split(commandLine, "&&")
		if len(SplittedCommandLine) > 1 && strings.TrimSpace(commandLine) != "&&" {
			session.CommandHistory = append(session.CommandHistory, commandLine)

			for _, commandLineLoop := range SplittedCommandLine {
				commandLine, command, commandArgs, commandLower = src.ReadCommandLine(commandLineLoop)
				LoopCommand(commandLine, command, commandLower, commandArgs)
			}

			continue
		}

		if LoopCommand(commandLine, command, commandLower, commandArgs) {
			continue
		}
	}

	EndOfSessions(originalStdout, originalStderr,
		LoopData.Session,
		LoopData.SessionData,
		prefix)
}
