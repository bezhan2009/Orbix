package Orbix

import (
	"errors"
	"fmt"
	ExOpenedCommReadUtils "goCmd/cmd/commands/Read/utils"
	"goCmd/cmd/commands/fcommands"
	"goCmd/src"
	"goCmd/system"
	"goCmd/system/errs"
	"goCmd/utils"
	"strings"
	"time"
)

func updateGlobalCommVars() {
	TEXCOM = ""
	commandExLogg = ""

	commExLoggSplit = []string{}
	commandArgsListStr = []string{}
	commandList = []string{}
	fullCommandArgs = []string{}
	argList = []string{}
	argSplit = []string{}

	iArgSplit = 0

	startTime = time.Time{}
}

func SetCommandVarValues(commandArgs *[]string, autocomplete bool) ([]string, string) {
	for iArg, arg := range *commandArgs {
		if !strings.Contains(arg, "$") {
			continue
		}

		if strings.Contains(arg, "+") {
			argSplit = strings.Split(arg, "+")

			var argSplitTemp []string
			for iArgSplit, arg = range argSplit {
				if string(arg[0]) == "$" && len(strings.TrimSpace(arg)) > 1 {
					if string(arg[len(arg)-1]) == "-" {
						argList = strings.Split(arg, "-")
						arg = argList[0]
						fullCommandArgs[iArg] = arg
						*commandArgs = fullCommandArgs
						continue
					}

					customVar, err := getCustomVar(arg[1:])
					if err != nil && !autocomplete {
						fmt.Println(system.Red(err))
						continue
					}

					if customVar != nil {
						argSplitTemp = append(argSplitTemp, customVar.(string))
					}
				}
			}

			sumArgs = ""
			for _, argCon := range argSplitTemp {
				sumArgs += argCon
			}

			fullCommandArgs = *commandArgs
			fullCommandArgs[iArg] = sumArgs
			continue
		}

		if string(arg[0]) == "$" && len(strings.TrimSpace(arg)) > 1 {
			if string(arg[len(arg)-1]) == "-" {
				argList = strings.Split(arg, "-")
				arg = argList[0]
				fullCommandArgs[iArg] = arg
				*commandArgs = fullCommandArgs
				continue
			}

			customVar, err := getCustomVar(arg[1:])
			if err != nil && !autocomplete {
				fmt.Println(system.Red(err))
			}

			if customVar != nil {
				fullCommandArgs[iArg] = customVar.(string)
				*commandArgs = fullCommandArgs
			}
		}
	}

	return fullCommandArgs, sumArgs
}

// replaceShortcuts заменяет строки в массиве на значения из карты shortcuts.
func replaceShortcuts(strings []string, shortcuts map[string]string) []string {
	for i, str := range strings {
		if shortcut, exists := shortcuts[str]; exists {
			strings[i] = shortcut // Замена строки на соответствующий shortcut
		}
	}
	return strings
}

