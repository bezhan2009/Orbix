package src

import (
	"errors"
	"fmt"
	"goCmd/cmd/commands"
	ExOpenedCommReadUtils "goCmd/cmd/commands/Read/utils"
	"goCmd/cmd/dirInfo"
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

func EndOfSessions(originalStdout, originalStderr *os.File,
	session *system.Session,
	sessionData *system.AppState,
	prefix string) {
	// Restore original outputs
	os.Stdout, os.Stderr = originalStdout, originalStderr

	PreviousSessionPath = session.Path
	session, _ = sessionData.GetSession(PreviousSessionPrefix)

	if err := commands.ChangeDirectory(session.Path); err != nil {
		fmt.Println(red("Error changing directory:", err))
	}

	sessionData.DeleteSession(prefix)

	system.OrbixWorking = false
	UnknownCommandsCounter = 0
}

func ExecLoopCommand(commandLower,
	prefix string,
	echoTime, runOnNewThread bool,
	execCommand structs.ExecuteCommandFuncParams) error {
	if UnknownCommandsCounter != 0 {
		return nil
	}

	execCommandCatchErrs := structs.ExecuteCommandCatchErrs{
		CommandLower:   commandLower,
		EchoTime:       &echoTime,
		RunOnNewThread: &runOnNewThread,
	}

	if strings.TrimSpace(commandLower) == "orbix" && *execCommand.IsWorking {
		PreviousSessionPrefix = prefix
	}

	if catchSyntaxErrs(execCommandCatchErrs) {
		return errors.New("continue loop")
	}

	if runOnNewThread {
		go ExecuteCommand(execCommand)
	} else {
		if echoTime {
			// Запоминаем время начала
			startTime := time.Now()
			ExecuteCommand(execCommand)
			// Выводим время выполнения
			TEXCOM := fmt.Sprintf("Command executed in: %s\n", time.Since(startTime))
			fmt.Println(green(TEXCOM))
		} else {
			ExecuteCommand(execCommand)
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
	GlobalSession.CommandHistory = session.CommandHistory

	if commandFile(strings.TrimSpace(commandLower)) {
		fullFileName(&commandArgs)
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
			startTime := time.Now()
			err = executeCommandOrbix(fullCommand, command, commandLower, dir)
			// Выводим время выполнения
			TEXCOM := fmt.Sprintf("Command executed in: %s\n", time.Since(startTime))
			fmt.Println(green(TEXCOM))

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
	if _, exists := editableVars[varName]; exists {
		return GetVariableValue(varName)
	}

	return nil, errors.New("Variable not found: " + varName)
}

func ExecCommandPromptLogic(
	firstCharIs,
	lastCharIs,
	isComHasFlag,
	echoTime,
	runOnNewThread bool,
	dir string,
	commandArgs *[]string,
	prompt, command, commandLine, commandInput, commandLower *string,
	session *system.Session) bool {
	if isComHasFlag && (echoTime || runOnNewThread) {
		*commandLine = removeFlags(*commandLine)
		*commandInput = removeFlags(*commandInput)
		*commandLine, *command, *commandArgs, *commandLower = readCommandLine(*commandLine) // Refactored input handling
	}

	if firstCharIs && lastCharIs {
		*commandLower = "print"
		*commandLine = fmt.Sprintf("print %s", *commandLine)
		*commandLine, *command, *commandArgs, *commandLower = readCommandLine(*commandLine) // Refactored input handling
	}

	commandArgsListStr := *commandArgs
	fullCommandArgs := []string{}
	commandExLogg := *command

	for _ = range *commandLine {
		if string(commandExLogg[0]) == "$" && len(strings.TrimSpace(commandExLogg)) > 1 {
			if string(commandExLogg[len(commandExLogg)-1]) == "-" {
				commandList := strings.Split(commandExLogg, "-")
				commandExLogg = commandList[0]
				break
			}

			customVar, err := getCustomVar(commandExLogg[1:])
			if err != nil {
				fmt.Println(red(err))
			}

			if customVar != nil {
				commandExLogg = customVar.(string)
				*commandLine, *command, *commandArgs, *commandLower = readCommandLine(commandExLogg) // Refactored input handling
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
				argList := strings.Split(arg, "-")
				arg = argList[0]
				fullCommandArgs[iArg] = arg
				*commandArgs = fullCommandArgs
				continue
			}

			customVar, err := getCustomVar(arg[1:])
			if err != nil {
				fmt.Println(red(err))
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
				fmt.Println(red(fmt.Sprintf("Syntax Error in '%s'\nPlease Close the '('", fullCommandArgs[iArg])))
				fmt.Println(fmt.Sprintf(" %s\n%s", fmt.Sprintf("%s%s", red(fmt.Sprintf("%s", fullCommandArgs[iArg])), yellow(")")), strings.Repeat(red("-"), (len(arg)-1)+2)+yellow("^")))
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

	session.Path = dir

	isValid := utils.ValidCommand(*commandLower, Commands)

	if len(strings.TrimSpace(*commandLower)) != len(strings.TrimSpace(*commandLine)) && isValid {
		session.CommandHistory = append(session.CommandHistory, *commandLine)
		GlobalSession.CommandHistory = session.CommandHistory
	}

	if !isValid {
		err := ExecExternalLoopCommand(
			session,
			dir,
			*command,
			*commandInput,
			*commandLine,
			*commandLower,
			*commandArgs,
			echoTime,
			runOnNewThread,
		)

		if err != nil {
			return true
		}

		if strings.TrimSpace(*commandInput) != "" {
			return true
		}

		return true
	}

	if strings.TrimSpace(*commandLower) == "prompt" {
		handlePromptCommand(*commandArgs, prompt)
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
		fmt.Printf("\n%s\n", green(RecoverText))
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
			isValid := utils.ValidCommand(commandLowerEx, AdditionalCommands)
			if !isValid {
				HandleUnknownCommandUtil(commandEx, Commands)
				return err
			} else {
				return errors.New("continue loop")
			}
		} else {
			return errors.New("continue loop")
		}
	} else {
		UnknownCommandsCounter = UnknownCommandsCounter + 1
	}

	return nil
}

func usingForLT(commandInput string) bool {
	if strings.TrimSpace(commandInput) != "" && strings.TrimSpace(commandInput) != "restart" {
		return false
	}

	return true
}

func execLtCommand(commandInput string) {
	user := "OneCom"

	// Initialize colors
	InitColors()

	if strings.TrimSpace(Location) == "" {
		Location = os.Getenv("CITY")
		if strings.TrimSpace(Location) == "" {
			Location = string(strings.TrimSpace(os.Getenv("USERS_LOCATION")))
		}
	}

	dir, _ := os.Getwd()
	dirC := dirInfo.CmdDir(dir)

	printPromptInfoWithoutGit(Location, user, dirC, commandInput)

	commandLine, command, commandArgs, commandLower := readCommandLine(commandInput) // Refactored input handling
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
	GlobalSession = session

	execCommand := structs.ExecuteCommandFuncParams{
		Command:       command,
		CommandLower:  commandLower,
		CommandArgs:   commandArgs,
		Dir:           dir,
		IsWorking:     &isWorking,
		IsPermission:  &PermissionDenied,
		Username:      user,
		SD:            &sessionData,
		SessionPrefix: "",
		Session:       &session,
		GlobalSession: &GlobalSession,
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
		IsWorking:      &isWorking,
		IsComHasFlag:   &isComHasFlag,
		Session:        &session,
		ExecCommand:    execCommand,
	}

	startTimePRCOMARGS := time.Now()
	continueLoop := processCommandArgs(processCommandParams)

	if continueLoop {
		if echoTime {
			TEXCOMARGS := fmt.Sprintf("Command executed in: %s\n", time.Since(startTimePRCOMARGS))
			fmt.Println(green(TEXCOMARGS))
			return
		}
		return
	}

	if isComHasFlag && (echoTime || runOnNewThread) {
		commandLine = removeFlags(commandLine)
		commandInput = removeFlags(commandInput)
		commandLine, command, commandArgs, commandLower = readCommandLine(commandLine) // Refactored input handling
	}

	if firstCharIs && lastCharIs {
		commandLower = "print"
		commandLine = fmt.Sprintf("print %s", commandLine)
		commandLine, command, commandArgs, commandLower = readCommandLine(commandLine) // Refactored input handling
	}

	isValid := utils.ValidCommand(commandLower, Commands)

	if len(strings.TrimSpace(commandLower)) != len(strings.TrimSpace(commandLine)) && isValid {
		session.CommandHistory = append(session.CommandHistory, commandLine)
		GlobalSession.CommandHistory = session.CommandHistory
	}

	if !isValid {
		session.CommandHistory = append(session.CommandHistory, commandLine)
		GlobalSession.CommandHistory = session.CommandHistory

		if commandFile(strings.TrimSpace(commandLower)) {
			fullFileName(&commandArgs)
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
				startTime := time.Now()
				executeCommandOrbix(fullCommand, command, commandLower, dir)
				// Выводим время выполнения
				TEXCOM := fmt.Sprintf("Command executed in: %s\n", time.Since(startTime))
				fmt.Println(green(TEXCOM))

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
		Dir:           dir,
		IsWorking:     &isWorking,
		IsPermission:  &PermissionDenied,
		Username:      "OneCom",
		SD:            &sessionData,
		SessionPrefix: "",
		Session:       &session,
		GlobalSession: &GlobalSession,
	}

	execCommandCatchErrs := structs.ExecuteCommandCatchErrs{
		EchoTime:       &echoTime,
		RunOnNewThread: &runOnNewThread,
	}

	if catchSyntaxErrs(execCommandCatchErrs) {
		return
	}

	if runOnNewThread {
		go ExecuteCommand(execCommand)
	} else {
		if echoTime {
			// Запоминаем время начала
			startTime := time.Now()
			ExecuteCommand(execCommand)
			// Выводим время выполнения
			TEXCOM := fmt.Sprintf("Command executed in: %s\n", time.Since(startTime))
			fmt.Println(green(TEXCOM))
		} else {
			ExecuteCommand(execCommand)
		}
	}
}

func processCommandArgs(processCommandParams structs.ProcessCommandParams) (continueLoop bool) {
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
				}
			}
		}
	}

	for index, commandLetter := range processCommandParams.CommandLine {
		if (string(commandLetter) == string('"') || string(commandLetter) == "'") && index == 0 {
			*processCommandParams.FirstCharIs = true
		} else if (string(commandLetter) == string('"') || string(commandLetter) == "'") && index == len(processCommandParams.CommandLine)-1 {
			*processCommandParams.LastCharIs = true
		}
	}

	if commandInt, err := strconv.Atoi(processCommandParams.Command); err == nil && len(processCommandParams.CommandArgs) == 0 {
		fmt.Println(magenta(commandInt))
		return true
	}

	if strings.TrimSpace(processCommandParams.CommandLower) == "neofetch" && *processCommandParams.IsWorking && system.OperationSystem == "windows" {
		neofetchUser := User

		if User == "" {
			neofetchUser = processCommandParams.Session.User
		}

		if *processCommandParams.RunOnNewThread {
			go ExCommUtils.NeofetchUtil(processCommandParams.ExecCommand, neofetchUser, Commands)
		} else {
			ExCommUtils.NeofetchUtil(processCommandParams.ExecCommand, neofetchUser, Commands)
		}

		if strings.TrimSpace(processCommandParams.CommandInput) != "" {
			*processCommandParams.IsWorking = false
		}

		return true
	}

	return false
}

func executeGoCode(code string) error {
	// Создаем временный файл для исходного кода
	tmpFile, err := ioutil.TempFile(".", "tempcode*.go")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name()) // Удаляем файл после выполнения

	// Записываем код в файл
	if _, err := tmpFile.Write([]byte(code)); err != nil {
		return err
	}

	// Закрываем временный файл
	tmpFile.Close()

	// Компилируем и запускаем код
	cmd := exec.Command("go", "run", tmpFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error: %s, output: %s", err.Error(), output)
	}

	fmt.Println(string(output))
	return nil
}
