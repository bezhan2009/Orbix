package src

import (
	"errors"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/fsnotify/fsnotify"
	"goCmd/cmd/commands"
	"goCmd/cmd/dirInfo"
	"goCmd/src/environment"
	"goCmd/src/service"
	"goCmd/src/user"
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
	runningPath := filepath.Join(system.Absdir, system.OrbixRunningUsersFileName)
	if sourceRunning, err := os.ReadFile(runningPath); err == nil {
		if !strings.Contains(string(sourceRunning), username) {
			if file, err := os.OpenFile(system.OrbixRunningUsersFileName, os.O_APPEND|os.O_WRONLY,
				0644); err == nil {
				defer func() {
					err = file.Close()
					if err != nil {
						return
					}
				}()
				if _, err := file.WriteString("\n" + username + "\n"); err != nil {
					fmt.Println(fmt.Sprintf("Error writing to %s: %s",
						system.OrbixRunningUsersFileName, err))
				}
			}
		}
	}
}

func checkUserInRunningFile(username string) bool {
	runningPath := filepath.Join(system.Absdir, system.OrbixRunningUsersFileName)
	sourceRunning, err := os.ReadFile(runningPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(sourceRunning), username)
}

func getUser(username string) string {
	if strings.TrimSpace(system.User) != "" {
		return system.User
	} else {
		return username
	}
}

func printPromptInfo(location, user, dirC, commandInput string, sd *system.Session) {
	if len(system.Prompt) > 2 {
		system.Prompt = string(system.Prompt[0:2])
	}

	fmt.Printf("\n%s%s%s%s%s%s%s%s %s%s%s%s%s%s%s%s%s%s%s\n",
		system.Yellow("╭"), system.Yellow("─"), system.Yellow("("), system.Cyan("Orbix@"+getUser(user)), system.Yellow(")"), system.Yellow("─"), system.Yellow("["),
		system.Yellow(location), system.Magenta(time.Now().Format("15:04")), system.Yellow("]"), system.Yellow("─"), system.Yellow("["),
		system.Cyan("~"), system.Cyan(dirC), system.Yellow("]"), system.Yellow(" git:"), system.Green("["), system.Green(sd.GitBranch), system.Green("]"))
	fmt.Printf("%s%s %s",
		system.Yellow("╰"), system.Green(strings.TrimSpace(system.Prompt)), system.Green(commandInput))

	if strings.TrimSpace(commandInput) != "" && len(os.Args) > 0 {
		fmt.Println()
	}
}

func PrintPromptInfoWithoutGit(location, user, dirC, commandInput string) {
	if len(system.Prompt) > 2 {
		system.Prompt = string(system.Prompt[0])
	}

	fmt.Printf("\n%s%s%s%s%s%s%s%s %s%s%s%s%s%s%s\n",
		system.Yellow("╭"), system.Yellow("─"), system.Yellow("("), system.Cyan("Orbix@"+getUser(user)), system.Yellow(")"), system.Yellow("─"), system.Yellow("["),
		system.Yellow(location), system.Magenta(time.Now().Format("15:04")), system.Yellow("]"), system.Yellow("─"), system.Yellow("["),
		system.Cyan("~"), system.Cyan(dirC), system.Yellow("]"))
	fmt.Printf("%s%s %s",
		system.Yellow("╰"), system.Green(strings.TrimSpace(system.Prompt)), system.Green(commandInput))

	if strings.TrimSpace(commandInput) != "" && len(os.Args) > 0 {
		fmt.Println()
	}
}