func ExecCommandPromptLogic(
	firstCharIs,
	lastCharIs,
	isComHasFlag,
	echoTime,
	runOnNewThread *bool,
	commandArgs *[]string,
	command, commandLine, commandInput, commandLower *string,
	session *system.Session) bool {
	defer updateGlobalCommVars()

	res := replaceShortcuts(strings.Fields(*commandLine), system.Shortcuts)

	shCmdLine := strings.Join(res, " ")

	*commandLine, *command, *commandArgs, *commandLower = src.ReadCommandLine(shCmdLine)

	if strings.ToLower(strings.TrimSpace(*commandLine)) == "r" && strings.TrimSpace(session.R) != "" {
		fmt.Println(session.R)

		*commandLine, *command, *commandArgs, *commandLower = src.ReadCommandLine(session.R)
	}

	if *commandLine == "cd.." {
		*commandLine, *command, *commandArgs, *commandLower = src.ReadCommandLine("cd ..")
		return false
	}

	if *isComHasFlag && (*echoTime || *runOnNewThread) {
		*commandLine = src.RemoveFlags(*commandLine)
		*commandInput = src.RemoveFlags(*commandInput)
		*commandLine, *command, *commandArgs, *commandLower = src.ReadCommandLine(*commandLine) // Refactored input handling
	}

	if *firstCharIs && *lastCharIs {
		*commandLower = "print"
		*commandLine = fmt.Sprintf("print %s", *commandLine)
		*commandLine, *command, *commandArgs, *commandLower = src.ReadCommandLine(*commandLine) // Refactored input handling
	}

	commandArgsListStr = *commandArgs
	commandExLogg = *command

	for _ = range *commandLine {
		if !strings.Contains(commandExLogg, "$") {
			break
		}

		if strings.Contains(commandExLogg, "+") {
			commExLoggSplit = strings.Split(commandExLogg, "+")

			commandExLogg = ""

			for _, arg := range commExLoggSplit {
				if string(arg[0]) == "$" && len(strings.TrimSpace(arg)) > 1 {
					if string(arg[len(arg)-1]) == "-" {
						argList = strings.Split(arg, "-")
						commandExLogg = argList[0]
						continue
					}

					customVar, err := getCustomVar(arg[1:])
					if err != nil {
						fmt.Println(system.Red(err))
						continue
					}

					if customVar != nil {
						commandExLogg += customVar.(string)
					}
				}
			}

			*commandLine, *command, *commandArgs, *commandLower = src.ReadCommandLine(commandExLogg) // Refactored input handling
			fullCommandArgs = append(fullCommandArgs, *commandArgs...)

			break
		}

		if string(commandExLogg[0]) == "$" && len(strings.TrimSpace(commandExLogg)) > 1 {
			if string(commandExLogg[len(commandExLogg)-1]) == "-" {
				commandList = strings.Split(commandExLogg, "-")
				commandExLogg = commandList[0]
				break
			}

			customVar, err := getCustomVar(commandExLogg[1:])
			if err != nil {
				fmt.Println(system.Red(err))
				continue
			}

			if customVar != nil {
				commandExLogg = customVar.(string)
				*commandLine, *command, *commandArgs, *commandLower = src.ReadCommandLine(commandExLogg) // Refactored input handling
				fullCommandArgs = append(fullCommandArgs, *commandArgs...)
			}
		}
	}

	fullCommandArgs = append(fullCommandArgs, commandArgsListStr...)
	*commandArgs = fullCommandArgs

	SetCommandVarValues(commandArgs, false)

	InQuotes := true

	for iArg := 0; iArg < len(fullCommandArgs); iArg++ {
		arg := fullCommandArgs[iArg]

		if strings.HasPrefix(arg, "(") && InQuotes {

			// Ищем конец команды в скобках
			endIdx := -1
			cntSpaces := 0
			for j := iArg; j < len(fullCommandArgs); j++ {
				if fullCommandArgs[j] == " " || fullCommandArgs[j] == "" {
					cntSpaces++
				}

				if strings.HasSuffix(fullCommandArgs[j], ")") {
					endIdx = j - cntSpaces
					break
				}
			}

			if endIdx != 2 && endIdx != -1 {
				fmt.Println(system.Red(fmt.Sprintf("Syntax Error in '%s'\nNot enough arguments to call func", fullCommandArgs[iArg])))
				return true
			}

			// Если нет закрывающей скобки
			if endIdx == -1 {
				fmt.Println(system.Red(fmt.Sprintf("Syntax Error in '%s'\nPlease Close the '('", fullCommandArgs[iArg])))
				fmt.Println(fmt.Sprintf(" %s\n%s", fmt.Sprintf("%s%s", system.Red(fmt.Sprintf("%s %s", fullCommandArgs[iArg], fullCommandArgs[iArg+1])), system.Yellow(")")), strings.Repeat(system.Red("─"), (len(fmt.Sprintf("%s %s", fullCommandArgs[iArg], fullCommandArgs[iArg+1]))-1)+2)+system.Yellow("ꜛ")))
				return true
			}

			// Объединяем команду внутри скобок
			innerCommand := strings.Join(fullCommandArgs[iArg:endIdx+1], " ")
			innerCommand = strings.TrimSuffix(strings.TrimPrefix(innerCommand, "("), ")")
			innerArgs := strings.Fields(innerCommand)

			// Выполняем команду
			dataRecord, err := ExecOpenedComms(innerArgs)
			if err != nil {
				fmt.Println(system.Red(fmt.Sprintf("Error: %s", err)))
				return true
			}

			// Заменяем в основном массиве команду на результат выполнения
			fullCommandArgs[iArg] = dataRecord

			// Удаляем оставшиеся элементы команды из аргументов
			fullCommandArgs = append(fullCommandArgs[:iArg+1], fullCommandArgs[endIdx+1:]...)
			continue
		}
	}

	*commandArgs = fullCommandArgs

	session.Path = system.UserDir

	isValid := utils.ValidCommand(*commandLower, system.Commands)

	if len(strings.TrimSpace(*commandLower)) != len(strings.TrimSpace(*commandLine)) && isValid {
		session.CommandHistory = append(session.CommandHistory, *commandLine)
		system.GlobalSession.CommandHistory = session.CommandHistory
	}

	if !isValid {
		system.ExecutingCommand = true

		_ = ExecExternalLoopCommand(
			session,
			system.UserDir,
			*command,
			*commandInput,
			*commandLine,
			*commandLower,
			*commandArgs,
			*echoTime,
			*runOnNewThread,
		)

		system.ExecutingCommand = false

		//if err != nil {
		//	return true
		//}
		//
		//if strings.TrimSpace(*commandInput) != "" {
		//	return true
		//}
		//
		//return true
	}

	return false
}

func ExecOpenedComms(commandArgs []string) (string, error) {
	if len(commandArgs) < 1 {
		return "", errs.CommandArgsNotFound
	}

	commandMap := map[string]func() ([]byte, error){
		"read":   func() ([]byte, error) { return ExOpenedCommReadUtils.File(commandArgs[1]) },
		"create": func() ([]byte, error) { return fcommands.CreateFile(commandArgs[1]) },
	}

	if handler, exists := commandMap[commandArgs[0]]; exists {
		bytes, err := handler()
		if err != nil {
			return "", err
		}

		return string(bytes), nil
	}

	return "", errors.New("command not found")
}
