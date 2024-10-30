package src

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/fsnotify/fsnotify"
	"goCmd/cmd/commands"
	ExCommUtils "goCmd/src/utils"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// New helper functions
func initializeRunningFile(username string) {
	// Check and initialize running.txt if not exists
	if _, err := os.Stat("running.txt"); os.IsNotExist(err) {
		if _, err = os.Create("running.txt"); err != nil {
			panic(err)
		}
	}

	// Check for username in running.txt and add if missing
	runningPath := filepath.Join(Absdir, "running.txt")
	if sourceRunning, err := os.ReadFile(runningPath); err == nil {
		if !strings.Contains(string(sourceRunning), username) {
			if file, err := os.OpenFile("running.txt", os.O_APPEND|os.O_WRONLY, 0644); err == nil {
				defer func() {
					err = file.Close()
					if err != nil {
						return
					}
				}()
				if _, err := file.WriteString("\n" + username + "\n"); err != nil {
					fmt.Println("Error writing to running.txt:", err)
				}
			}
		}
	}
}

func checkUserInRunningFile(username string) bool {
	runningPath := filepath.Join(Absdir, "running.txt")
	sourceRunning, err := os.ReadFile(runningPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(sourceRunning), username)
}

func getUser(username string) string {
	if strings.TrimSpace(User) != "" {
		return User
	} else {
		return username
	}
}

func printPromptInfo(location, user, dirC, commandInput string, sd *system.Session) {
	if len(Prompt) > 1 {
		Prompt = string(Prompt[0])
	}

	fmt.Printf("\n%s%s%s%s%s%s%s%s %s%s%s%s%s%s%s%s%s%s%s\n",
		yellow("┌"), yellow("─"), yellow("("), cyan("Orbix@"+getUser(user)), yellow(")"), yellow("─"), yellow("["),
		yellow(location), magenta(time.Now().Format("15:04")), yellow("]"), yellow("─"), yellow("["),
		cyan("~"), cyan(dirC), yellow("]"), yellow(" git:"), green("["), green(sd.GitBranch), green("]"))
	fmt.Printf("%s%s%s %s",
		yellow("└"), yellow("─"), green(strings.TrimSpace(Prompt)), green(commandInput))

	if strings.TrimSpace(commandInput) != "" && len(os.Args) > 0 {
		fmt.Println()
	}
}

func printPromptInfoWithoutGit(location, user, dirC, commandInput string) {
	if len(Prompt) > 1 {
		Prompt = string(Prompt[0])
	}

	fmt.Printf("\n%s%s%s%s%s%s%s%s %s%s%s%s%s%s%s\n",
		yellow("┌"), yellow("─"), yellow("("), cyan("Orbix@"+getUser(user)), yellow(")"), yellow("─"), yellow("["),
		yellow(location), magenta(time.Now().Format("15:04")), yellow("]"), yellow("─"), yellow("["),
		cyan("~"), cyan(dirC), yellow("]"))
	fmt.Printf("%s%s%s %s",
		yellow("└"), yellow("─"), green(strings.TrimSpace(Prompt)), green(commandInput))

	if strings.TrimSpace(commandInput) != "" && len(os.Args) > 0 {
		fmt.Println()
	}
}

func readCommandLine(commandInput string) (string, string, []string, string) {
	var commandLine string
	if commandInput != "" {
		commandLine = strings.TrimSpace(commandInput)
	} else {
		// Чтение ввода
		commandLine = strings.TrimSpace(prompt.Input("", autoComplete))
	}

	commandParts := utils.SplitCommandLine(commandLine)
	if len(commandParts) == 0 {
		return "", "", nil, ""
	}

	command := commandParts[:1]

	return commandLine, command[0], commandParts[1:], strings.ToLower(commandParts[0])
}

func processCommand(commandLower string) (bool, error) {
	if strings.TrimSpace(commandLower) == "cd" && GitCheck {
		return true, nil
	}

	if strings.TrimSpace(commandLower) == "git" && GitCheck {
		return true, nil
	}

	if commandLower == "signout" {
		return false, fmt.Errorf("signout")
	}

	return false, nil
}

// Функция обработки каждой команды
func handleCommand(command string) (ok bool) {
	// Пример конкатенации строк через оператор "+"
	if strings.Contains(command, "+") {
		parts := strings.Split(command, "+")
		var result string
		for _, part := range parts {
			// Убираем кавычки и пробелы
			cleanPart := strings.Trim(part, "\" ")
			result += cleanPart
		}

		if result != command {
			fmt.Println(result)
			return true
		}

		return false
	} else {
		// Вывод просто строки (если нет "+")
		cleanCommand := strings.Trim(command, "\"")
		if command != cleanCommand {
			fmt.Println(cleanCommand)
			return true
		}

		return false
	}
}