func ReadCommandLine(commandInput string) (string, string, []string, string) {
	var commandLine string
	if commandInput != "" {
		commandLine = strings.TrimSpace(commandInput)
	} else {
		// Чтение ввода
		commandLine = strings.TrimSpace(prompt.Input("", service.AutoComplete))
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

	return commandLine, command[0], commandParts[1:], strings.ToLower(commandParts[0])
}

func ProcessCommand(commandLower string) bool {
	if strings.TrimSpace(commandLower) == "cd" && system.GitCheck {
		return true
	}

	if strings.TrimSpace(commandLower) == "git" && system.GitCheck {
		return true
	}

	return false
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
						fmt.Print(system.Red("\nUser not authorized. to continue, press Enter:"))
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

func OpenNewWindowForCommand(executeCommand structs.ExecuteCommandFuncParams) {
	var cmd *exec.Cmd

	// Преобразуем команду в формат для запуска в новом окне
	commandToExecute := strings.Join(executeCommand.CommandArgs, " ")
	dir, _ := os.Getwd()
	newOrbix := func() {
		if len(executeCommand.CommandArgs) < 1 {
			err := commands.ChangeDirectory(system.Absdir)
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

func CatchSyntaxErrs(execCommandCatchErrs structs.ExecuteCommandCatchErrs) (findErr bool) {
	if *execCommandCatchErrs.EchoTime && *execCommandCatchErrs.RunOnNewThread && !(execCommandCatchErrs.CommandLower == "orbix") {
		fmt.Println(system.Red("You cannot take timing and running on new thread at the same time"))
		return true
	}

	return false
}

// RemoveFlags удаляет части строки, если они содержатся в OrbixFlags
func RemoveFlags(input string) string {
	// Разделяем строку на части
	parts := strings.Fields(input)
	var result []string

	// Проходим по всем частям
	for _, part := range parts {
		// Проверяем, есть ли текущая часть в OrbixFlags
		if !contains(system.Flags, part) {
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

func CommandFile(command string) bool {
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

func FullFileName(commandArgs *[]string) {
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

func OrbixPrompt(session *system.Session, prompt, dir, username, commandInput string, isWorking, isPermission bool, colorsMap map[string]func(...interface{}) string) {
	if session.IsAdmin {
		if prompt == "" {
			fmt.Printf("ORB %s>%s", dir, system.Green(commandInput))
		} else {
			splitPrompt := strings.Split(prompt, ", ")
			fmt.Print(colorsMap[splitPrompt[1]](splitPrompt[0]))
		}
	}

	dirC := dirInfo.CmdDir(dir)
	Orbixuser := session.User
	if Orbixuser == "" {
		Orbixuser = dirInfo.CmdUser(dir)
	}

	if username != "" {
		Orbixuser = username
	}

	if !session.IsAdmin {
		// Single Orbixuser check outside repeated prompt formatting
		if !system.Unauthorized {
			go func() {
				watchFile(system.RunningPath, Orbixuser, &isWorking, &isPermission)
			}()
		}

		if prompt == "" {
			if system.GitCheck {
				printPromptInfo(system.Location, Orbixuser, dirC, commandInput, session) // New helper function
			} else {
				PrintPromptInfoWithoutGit(system.Location, Orbixuser, dirC, commandInput) // New helper function
			}
		} else {
			splitPrompt := strings.Split(prompt, ", ")
			fmt.Print(colorsMap[splitPrompt[1]](splitPrompt[0]))
		}
	}
}

func LoadConfigs() {
	fmt.Print(system.Cyan("Loading configs"))
	utils.AnimatedPrint("...\n", "cyan")

	err := environment.LoadUserConfigs()
	if err != nil {
		fmt.Println(system.Red("Error Loading configs:", err))
	} else {
		fmt.Println(system.Green("Successfully Loaded configs"))
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
		fmt.Println(system.Red("Session does not exist!"))
		return nil
	}

	if session == nil {
		fmt.Println(system.Red("Session is nil!"))
		return nil
	}

	system.Prefix = fmt.Sprintf(*prefix)

	// Initialize Global Vars
	go system.InitSession(session)

	session.PreviousPath = system.PreviousSessionPath
	fmt.Println(system.Green(session.PreviousPath))
	if system.PreviousSessionPrefix != "" {
		session, _ = sessionData.GetSession(system.PreviousSessionPrefix)
	}

	system.GlobalSession = *session

	dir, _ := os.Getwd()
	system.Path = dir

	return session
}

func DefineUser(commandInput string,
	rebooted structs.RebootedData,
	sessionData *system.AppState) (string, error) {
	var username string

	// Check if password directory is empty once and handle errors here
	isEmpty, err := user.IsPasswordDirectoryEmpty()
	if err != nil {
		service.AnimatedPrint(fmt.Sprintf("Error checking password directory: %s\n", err.Error()), "red")
		return "", errors.New("ErrCheckPasswordDirectory")
	}

	if strings.TrimSpace(rebooted.Username) != "" {
		username = strings.TrimSpace(rebooted.Username)
	} else if !isEmpty && commandInput == "" {
		dir, _ := os.Getwd()
		OrbixUser := dirInfo.CmdUser(dir)

		nameUser, isSuccess := user.CheckUser(OrbixUser, sessionData)
		if !isSuccess {
			return "", errors.New("ErrSuccess")
		}
		system.Unauthorized = false
		username = nameUser
		if username != OrbixUser {
			initializeRunningFile(username)
		}

		if OrbixUser == username {
			sessionData.IsAdmin = true
			sessionData.User = OrbixUser
		} else {
			sessionData.IsAdmin = false
			sessionData.User = username
		}
	}

	return username, nil
}

func IgnoreSI(signalChan chan os.Signal,
	sessionData *system.AppState,
	prompt, commandInput, username string) bool {
	colorsMap := system.GetColorsMap()
	if system.SessionsStarted > 1 {
		return true
	}

	for {
		sig := <-signalChan

		if sig == syscall.SIGHUP {
			user.DeleteUserFromRunningFile(system.UserName)
			os.Exit(1)
		}

		if !system.ExecutingCommand {
			fmt.Println(system.Red("^C"))
			if !system.GlobalSession.IsAdmin {
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
					if system.GitCheck {
						gitBranch, _ := system.GetCurrentGitBranch()
						printPromptInfo(system.Location, user, dirC, commandInput, &system.Session{Path: dir, GitBranch: gitBranch})
					} else {
						PrintPromptInfoWithoutGit(system.Location, user, dirC, commandInput)
					}
				} else {
					splitPrompt := strings.Split(prompt, ", ")
					fmt.Print(colorsMap[splitPrompt[1]](splitPrompt[0]))
				}
			} else {
				dir, _ := os.Getwd()
				if prompt == "" {
					fmt.Printf("ORB %s>%s", dir, system.Green(commandInput))
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
	if strings.TrimSpace(system.Location) == "" {
		system.Location = os.Getenv("CITY")
		if strings.TrimSpace(system.Location) == "" {
			system.Location = string(strings.TrimSpace(os.Getenv("USERS_LOCATION")))
		}
	}
}

func InitOrbixFn(RestartAfterInit *bool,
	echo bool,
	commandInput string,
	rebooted structs.RebootedData,
	SD *system.AppState) *system.AppState {
	system.Prompt = string(strings.TrimSpace(os.Getenv("PROMPT")))
	system.SessionsStarted = system.SessionsStarted + 1

	setLocation()

	// Initialize colors
	system.InitColors()

	if strings.TrimSpace(strings.ToLower(system.OperationSystem)) == "windows" {
		system.Commands = append(system.Commands, system.Command{Name: "neofetch", Description: "Displays information about the system"})
		system.AdditionalCommands = append(system.AdditionalCommands, system.Command{Name: "neofetch", Description: "Displays information about the system"})
	}

	if system.RebootAttempts > 5 {
		system.OrbixWorking = false
		fmt.Println(system.Red("Max retry attempts reached. Exiting..."))
		log.Println("Max retry attempts reached. Exiting...")
		return nil
	}

	system.OrbixWorking = true

	if strings.TrimSpace(commandInput) == "restart" {
		*RestartAfterInit = true
	}

	if SD == nil {
		fmt.Println(system.Red("Fatal: App State is nil!"))
		os.Exit(1)
	}

	if err := commands.ChangeDirectory(system.Absdir); err != nil {
		fmt.Println(system.Red(err))
	}

	sessionData := SD

	if !echo && commandInput == "" {
		fmt.Println(system.Red("You cannot enable echo with an empty Input command!"))
		return nil
	}

	if echo && rebooted.Username == "" && commandInput == "" {
		environment.SystemInformation()
	}

	return sessionData
}
