package src

import (
	"errors"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/fsnotify/fsnotify"
	_chan "goCmd/chan"
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
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

var UnknownCommandsCounter uint
var dirC string

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
	if strings.TrimSpace(system.User) == "" {
		return system.User
	} else {
		return username
	}
}

func printPromptInfo(location, user, dirC, commandInput *string, sd *system.Session) {
	// Обрезаем Prompt, если он длинный
	if len(system.Prompt) > 2 {
		system.Prompt = string(system.Prompt[:2])
	}

	// Сохраняем форматированные данные в переменные
	gitBranch := system.Green(sd.GitBranch)
	userInfo := system.Cyan("Orbix@" + getUser(*user))
	locationInfo := system.Yellow(*location)
	dirInfo := system.Cyan(*dirC)
	currentTime := system.Magenta(time.Now().Format("15:04"))
	prompt := strings.TrimSpace(system.Prompt)
	input := strings.TrimSpace(*commandInput)

	// Формируем строки для вывода
	header := fmt.Sprintf(
		"\n%s%s%s%s%s%s%s%s %s%s%s%s%s%s%s%s%s%s%s",
		system.Yellow("╭"), system.Yellow("─"), system.Yellow("("),
		userInfo, system.Yellow(")"), system.Yellow("─"), system.Yellow("["),
		locationInfo, currentTime, system.Yellow("]"), system.Yellow("─"), system.Yellow("["),
		system.Cyan("~"), dirInfo, system.Yellow("]"), system.Yellow(" git:"), system.Yellow("["), gitBranch, system.Yellow("]"),
	)

	footer := fmt.Sprintf("%s%s %s", system.Yellow("╰"), system.Green(prompt), system.Green(*commandInput))

	// Печатаем информацию
	fmt.Println(header)
	fmt.Print(footer)

	// Выводим перенос строки, если есть команды
	if input != "" && len(os.Args) > 0 {
		fmt.Println()
	}
}

func PrintPromptInfoWithoutGit(location, user, dirC, commandInput *string) {
	// Обрезаем Prompt, если он длинный
	if len(system.Prompt) > 2 {
		system.Prompt = string(system.Prompt[:1])
	}

	// Сохраняем форматированные данные в переменные
	userInfo := system.Cyan("Orbix@" + getUser(*user))
	locationInfo := system.Yellow(*location)
	dirInfo := system.Cyan(*dirC)
	currentTime := system.Magenta(time.Now().Format("15:04"))
	prompt := strings.TrimSpace(system.Prompt)

	// Формируем строки для вывода
	header := fmt.Sprintf(
		"\n%s%s%s%s%s%s%s%s %s%s%s%s%s%s%s",
		system.Yellow("╭"), system.Yellow("─"), system.Yellow("("),
		userInfo, system.Yellow(")"), system.Yellow("─"), system.Yellow("["),
		locationInfo, currentTime, system.Yellow("]"), system.Yellow("─"), system.Yellow("["),
		system.Cyan("~"), dirInfo, system.Yellow("]"),
	)

	footer := fmt.Sprintf("%s%s %s", system.Yellow("╰"), system.Green(prompt), system.Green(*commandInput))

	// Печатаем информацию
	fmt.Println(header)
	fmt.Print(footer)

	// Выводим перенос строки, если есть команды
	if *commandInput != "" && len(os.Args) > 0 {
		fmt.Println()
	}
}

