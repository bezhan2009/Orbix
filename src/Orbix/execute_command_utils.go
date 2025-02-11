package Orbix

import (
	"errors"
	"fmt"
	"goCmd/cmd/commands"
	ExOpenedCommReadUtils "goCmd/cmd/commands/Read/utils"
	"goCmd/cmd/commands/fcommands"
	"goCmd/cmd/dirInfo"
	"goCmd/src"
	"goCmd/src/environment"
	"goCmd/src/handlers"
	"goCmd/src/user"
	ExCommUtils "goCmd/src/utils"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/system/errs"
	"goCmd/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	TEXCOM             string
	commandExLogg      string
	sumArgs            string
	commandArgsListStr []string
	commandList        []string
	fullCommandArgs    []string
	argList            []string
	commExLoggSplit    []string
	argSplit           []string
	iArgSplit          int

	startTime time.Time
)

func EndOfSessions(originalStdout, originalStderr *os.File,
	session *system.Session,
	sessionData *system.AppState,
	prefix string) {
	system.CntLaunchedOrbixes--

	// Restore original outputs
	os.Stdout, os.Stderr = originalStdout, originalStderr

	system.PreviousSessionPath = session.Path
	session, _ = sessionData.GetSession(system.PreviousSessionPrefix)

	if strings.TrimSpace(session.Path) != "" {
		if err := commands.ChangeDirectory(session.Path); err != nil {
			fmt.Println(system.Red("Error changing directory:", err))
		}
	}

	sessionData.DeleteSession(prefix)

	system.OrbixWorking = false
	src.UnknownCommandsCounter = 0
}

func ExecLoopCommand(commandLower,
	prefix *string,
	echoTime, runOnNewThread *bool,
	execCommand *structs.ExecuteCommandFuncParams) error {

	if src.UnknownCommandsCounter != 0 {
		return nil
	}

	execCommandCatchErrs := structs.ExecuteCommandCatchErrs{
		CommandLower:   *commandLower,
		EchoTime:       echoTime,
		RunOnNewThread: runOnNewThread,
	}

	if strings.TrimSpace(*commandLower) == "orbix" && *execCommand.LoopData.IsWorking {
		system.PreviousSessionPrefix = *prefix
	}

	if src.CatchSyntaxErrs(execCommandCatchErrs) {
		return errors.New("continue loop")
	}

	if *runOnNewThread {
		go Command(execCommand)
	} else {
		if *echoTime {
			// Запоминаем время начала
			startTime = time.Now()
			Command(execCommand)
			// Выводим время выполнения
			TEXCOM = fmt.Sprintf("Command executed in: %s\n",
				time.Since(startTime))
			fmt.Println(system.Green(TEXCOM))
		} else {
			Command(execCommand)
		}
	}

	return nil
}

func ExecExternalLoopCommand(session *system.Session,
	dir, command, commandInput, commandLine, commandLower string,
	commandArgs []string,
	echoTime, runOnNewThread bool) error {
	errLtCommand := errors.New("LtCommand")

	session.CommandHistory = append(session.CommandHistory, commandLine)
	system.GlobalSession.CommandHistory = session.CommandHistory

	if src.CommandFile(strings.TrimSpace(commandLower)) {
		src.FullFileName(&commandArgs)
	}

	fullCommand := append([]string{command}, commandArgs...)

	var err error

	if runOnNewThread {
		go func() {
			err = executeCommandOrbix(fullCommand, command, commandLower, dir)
			if err != nil {
				fmt.Println(system.Red("Error executing orbix command:", err))
			}
		}()

		if strings.TrimSpace(commandInput) != "" {
			return errLtCommand
		}
	} else {
		if echoTime {
			// Запоминаем время начала
			startTime = time.Now()
			err = executeCommandOrbix(fullCommand, command, commandLower, dir)
			// Выводим время выполнения
			TEXCOM = fmt.Sprintf("Command executed in: %s\n", time.Since(startTime))
			fmt.Println(system.Green(TEXCOM))

			if strings.TrimSpace(commandInput) != "" {
				return errLtCommand
			}
		} else {
			err = executeCommandOrbix(fullCommand, command, commandLower, dir)
			if err != nil {
				return err
			}

			if strings.TrimSpace(commandInput) != "" {
				return errLtCommand
			}
		}
	}

	return err
}

func getCustomVar(varName string) (interface{}, error) {
	if _, exists := system.EditableVars[varName]; exists {
		return environment.GetVariableValue(varName)
	}

	return nil, errors.New("Variable not found: " + varName)
}

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

func RecoverAndRestore(rebooted *structs.RebootedData) {
	if rebooted.Recover != nil {
		RecoverText := fmt.Sprintf("Successfully recovered from the panic: %v", rebooted.Recover)
		fmt.Printf("\n%s\n", system.Green(RecoverText))
		rebooted.Recover = nil
	}

	system.OrbixRecovering = false
}

// Command execution logic that can be run in a new thread
func executeCommandOrbix(fullCommandEx []string,
	commandEx, commandLowerEx, dirEx string) error {

	var fullCommandExCopy []string

	fullCommandExCopy = append(fullCommandExCopy, fullCommandEx...)

	fullPath := filepath.Join(dirEx, commandEx)
	fullCommandEx[0] = fullPath

	// Проверяем расширение файла и выбираем подходящий интерпретатор
	extension := strings.ToLower(filepath.Ext(fullPath))

	if extension == ".exe" {
		err := utils.ExternalCommand(fullCommandEx)
		if err != nil {
			fullCommandEx = fullCommandExCopy
		}

		return err
	} else {
		fullCommandEx = fullCommandExCopy
	}

	err := utils.ExternalCommand(fullCommandEx)
	if err != nil {
		fullPath = filepath.Join(dirEx, commandEx)
		fullCommandEx[0] = fullPath

		// Проверяем расширение файла и выбираем подходящий интерпретатор
		extension = strings.ToLower(filepath.Ext(fullPath))
		var cmd *exec.Cmd
		switch extension {
		case ".ps1":
			cmd = exec.Command("powershell", "-File", fullPath)
		case ".py":
			cmd = exec.Command("python", fullPath)
		default:
			cmd = exec.Command("cmd.exe", "/C", fullPath)
		}

		return cmdRun(cmd,
			commandEx, commandLowerEx)
	} else {
		src.UnknownCommandsCounter = src.UnknownCommandsCounter + 1
	}

	return nil
}

