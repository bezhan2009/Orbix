package Orbix

import (
	"errors"
	"fmt"
	"goCmd/cmd/commands"
	ExOpenedCommReadUtils "goCmd/cmd/commands/Read/utils"
	"goCmd/cmd/dirInfo"
	"goCmd/src"
	"goCmd/src/environment"
	"goCmd/src/handlers"
	"goCmd/src/user"
	ExCommUtils "goCmd/src/utils"
	"goCmd/structs"
	"goCmd/system"
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
	commandArgsListStr []string
	commandList        []string
	fullCommandArgs    []string
	argList            []string

	startTime time.Time
)

func EndOfSessions(originalStdout, originalStderr *os.File,
	session *system.Session,
	sessionData *system.AppState,
	prefix string) {
	// Restore original outputs
	os.Stdout, os.Stderr = originalStdout, originalStderr

	system.PreviousSessionPath = session.Path
	session, _ = sessionData.GetSession(system.PreviousSessionPrefix)

	if err := commands.ChangeDirectory(session.Path); err != nil {
		fmt.Println(system.Red("Error changing directory:", err))
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
	commandArgsListStr = []string{}
	commandList = []string{}
	fullCommandArgs = []string{}
	argList = []string{}
	startTime = time.Time{}
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
		if string(commandExLogg[0]) == "$" && len(strings.TrimSpace(commandExLogg)) > 1 {
			if string(commandExLogg[len(commandExLogg)-1]) == "-" {
				commandList = strings.Split(commandExLogg, "-")
				commandExLogg = commandList[0]
				break
			}

			customVar, err := getCustomVar(commandExLogg[1:])
			if err != nil {
				fmt.Println(system.Red(err))
			}

			if customVar != nil {
				commandExLogg = customVar.(string)
				*commandLine, *command, *commandArgs, *commandLower = src.ReadCommandLine(commandExLogg) // Refactored input handling
				fullCommandArgs = append(fullCommandArgs, *commandArgs...)
			}
		}

		break
	}

	fullCommandArgs = append(fullCommandArgs, commandArgsListStr...)
	*commandArgs = fullCommandArgs

	for iArg, arg := range *commandArgs {
		if string(arg[0]) == "$" && len(strings.TrimSpace(arg)) > 1 {
			if string(arg[len(arg)-1]) == "-" {
				argList = strings.Split(arg, "-")
				arg = argList[0]
				fullCommandArgs[iArg] = arg
				*commandArgs = fullCommandArgs
				continue
			}

			customVar, err := getCustomVar(arg[1:])
			if err != nil {
				fmt.Println(system.Red(err))
			}

			if customVar != nil {
				fullCommandArgs[iArg] = customVar.(string)
				*commandArgs = fullCommandArgs
			}
		}
	}

	InQuotes := true

	for iArg := 0; iArg < len(fullCommandArgs); iArg++ {
		arg := fullCommandArgs[iArg]

		if strings.HasPrefix(arg, "(") && InQuotes {
			// Ищем конец команды в скобках
			endIdx := -1
			for j := iArg; j < len(fullCommandArgs); j++ {
				if strings.HasSuffix(fullCommandArgs[j], ")") {
					endIdx = j
					continue
				}
			}

			// Если нет закрывающей скобки
			if endIdx == -1 {
				fmt.Println(system.Red(fmt.Sprintf("Syntax Error in '%s'\nPlease Close the '('", fullCommandArgs[iArg])))
				fmt.Println(fmt.Sprintf(" %s\n%s", fmt.Sprintf("%s%s", system.Red(fmt.Sprintf("%s", fullCommandArgs[iArg])), system.Yellow(")")), strings.Repeat(system.Red("-"), (len(arg)-1)+2)+system.Yellow("^")))
				continue
			}

			// Объединяем команду внутри скобок
			innerCommand := strings.Join(fullCommandArgs[iArg:endIdx+1], " ")
			innerCommand = strings.TrimSuffix(strings.TrimPrefix(innerCommand, "("), ")")
			innerArgs := strings.Fields(innerCommand)

			// Выполняем команду
			dataRecord, err := ExecOpenedComms(innerArgs)
			if err != nil {
				fmt.Println(fmt.Sprintf("Ошибка выполнения: %s", err))
				continue
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
		err := ExecExternalLoopCommand(
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

		if err != nil {
			return true
		}

		if strings.TrimSpace(*commandInput) != "" {
			return true
		}

		return true
	}

	return false
}

func ExecOpenedComms(commandArgs []string) (string, error) {
	commandMap := map[string]func() ([]byte, error){
		"read": func() ([]byte, error) { return ExOpenedCommReadUtils.File(commandArgs[1]) },
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
}

// Command execution logic that can be run in a new thread
func executeCommandOrbix(fullCommandEx []string,
	commandEx, commandLowerEx, dirEx string) error {
	err := utils.ExternalCommand(fullCommandEx)
	if err != nil {
		fullPath := filepath.Join(dirEx, commandEx)
		fullCommandEx[0] = fullPath

		// Проверяем расширение файла и выбираем подходящий интерпретатор
		extension := strings.ToLower(filepath.Ext(fullPath))
		var cmd *exec.Cmd
		switch extension {
		case ".bat":
			cmd = exec.Command("cmd.exe", "/C", fullPath)
		case ".ps1":
			cmd = exec.Command("powershell", "-File", fullPath)
		case ".py":
			cmd = exec.Command("python", fullPath)
		default:
			cmd = exec.Command(fullPath)
		}

		// Запускаем команду и обрабатываем ошибки
		err = cmd.Run()
		if err != nil {
			isValid := utils.ValidCommand(commandLowerEx, system.AdditionalCommands)
			if !isValid {
				handlers.HandleUnknownCommandUtil(commandEx, commandLowerEx, system.Commands)
				return err
			} else {
				return errors.New("continue loop")
			}
		} else {
			return errors.New("continue loop")
		}
	} else {
		src.UnknownCommandsCounter = src.UnknownCommandsCounter + 1
	}

	return nil
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

	isValid := utils.ValidCommand(commandLower, system.Commands)

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
		for _, commandLetter := range processCommandParams.CommandLine {
			if commandLetter == '-' {
				*processCommandParams.IsComHasFlag = true
				break // Прерываем цикл, если флаг найден
			}
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
	}

	for index, commandLetter := range processCommandParams.CommandLine {
		if (string(commandLetter) == string('"') || string(commandLetter) == "'") && index == 0 {
			*processCommandParams.FirstCharIs = true
		} else if (string(commandLetter) == string('"') || string(commandLetter) == "'") && index == len(processCommandParams.CommandLine)-1 {
			*processCommandParams.LastCharIs = true
		}
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
			go ExCommUtils.NeofetchUtil(&processCommandParams.ExecCommand, neofetchUser, system.Commands)
		} else {
			ExCommUtils.NeofetchUtil(&processCommandParams.ExecCommand, neofetchUser, system.Commands)
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