func commandVar(commandLower string) bool {
	return commandLower == "setvar" ||
		commandLower == "delvar" ||
		commandLower == "getvar"
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
	if commandVar(strings.ToLower(commandLineSplit[0])) {
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
	if strings.TrimSpace(commandLower) == "cd" && system.GitExists {
		return true
	}

	if strings.TrimSpace(commandLower) == "git" && system.GitExists {
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
					if !checkUserInRunningFile(username) && *isWorking && system.User == username {
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

func OpenNewWindowForCommand(executeCommand *structs.ExecuteCommandFuncParams) {
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
		command == "format" ||
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

func customPrompt(commandInput, prompt *string,
	colorsMap map[string]func(...interface{}) string) {
	if strings.TrimSpace(*commandInput) != "" {
		splitPrompt := strings.Split(*prompt, ", ")
		fmt.Printf("%s%s", colorsMap[splitPrompt[1]](splitPrompt[0]), system.Green(*commandInput))
	} else {
		splitPrompt := strings.Split(*prompt, ", ")
		fmt.Print(colorsMap[splitPrompt[1]](splitPrompt[0]))
	}
}

func printOldPrompt(commandInput, dir *string) {
	if strings.TrimSpace(*commandInput) != "" {
		fmt.Printf("ORB %s>%s", *dir, system.Green(*commandInput))
	} else {
		fmt.Printf("ORB %s>", *dir)
	}
}

func OrbixPrompt(session *system.Session,
	prompt, commandInput *string,
	isWorking, isPermission *bool,
	colorsMap *map[string]func(...interface{}) string) {
	if session.IsAdmin {
		if *prompt == "" {
			printOldPrompt(commandInput, &system.UserDir)
		} else {
			customPrompt(commandInput, prompt,
				*colorsMap)
		}

		return
	}

	Orbixuser, _ := environment.GetVariableValue("user")
	if Orbixuser == "" {
		Orbixuser = dirInfo.CmdUser(&system.UserDir)
	}

	OrbixuserStr := fmt.Sprintf("%s", Orbixuser)

	if !session.IsAdmin {
		dirC = dirInfo.CmdDir(system.UserDir)

		// Single user check outside repeated prompt formatting
		if !system.Unauthorized {
			go func() {
				watchFile(system.RunningPath, OrbixuserStr, isWorking, isPermission)
			}()
		}

		if *prompt == "" {
			if system.GitExists {
				printPromptInfo(&system.Location,
					&OrbixuserStr,
					&dirC,
					commandInput,
					session) // New helper function
			} else {
				PrintPromptInfoWithoutGit(&system.Location,
					&OrbixuserStr,
					&dirC,
					commandInput) // New helper function
			}
		} else {
			customPrompt(commandInput, prompt,
				*colorsMap)
		}
	}
}

func printInfo(s interface{}, echo bool) {
	if !echo {
		return
	}

	fmt.Print(s)
}

func LoadConfigs(echo bool) error {
	printInfo(system.Cyan("Loading configs"), echo)
	if echo {
		utils.AnimatedPrint("...\n", "cyan")
	}

	_chan.LoadConfigsFn = environment.LoadUserConfigs

	err := environment.LoadUserConfigs()
	if err != nil {
		printInfo(system.Red("Error Loading configs:", err), echo)
		println()
	} else {
		printInfo(system.Green("Successfully Loaded configs"), echo)
		println()
	}

	return err
}

func InitSession(prefix *string,
	rebooted structs.RebootedData,
	OrbixLoopData structs.OrbixLoopData) *system.Session {
	system.CntLaunchedOrbixes++

	dirC = dirInfo.CmdDir(system.UserDir)

	if rebooted.Prefix != "" {
		*prefix = rebooted.Prefix
	} else {
		*prefix = OrbixLoopData.SessionData.NewSessionData(
			OrbixLoopData.SessionData.Path,
			OrbixLoopData.SessionData.User,
			OrbixLoopData.SessionData.GitBranch,
			OrbixLoopData.SessionData.IsAdmin,
		)
	}

	session, exists := OrbixLoopData.SessionData.GetSession(*prefix)
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
	go system.InitSession(OrbixLoopData.Username,
		session)

	session.PreviousPath = system.PreviousSessionPath
	if system.PreviousSessionPrefix != "" {
		session, _ = OrbixLoopData.SessionData.GetSession(system.PreviousSessionPrefix)
	}

	system.GlobalSession = *session

	system.Path = system.UserDir

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
		OrbixUser := dirInfo.CmdUser(&dir)

		nameUser, isSuccess, errUser := user.CheckUser(OrbixUser, sessionData)
		if !isSuccess {
			if _chan.UserStatusAuth {
				system.Unauthorized = false
				_chan.UpdateChan("system__user")
			} else {
				system.Unauthorized = true
			}

			return "", errUser
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

func ignoreSI(signalChan chan os.Signal,
	sessionData *system.AppState,
	prompt, commandInput, username *string) bool {
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

				dirC = dirInfo.CmdDir(dir)
				userName := sessionData.User
				if userName == "" {
					userName = dirInfo.CmdUser(&dir)
				}

				if *username != "" {
					userName = *username
				}

				fmt.Println()
				if system.ExecutingCommand {
					return true
				}

				if *prompt == "" {
					if system.GitExists {
						gitBranch, _ := system.GetCurrentGitBranch()
						printPromptInfo(&system.Location, &userName, &dirC, commandInput, &system.Session{Path: dir, GitBranch: gitBranch})
					} else {
						PrintPromptInfoWithoutGit(&system.Location, &userName, &dirC, commandInput)
					}
				} else {
					customPrompt(commandInput, prompt,
						colorsMap)
				}
			} else {
				dir, _ := os.Getwd()
				if *prompt == "" {
					fmt.Printf("ORB %s>%s", dir, system.Green(*commandInput))
				} else {
					customPrompt(commandInput, prompt,
						colorsMap)
				}
			}
		}
	}

	return false
}

func IgnoreSiC(commandInput, prompt *string,
	OrbixLoopData *structs.OrbixLoopData) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if ignoreSI(signalChan,
			OrbixLoopData.SessionData,
			prompt, commandInput, &OrbixLoopData.Username) {
			return
		}
	}()
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

func UsingForLT(commandInput string) bool {
	if strings.TrimSpace(commandInput) != "" && strings.TrimSpace(commandInput) != "restart" {
		return true
	}

	return false
}

func OrbixUser(commandInput string,
	echo bool,
	rebooted *structs.RebootedData,
	SD *system.AppState,
	ExecLtCommand func(commandInput string)) (LoopData structs.OrbixLoopData, LoadUserConfigsFn func(echo bool) error) {
	LoadUserConfigsFn = LoadConfigs

	if UsingForLT(commandInput) {

		// Load User Configs
		_ = LoadConfigs(false)

		ExecLtCommand(commandInput)

		isWorking := false
		isPermission := false
		RestartAfterInit := false

		return structs.OrbixLoopData{
			IsWorking:        &isWorking,
			IsPermission:     &isPermission,
			Username:         "",
			SessionData:      &system.AppState{},
			RestartAfterInit: &RestartAfterInit,
		}, LoadUserConfigsFn
	}

	RestartAfterInit := false

	sessionData := InitOrbixFn(&RestartAfterInit,
		echo,
		commandInput,
		*rebooted,
		SD)

	isWorking := true
	isPermission := true
	if commandInput != "" {
		isPermission = false
	}

	username, err := DefineUser(commandInput,
		*rebooted,
		sessionData)
	if err != nil {
		isWorking = false
		isPermission = false
		RestartAfterInit = false

		return structs.OrbixLoopData{
			IsWorking:        &isWorking,
			IsPermission:     &isPermission,
			Username:         "",
			SessionData:      &system.AppState{},
			RestartAfterInit: &RestartAfterInit,
		}, LoadUserConfigsFn
	}

	// Load User Configs
	_ = LoadConfigs(true)

	if username != "" {
		system.EditableVars["user"] = &username
	}

	return structs.OrbixLoopData{
		IsWorking:        &isWorking,
		IsPermission:     &isPermission,
		Username:         username,
		SessionData:      sessionData,
		RestartAfterInit: &RestartAfterInit,
	}, nil
}

func EdgeCases(OrbixLoopData structs.OrbixLoopData,
	rebooted structs.RebootedData,
	RecoverAndRestore func(rebooted *structs.RebootedData)) {
	if len(OrbixLoopData.Session.CommandHistory) < 10 {
		go system.InitSession(OrbixLoopData.Username,
			OrbixLoopData.Session)
	}

	if system.RebootAttempts != 0 {
		RecoverAndRestore(&rebooted)
		system.RebootAttempts = 0
	}
}

func PrepareOrbix() {
	_chan.User = system.User
	_chan.UserName = system.UserName
	_chan.DirUser, _ = os.Getwd()
}

func RestoreOrbix() {
	system.User = _chan.User
	system.UserName = _chan.UserName

	err := commands.ChangeDirectory(_chan.DirUser)
	if err != nil {
		fmt.Println(system.Red("Error changing directory:", err))
	}

	system.UserDir, _ = os.Getwd()
}