func cmdRun(cmd *exec.Cmd,
	commandEx, commandLowerEx string) (err error) {
	// Запускаем команду и обрабатываем ошибки
	err = cmd.Run()
	if err != nil {
		isValid := utils.ValidCommandFast(commandLowerEx, system.AdditionalCmdMap)
		if !isValid {
			handlers.HandleUnknownCommandUtil(commandEx, commandLowerEx, system.CmdMap)
			return err
		} else {
			return errors.New("continue loop")
		}
	} else {
		return errors.New("continue loop")
	}
}

func ExecLtCommand(commandInput string) {
	user := "OneCom"

	// Initialize colors
	system.InitColors()

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

func ProcessCommandArgs(processCommandParams *structs.ProcessCommandParams) (continueLoop bool) {
	*processCommandParams.RunOnNewThread = false
	*processCommandParams.EchoTime = false
	*processCommandParams.FirstCharIs = false
	*processCommandParams.LastCharIs = false

	if processCommandParams.CommandLower == "signout" {
		user.DeleteUserFromRunningFile(system.UserName)
		*processCommandParams.LoopData.IsWorking = false
		return true
	}

	if len(processCommandParams.CommandArgs) > 0 {
		*processCommandParams.IsComHasFlag = strings.Contains(processCommandParams.CommandLine, "-")
	}

	if *processCommandParams.IsComHasFlag {
		// Проходим по всем аргументам
		for i := len(processCommandParams.CommandArgs) - 1; i >= 0; i-- {
			arg := strings.ToLower(strings.TrimSpace(processCommandParams.CommandArgs[i]))

			switch arg {
			case "--run-in-new-thread":
				*processCommandParams.RunOnNewThread = true
				// Удаляем аргумент из списка
				processCommandParams.CommandArgs = append(processCommandParams.CommandArgs[:i], processCommandParams.CommandArgs[i+1:]...)
			case "--timing", "-t":
				*processCommandParams.EchoTime = true
				// Удаляем аргумент из списка
				processCommandParams.CommandArgs = append(processCommandParams.CommandArgs[:i], processCommandParams.CommandArgs[i+1:]...)
			default:
				*processCommandParams.IsComHasFlag = false
				*processCommandParams.RunOnNewThread = false
				*processCommandParams.EchoTime = false
			}
		}
	} else {
		*processCommandParams.RunOnNewThread = false
		*processCommandParams.EchoTime = false
	}

	if rune(processCommandParams.CommandLine[0]) == '"' || string(processCommandParams.CommandLine[0]) == "'" {
		*processCommandParams.FirstCharIs = true
	}

	if rune(processCommandParams.CommandLine[len(processCommandParams.CommandLine)-1]) == '"' || string(processCommandParams.CommandLine[len(processCommandParams.CommandLine)-1]) == "'" {
		*processCommandParams.LastCharIs = true
	}

	if commandInt, err := strconv.Atoi(processCommandParams.Command); err == nil && len(processCommandParams.CommandArgs) == 0 && commandInt < system.MaxInt {
		fmt.Println(system.Magenta(commandInt))
		return true
	}

	if strings.TrimSpace(processCommandParams.CommandLower) == "neofetch" && *processCommandParams.LoopData.IsWorking && system.OperationSystem == "windows" {
		neofetchUser := system.User

		if system.User == "" {
			neofetchUser = processCommandParams.LoopData.Session.User
		}

		if *processCommandParams.RunOnNewThread {
			go ExCommUtils.NeofetchUtil(&processCommandParams.ExecCommand, neofetchUser, system.CmdMap)
		} else {
			ExCommUtils.NeofetchUtil(&processCommandParams.ExecCommand, neofetchUser, system.CmdMap)
		}

		if strings.TrimSpace(processCommandParams.CommandInput) != "" {
			*processCommandParams.LoopData.IsWorking = false
		}

		return true
	}

	return false
}

func executeGoCode(code string) {
	var execGoCode string

	if !strings.HasSuffix(code, "package") || !strings.HasSuffix(code, "main") {
		execGoCode = fmt.Sprintf(`
package main

import "fmt"

func main() {
	%s
}
`, code)
	}
	// Создаем временный файл для исходного кода
	tmpFile, err := ioutil.TempFile(".", "tempcode*.go")
	if err != nil {
		fmt.Println(system.Red("Error creating TempFile:", err))
		return
	}
	defer os.Remove(tmpFile.Name()) // Удаляем файл после выполнения

	// Записываем код в файл
	if _, err = tmpFile.Write([]byte(execGoCode)); err != nil {
		fmt.Println(system.Red("Error writing to TempFile:", err))
		return
	}

	// Закрываем временный файл
	tmpFile.Close()

	// Компилируем и запускаем код
	cmd := exec.Command("go", "run", tmpFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(system.Red("Error running Go code:", err))
		return
	}

	fmt.Println(system.Magenta(string(output)))
}
