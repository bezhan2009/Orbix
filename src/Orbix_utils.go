package src

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	_chan "goCmd/chan"
	"goCmd/cmd/commands"
	"goCmd/cmd/dirInfo"
	"goCmd/src/environment"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var UnknownCommandsCounter uint
var dirC string

func getUser(username string) string {
	if strings.TrimSpace(system.User) == "" {
		return system.User
	} else {
		return username
	}
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
					if !checkUserInRunningFile(username) && *isWorking && system.User == username && !system.OrbixRecovering {
						time.Sleep(system.RetryDelay)
						if !checkUserInRunningFile(username) && *isWorking && system.User == username && !system.OrbixRecovering {

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
							return
						}
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

func UsingForLT(commandInput string) bool {
	if strings.TrimSpace(commandInput) != "" && strings.TrimSpace(commandInput) != "restart" {
		return true
	}

	return false
}

func EdgeCases(OrbixLoopData *structs.OrbixLoopData,
	session *system.Session,
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

	if strings.TrimSpace(OrbixLoopData.Username) == "" {
		nickname := GetUserNickname()
		system.User = nickname
		session.User = nickname
		OrbixLoopData.Username = nickname
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
