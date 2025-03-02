package Orbix

import (
	"fmt"
	"goCmd/cmd/dirInfo"
	"goCmd/src"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
	"os"
	"strings"
	"time"
)

func ExecLtCommand(commandInput string) {
	user := "OneCom"

	if strings.TrimSpace(system.Location) == "" {
		system.Location = os.Getenv("CITY")
		if strings.TrimSpace(system.Location) == "" {
			system.Location = string(strings.TrimSpace(os.Getenv("USERS_LOCATION")))
		}
	}

	dir, _ := os.Getwd()
	dirC := dirInfo.CmdDir(dir)

	src.PrintPromptInfoWithoutGit(&system.Location, &user, &dirC, &commandInput)

	commandLine,
		command,
		commandArgs,
		commandLower := src.ReadCommandLine(commandInput) // Refactored input handling
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
	PermissionDenied := false
	sessionData := system.AppState{}
	session := system.Session{Path: dir, PreviousPath: dir, User: user, IsAdmin: false, GitBranch: "main", CommandHistory: []string{}}
	system.GlobalSession = session

	execCommand := structs.ExecuteCommandFuncParams{
		Command:       command,
		CommandLower:  commandLower,
		CommandArgs:   commandArgs,
		SessionPrefix: "",
		LoopData: structs.OrbixLoopData{
			Session:      &session,
			SessionData:  &sessionData,
			IsWorking:    &isWorking,
			IsPermission: &PermissionDenied,
			Username:     user,
			CommandInput: commandInput,
		},
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
		IsComHasFlag:   &isComHasFlag,
		ExecCommand:    execCommand,
		LoopData: structs.OrbixLoopData{
			Session:      &session,
			SessionData:  &sessionData,
			IsWorking:    &isWorking,
			IsPermission: &PermissionDenied,
			Username:     user,
			CommandInput: commandInput,
		},
	}

	startTimePRCOMARGS := time.Now()
	continueLoop := ProcessCommandArgs(&processCommandParams)

	if continueLoop {
		if echoTime {
			TEXCOMARGS := fmt.Sprintf("Command executed in: %s\n", time.Since(startTimePRCOMARGS))
			fmt.Println(system.Green(TEXCOMARGS))
			return
		}
		return
	}

	if isComHasFlag && (echoTime || runOnNewThread) {
		commandLine = src.RemoveFlags(commandLine)
		commandInput = src.RemoveFlags(commandInput)
		commandLine, command, commandArgs, commandLower = src.ReadCommandLine(commandLine) // Refactored input handling
	}

	if firstCharIs && lastCharIs {
		commandLower = "print"
		commandLine = fmt.Sprintf("print %s", commandLine)
		commandLine, command, commandArgs, commandLower = src.ReadCommandLine(commandLine) // Refactored input handling
	}

	isValid := utils.ValidCommandFast(commandLower, system.CmdMap)

	if len(strings.TrimSpace(commandLower)) != len(strings.TrimSpace(commandLine)) && isValid {
		session.CommandHistory = append(session.CommandHistory, commandLine)
		system.GlobalSession.CommandHistory = session.CommandHistory
	}

	if !isValid {
		session.CommandHistory = append(session.CommandHistory, commandLine)
		system.GlobalSession.CommandHistory = session.CommandHistory

		if src.CommandFile(strings.TrimSpace(commandLower)) {
			src.FullFileName(&commandArgs)
		}

		fullCommand := append([]string{command}, commandArgs...)

		if runOnNewThread {
			go executeCommandOrbix(fullCommand, command, commandLower, dir)

			if strings.TrimSpace(commandInput) != "" {
				return
			}
		} else {
			if echoTime {
				// Запоминаем время начала
				startTime = time.Now()
				executeCommandOrbix(fullCommand, command, commandLower, dir)
				// Выводим время выполнения
				TEXCOM = fmt.Sprintf("Command executed in: %s\n", time.Since(startTime))
				fmt.Println(system.Green(TEXCOM))

				if strings.TrimSpace(commandInput) != "" {
					return
				}
			} else {
				executeCommandOrbix(fullCommand, command, commandLower, dir)

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
		SessionPrefix: "",
		LoopData: structs.OrbixLoopData{
			Session:      &session,
			SessionData:  &sessionData,
			IsWorking:    &isWorking,
			IsPermission: &PermissionDenied,
			Username:     user,
			CommandInput: commandInput,
		},
	}

	execCommandCatchErrs := structs.ExecuteCommandCatchErrs{
		EchoTime:       &echoTime,
		RunOnNewThread: &runOnNewThread,
	}

	if src.CatchSyntaxErrs(execCommandCatchErrs) {
		return
	}

	if runOnNewThread {
		go Command(&execCommand)
	} else {
		if echoTime {
			// Запоминаем время начала
			startTime = time.Now()
			Command(&execCommand)
			// Выводим время выполнения
			TEXCOM = fmt.Sprintf("Command executed in: %s\n", time.Since(startTime))
			fmt.Println(system.Green(TEXCOM))
		} else {
			Command(&execCommand)
		}
	}
}
