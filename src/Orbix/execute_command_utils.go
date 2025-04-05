package Orbix

import (
	"errors"
	"fmt"
	"goCmd/pkg/cache"
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
	sumArgs            string
	commandArgsListStr []string
	commandList        []string
	fullCommandArgs    []string
	argList            []string
	commExLoggSplit    []string
	argSplit           []string
	iArgSplit          int

	startTime      time.Time
	printed_timing string

	commandLineFr *string
	isFormated    bool
)

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
			if printed_timing == "" {
				TEXCOM = fmt.Sprintf("Command executed in: %s\n",
					time.Since(startTime))
				fmt.Println(system.Green(TEXCOM))
			} else {
				printed_timing = ""
			}
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
			printed_timing = TEXCOM

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

func executeCommandOrbix(fullCommandEx []string, commandEx, commandLowerEx, dirEx string) error {
	// Если команда с путём, формируем полный путь относительно dirEx
	fullPath := filepath.Join(dirEx, commandEx)
	extension := strings.ToLower(filepath.Ext(fullPath))
	fullCommandEx[0] = fullPath

	if extension == ".exe" {
		return utils.ExternalCommand(fullCommandEx)
	}

	// Если команда задана без разделителя, ищем её в PATH
	if !strings.Contains(commandEx, string(os.PathSeparator)) {
		fullCommandEx[0] = commandEx
		if err := utils.ExternalCommand(fullCommandEx); err != nil {
			// Можно логировать ошибку, если нужно
			return err
		}
		src.UnknownCommandsCounter++
		return nil
	}

	// Если файл не найден, используем fallback
	if cache.GetCommandFromCache(fullPath) == "" {
		return fallbackExecution(fullCommandEx, commandEx, commandLowerEx, dirEx)
	}

	if err := utils.ExternalCommand(fullCommandEx); err != nil {
		return fallbackExecution(fullCommandEx, commandEx, commandLowerEx, dirEx)
	}

	src.UnknownCommandsCounter++
	return nil
}

func fallbackExecution(fullCommandEx []string, commandEx, commandLowerEx, dirEx string) error {
	fullPath := filepath.Join(dirEx, commandEx)
	fullCommandEx[0] = fullPath
	extension := strings.ToLower(filepath.Ext(fullPath))
	var cmd *exec.Cmd
	switch extension {
	case ".ps1":
		cmd = exec.Command("powershell", "-File", fullPath)
	case ".py":
		cmd = exec.Command("python", fullPath)
	default:
		cmd = exec.Command("cmd.exe", "/C", fullPath)
	}
	return cmdRun(cmd, commandEx, commandLowerEx)
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