func createNewSession(path, user, gitBranch string, isAdmin bool) *system.Session {
	session := &system.Session{
		Path:      path,
		User:      user,
		GitBranch: gitBranch,
		IsAdmin:   isAdmin,
	}
	return session
}

func restorePreviousSession(sessionData *system.AppState, prefix string) *system.Session {
	session, exists := sessionData.GetSession(prefix)
	if !exists {
		fmt.Println(red("Session does not exist!"))
		return nil
	}
	return session
}

// Map, ограничивающий изменяемые переменные
var editableVars = map[string]interface{}{
	"location": &Location,
	"prompt":   &Prompt,
	"user":     &User,
}

var availableEditableVars = []string{"location", "prompt", "user"}

func watchFile(runningPath string, username string, isWorking *bool, isPermission *bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = watcher.Close()
		if err != nil {
			return
		}
	}()

	done := make(chan bool)

	// Запускаем горутину для отслеживания событий
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write && *isWorking {
					if !checkUserInRunningFile(username) && *isWorking {
						fmt.Print(red("\nUser not authorized. to continue, press Enter:"))
						devNull, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0666)
						func() {
							err = devNull.Close()
							if err != nil {
								return
							}
						}()

						os.Stdout, os.Stderr = devNull, devNull

						*isWorking = false
						*isPermission = false
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("Error:", err)
			}
		}
	}()

	// Добавляем файл для наблюдения
	err = watcher.Add(runningPath)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}

func openNewWindowForCommand(executeCommand structs.ExecuteCommandFuncParams) {
	var cmd *exec.Cmd

	// Преобразуем команду в формат для запуска в новом окне
	commandToExecute := strings.Join(executeCommand.CommandArgs, " ")
	dir, _ := os.Getwd()
	newOrbix := func() {
		if len(executeCommand.CommandArgs) < 1 {
			err := commands.ChangeDirectory(Absdir)
			if err != nil {
				fmt.Println("Error changing directory:", err)
			}

			commandToExecute = "go run orbix.go"
		}
	}

	// Определяем ОС и выбираем способ запуска нового окна
	switch system.OperationSystem {
	case "windows":
		newOrbix()
		// Для Windows запускаем новое окно с помощью cmd
		cmd = exec.Command("cmd", "/c", "start", "cmd", "/k", commandToExecute)
	case "linux":
		newOrbix()
		// Для Linux используем gnome-terminal, xterm или другой эмулятор терминала
		cmd = exec.Command("gnome-terminal", "--", "bash", "-c", commandToExecute)
	case "darwin":
		newOrbix()
		// Для MacOS запускаем новое окно в приложении Terminal
		cmd = exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Terminal" to do script "%s"`, commandToExecute))
	default:
		// Если ОС неизвестна, выводим ошибку
		fmt.Println("Unsupported OS")
		return
	}

	// Запускаем команду
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting command in new window:", err)
	}

	err = commands.ChangeDirectory(dir)
	if err != nil {
		fmt.Println("Error changing directory:", err)
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

func catchSyntaxErrs(execCommandCatchErrs structs.ExecuteCommandCatchErrs) (findErr bool) {
	if *execCommandCatchErrs.EchoTime && *execCommandCatchErrs.RunOnNewThread {
		fmt.Println(red("You cannot take timing and running on new thread at the same time"))
		return true
	}

	return false
}

// removeFlags удаляет части строки, если они содержатся в OrbixFlags
func removeFlags(input string) string {
	// Разделяем строку на части
	parts := strings.Fields(input)
	var result []string

	// Проходим по всем частям
	for _, part := range parts {
		// Проверяем, есть ли текущая часть в OrbixFlags
		if !contains(OrbixFlags, part) {
			// Если часть не является флагом, добавляем её в результат
			result = append(result, part)
		}
	}

	// Соединяем оставшиеся части в строку и возвращаем
	return strings.Join(result, " ")
}

// contains проверяет, находится ли элемент в списке
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func commandFile(command string) bool {
	return command == "py" ||
		command == "read" ||
		command == "edit" ||
		command == "create" ||
		command == "rem" ||
		command == "rename" ||
		command == "del" ||
		command == "delete" ||
		command == "cf" ||
		command == "df" ||
		command == "rustc" ||
		command == "cl"
}

func fullFileName(commandArgs *[]string) {
	if len(*commandArgs) == 0 {
		return
	}

	if len(*commandArgs) == 1 {
		return
	}

	var fileName string

	if len(*commandArgs) > 1 {
		for _, arg := range *commandArgs {
			fileName += arg + " "
		}

		fileName = strings.TrimSpace(fileName)

		resultSlice := []string{fileName}
		*commandArgs = resultSlice
	} else {
		return
	}
}
