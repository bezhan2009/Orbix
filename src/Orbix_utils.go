package src

import (
	"errors"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/fsnotify/fsnotify"
	"goCmd/cmd/commands"
	"goCmd/cmd/dirInfo"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var UnknownCommandsCounter uint

// New helper functions
func initializeRunningFile(username string) {
	// Check and initialize running.txt if not exists
	if _, err := os.Stat(system.OrbixRunningUsersFileName); os.IsNotExist(err) {
		if _, err = os.Create(system.OrbixRunningUsersFileName); err != nil {
			panic(err)
		}
	}

	// Check for username in running.txt and add if missing
	runningPath := filepath.Join(Absdir, system.OrbixRunningUsersFileName)
	if sourceRunning, err := os.ReadFile(runningPath); err == nil {
		if !strings.Contains(string(sourceRunning), username) {
			if file, err := os.OpenFile(system.OrbixRunningUsersFileName, os.O_APPEND|os.O_WRONLY, 0644); err == nil {
				defer func() {
					err = file.Close()
					if err != nil {
						return
					}
				}()
				if _, err := file.WriteString("\n" + username + "\n"); err != nil {
					fmt.Println(fmt.Sprintf("Error writing to %s: %s", system.OrbixRunningUsersFileName, err))
				}
			}
		}
	}
}

func checkUserInRunningFile(username string) bool {
	runningPath := filepath.Join(Absdir, system.OrbixRunningUsersFileName)
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
	if len(Prompt) > 2 {
		Prompt = string(Prompt[0:2])
	}

	fmt.Printf("\n%s%s%s%s%s%s%s%s %s%s%s%s%s%s%s%s%s%s%s\n",
		yellow("╭"), yellow("─"), yellow("("), cyan("Orbix@"+getUser(user)), yellow(")"), yellow("─"), yellow("["),
		yellow(location), magenta(time.Now().Format("15:04")), yellow("]"), yellow("─"), yellow("["),
		cyan("~"), cyan(dirC), yellow("]"), yellow(" git:"), green("["), green(sd.GitBranch), green("]"))
	fmt.Printf("%s%s %s",
		yellow("╰"), green(strings.TrimSpace(Prompt)), green(commandInput))

	if strings.TrimSpace(commandInput) != "" && len(os.Args) > 0 {
		fmt.Println()
	}
}

func printPromptInfoWithoutGit(location, user, dirC, commandInput string) {
	if len(Prompt) > 2 {
		Prompt = string(Prompt[0])
	}

	fmt.Printf("\n%s%s%s%s%s%s%s%s %s%s%s%s%s%s%s\n",
		yellow("╭"), yellow("─"), yellow("("), cyan("Orbix@"+getUser(user)), yellow(")"), yellow("─"), yellow("["),
		yellow(location), magenta(time.Now().Format("15:04")), yellow("]"), yellow("─"), yellow("["),
		cyan("~"), cyan(dirC), yellow("]"))
	fmt.Printf("%s%s %s",
		yellow("╰"), green(strings.TrimSpace(Prompt)), green(commandInput))

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

	commandLineSplit := strings.Split(commandLine, " ")
	if commandLineSplit[0] == "setvar" || commandLineSplit[0] == "delvar" || commandLineSplit[0] == "getvar" {
		command := commandLineSplit[:1]

		return commandLine, command[0], commandLineSplit[1:], strings.ToLower(commandLineSplit[0])
	}

	commandParts := utils.SplitCommandLine(commandLine)
	if len(commandParts) == 0 {
		return "", "", nil, ""
	}

	command := commandParts[:1]

	if strings.ToLower(commandParts[0]) == "cd" {
		system.UserDir, _ = os.Getwd()
	}

	return commandLine, command[0], commandParts[1:], strings.ToLower(commandParts[0])
}

func processCommand(commandLower string) (bool, error) {
	if strings.TrimSpace(commandLower) == "cd" && GitCheck {
		return true, nil
	}

	if strings.TrimSpace(commandLower) == "git" && GitCheck {
		return true, nil
	}

	if strings.TrimSpace(commandLower) == "signout" {
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

func catchSyntaxErrs(execCommandCatchErrs structs.ExecuteCommandCatchErrs) (findErr bool) {
	if *execCommandCatchErrs.EchoTime && *execCommandCatchErrs.RunOnNewThread && !(execCommandCatchErrs.CommandLower == "orbix") {
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
		command == "del_var" ||
		command == "del" ||
		command == "gocode" ||
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

func RecoverFromThePanic(commandInput string,
	r any,
	echo bool,
	SD *system.AppState) {
	PanicText := fmt.Sprintf("Panic recovered: %v", r)
	fmt.Printf("\n%s\n", red(PanicText))

	if RebootAttempts > system.MaxRetryAttempts {
		fmt.Println(red("Max retry attempts reached. Exiting..."))
		log.Println("Max retry attempts reached. Exiting...")
		os.Exit(1)
	}

	RebootAttempts += 1

	fmt.Println(yellow("Recovering from panic"))

	log.Printf("Panic recovered: %v", r)

	var reboot = structs.RebootedData{
		Username: system.UserName,
		Recover:  r,
		Prefix:   Prefix,
	}

	Orbix(strings.TrimSpace(commandInput),
		echo,
		reboot,
		SD)
}

func OrbixPrompt(session *system.Session, prompt, dir, username, commandInput string, isWorking, isPermission bool, colorsMap map[string]func(...interface{}) string) {
	if session.IsAdmin {
		if prompt == "" {
			fmt.Printf("ORB %s>%s", dir, green(commandInput))
		} else {
			splitPrompt := strings.Split(prompt, ", ")
			fmt.Print(colorsMap[splitPrompt[1]](splitPrompt[0]))
		}
	}

	dirC := dirInfo.CmdDir(dir)
	user := session.User
	if user == "" {
		user = dirInfo.CmdUser(dir)
	}

	if username != "" {
		user = username
	}

	if !session.IsAdmin {
		// Single user check outside repeated prompt formatting
		if !Unauthorized {
			go func() {
				watchFile(RunningPath, user, &isWorking, &isPermission)
			}()
		}

		if prompt == "" {
			if GitCheck {
				printPromptInfo(Location, user, dirC, commandInput, session) // New helper function
			} else {
				printPromptInfoWithoutGit(Location, user, dirC, commandInput) // New helper function
			}
		} else {
			splitPrompt := strings.Split(prompt, ", ")
			fmt.Print(colorsMap[splitPrompt[1]](splitPrompt[0]))
		}
	}
}

func InitSession(prefix *string,
	rebooted structs.RebootedData,
	sessionData *system.AppState) *system.Session {
	if rebooted.Prefix != "" {
		*prefix = rebooted.Prefix
	} else {
		*prefix = sessionData.NewSessionData(sessionData.Path, sessionData.User, sessionData.GitBranch, sessionData.IsAdmin)
	}

	session, exists := sessionData.GetSession(*prefix)
	if !exists {
		fmt.Println(red("Session does not exist!"))
		return nil
	}

	if session == nil {
		fmt.Println(red("Session is nil!"))
		return nil
	}

	Prefix = fmt.Sprintf(*prefix)

	// Initialize Global Vars
	go Init(session)

	// Load User Configs
	fmt.Print(cyan("Loading configs"))
	utils.AnimatedPrint("...\n", "cyan")

	err := LoadUserConfigs()
	if err != nil {
		fmt.Println(red("Error Loading configs:", err))
	} else {
		fmt.Println(green("Successfully Loaded configs"))
	}

	session.PreviousPath = PreviousSessionPath
	fmt.Println(green(session.PreviousPath))
	if PreviousSessionPrefix != "" {
		session, _ = sessionData.GetSession(PreviousSessionPrefix)
	}

	GlobalSession = *session

	dir, _ := os.Getwd()
	system.Path = dir

	return session
}

func defineUser(commandInput string,
	rebooted structs.RebootedData,
	sessionData *system.AppState) (string, error) {
	var username string

	// Check if password directory is empty once and handle errors here
	isEmpty, err := isPasswordDirectoryEmpty()
	if err != nil {
		animatedPrint(fmt.Sprintf("Error checking password directory: %s\n", err.Error()), "red")
		return "", errors.New("ErrCheckPasswordDirectory")
	}

	if strings.TrimSpace(rebooted.Username) != "" {
		username = strings.TrimSpace(rebooted.Username)
	} else if !isEmpty && commandInput == "" {
		dir, _ := os.Getwd()
		user := dirInfo.CmdUser(dir)

		nameUser, isSuccess := CheckUser(user, sessionData)
		if !isSuccess {
			return "", errors.New("ErrSuccess")
		}
		Unauthorized = false
		username = nameUser
		if username != user {
			initializeRunningFile(username)
		}

		if user == username {
			sessionData.IsAdmin = true
			sessionData.User = user
		} else {
			sessionData.IsAdmin = false
			sessionData.User = username
		}
	}

	return username, nil
}

func ignoreSI(signalChan chan os.Signal,
	sessionData *system.AppState,
	prompt, commandInput, username string) bool {
	colorsMap := system.GetColorsMap()
	if SessionsStarted > 1 {
		return true
	}

	for {
		sig := <-signalChan

		if sig == syscall.SIGHUP {
			DeleteUserFromRunningFile(system.UserName)
			os.Exit(1)
		}

		if !ExecutingCommand {
			fmt.Println(red("^C"))
			if !GlobalSession.IsAdmin {
				dir, _ := os.Getwd()

				dirC := dirInfo.CmdDir(dir)
				user := sessionData.User
				if user == "" {
					user = dirInfo.CmdUser(dir)
				}

				if username != "" {
					user = username
				}

				fmt.Println()
				if prompt == "" {
					if GitCheck {
						gitBranch, _ := GetCurrentGitBranch()
						printPromptInfo(Location, user, dirC, commandInput, &system.Session{Path: dir, GitBranch: gitBranch})
					} else {
						printPromptInfoWithoutGit(Location, user, dirC, commandInput)
					}
				} else {
					splitPrompt := strings.Split(prompt, ", ")
					fmt.Print(colorsMap[splitPrompt[1]](splitPrompt[0]))
				}
			} else {
				dir, _ := os.Getwd()
				if prompt == "" {
					fmt.Printf("ORB %s>%s", dir, green(commandInput))
				} else {
					splitPrompt := strings.Split(prompt, ", ")
					fmt.Print(colorsMap[splitPrompt[1]](splitPrompt[0]))
				}
			}
		}
	}

	return false
}

func setLocation() {
	if strings.TrimSpace(Location) == "" {
		Location = os.Getenv("CITY")
		if strings.TrimSpace(Location) == "" {
			Location = string(strings.TrimSpace(os.Getenv("USERS_LOCATION")))
		}
	}
}

func initOrbixFn(RestartAfterInit *bool,
	echo bool,
	commandInput string,
	rebooted structs.RebootedData,
	SD *system.AppState) *system.AppState {
	Prompt = string(strings.TrimSpace(os.Getenv("PROMPT")))
	SessionsStarted = SessionsStarted + 1

	setLocation()

	// Initialize colors
	InitColors()

	if strings.TrimSpace(strings.ToLower(system.OperationSystem)) == "windows" {
		Commands = append(Commands, structs.Command{Name: "neofetch", Description: "Displays information about the system"})
		AdditionalCommands = append(AdditionalCommands, structs.Command{Name: "neofetch", Description: "Displays information about the system"})
	}

	if RebootAttempts > 5 {
		system.OrbixWorking = false
		fmt.Println(red("Max retry attempts reached. Exiting..."))
		log.Println("Max retry attempts reached. Exiting...")
		return nil
	}

	system.OrbixWorking = true

	if strings.TrimSpace(commandInput) == "restart" {
		*RestartAfterInit = true
	}

	if SD == nil {
		fmt.Println(red("Fatal: App State is nil!"))
		os.Exit(1)
	}

	if err := commands.ChangeDirectory(Absdir); err != nil {
		fmt.Println(red(err))
	}

	sessionData := SD

	if !echo && commandInput == "" {
		fmt.Println(red("You cannot enable echo with an empty Input command!"))
		return nil
	}

	if echo && rebooted.Username == "" && commandInput == "" {
		SystemInformation()
	}

	return sessionData
}

func restartAfterInit(SD *system.AppState,
	sessionData *system.AppState,
	rebooted structs.RebootedData,
	prefix,
	username string,
	echo bool) {
	SD.User = username
	SD.IsAdmin = sessionData.IsAdmin
	rebooted.Prefix = prefix
	if len(os.Args) > 1 {
		return
	}

	Orbix("",
		echo,
		rebooted,
		SD)
}
