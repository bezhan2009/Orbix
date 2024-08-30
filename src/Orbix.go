package src

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithoutSignature/CD"
	"goCmd/cmd/dirInfo"
	"goCmd/system"
	"goCmd/utils"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	Absdir, _       = filepath.Abs("")
	DirUser, _      = filepath.Abs("")
	UsernameFromDir = dirInfo.CmdUser(DirUser)
)

func Orbix(commandInput string, echo bool) {
	Init()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	var printUserDir string

	if !echo && commandInput == "" {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Println(red("You cannot enable echo with an empty Input command!"))
		return
	}
	if echo {
		utils.SystemInformation()
	}

	isWorking := true
	isPermission := true

	if commandInput != "" {
		isPermission = false
	}

	var promptText string

	isEmpty, err := isPasswordDirectoryEmpty()
	if err != nil {
		animatedPrint("Ошибка при проверке директории с паролями:" + err.Error() + "\n")
		return
	}

	username := ""
	if !isEmpty && commandInput == "" {
		dir, _ := os.Getwd()
		user := dirInfo.CmdUser(dir)
		nameuser, isSuccess := CheckUser(user)
		if !isSuccess {
			return
		}

		username = nameuser
		// Создаем файл running.txt, если его нет
		if _, err := os.Stat("running.txt"); os.IsNotExist(err) {
			_, err := os.Create("running.txt")
			if err != nil {
				panic(err)
			}
		}

		// Проверка наличия username в файле running.txt
		runningPath := Absdir
		runningPath += "\\running.txt"
		sourceRunning, errReading := os.ReadFile(runningPath)
		if errReading == nil {
			dataRunning := string(sourceRunning)
			lines := strings.Split(dataRunning, "\n")
			found := false
			for _, line := range lines {
				if strings.TrimSpace(line) == username {
					found = true
					break
				}
			}
			if !found {
				// Добавляем username в конец файла running.txt
				file, err := os.OpenFile("running.txt", os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("Ошибка открытия файла running.txt:", err)
					return
				}
				defer file.Close()

				if _, err := file.WriteString(username + "\n"); err != nil {
					fmt.Println("Ошибка записи в файл running.txt:", err)
					return
				}
			}
		}
	}

	// Создаем канал для получения сигналов
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM) // Обработка Ctrl+C и других сигналов завершения

	// Запускаем горутину для обработки сигналов
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-signalChan
		fmt.Println("\nПолучен сигнал прерывания. Завершение программы...")
		isWorking = false
	}()

	// Сохранить оригинальные выводы
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	// Перенаправить вывод в нулевой девайс
	devNull, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0666)
	defer devNull.Close()
	//
	//if !isWorking {
	//	removeUserFromRunningFile(username)
	//}
	//
	for isWorking {
		if !echo {
			os.Stdout = devNull
			os.Stderr = devNull
		} else {
			os.Stdout = originalStdout
			os.Stderr = originalStderr
		}

		dir, _ := os.Getwd()

		if promptText == "" && echo && system.IsAdmin {
			fmt.Printf("\nORB %s>%s", dir, green(commandInput))
		}

		if !system.IsAdmin {
			runFilePath := Absdir
			runFilePath += "\\isRun.txt"

			currentBranchGit, errGitBranch := GetCurrentGitBranch()
			if errGitBranch != nil {
				currentBranchGit = "" // Ensure currentBranchGit is empty on error
			}

			// Ensure absolute path for activeUser.txt
			activeUserFilePath := DirUser
			activeUserFilePath += "\\activeUser.txt"

			_, err := os.Create(activeUserFilePath)
			if err != nil {
				fmt.Println(red(err))
			}

			err = os.WriteFile(activeUserFilePath, []byte(username), 0644)
			if err != nil {
				fmt.Println(red(err))
			}

			dirC := dirInfo.CmdDir(dir)
			user := dirInfo.CmdUser(dir)

			if username != "" {
				user = username
				printUserDir = user
			}

			currentTime := time.Now().Format("15:04")
			location := os.Getenv("CITY")
			if location == "" {
				location = "Unknown City"
			}

			// Проверка наличия username в файле running.txt
			runningPath := Absdir
			runningPath += "\\running.txt"
			sourceRunning, errReading := os.ReadFile(runningPath)
			if errReading == nil {
				dataRunning := string(sourceRunning)
				lines := strings.Split(dataRunning, "\n")
				found := false

				for _, line := range lines {
					if strings.TrimSpace(line) == username {
						found = true
						break
					}
				}

				if !found {
					fmt.Println(red("Пользователь не авторизован."))
					isWorking = false
					isPermission = false
					continue
				}
			}
			if promptText != "" && echo {
				fmt.Println("\n" + promptText)
			} else if echo {
				var gitInfo string
				if currentBranchGit != "" {
					gitInfo = fmt.Sprintf(" %s%s", yellow("git:"), green("[", currentBranchGit, "]"))
				}
				if printUserDir == "" {
					printUserDir = UsernameFromDir
				}

				fmt.Printf("\n%s%s%s%s%s%s%s%s %s%s%s%s%s%s%s%s\n",
					yellow("┌"), yellow("─"), yellow("("), cyan("Orbix@"+printUserDir), yellow(")"), yellow("─"), yellow("["),
					yellow(location), magenta(currentTime), yellow("]"), yellow("─"), yellow("["),
					cyan("~"), cyan(dirC), yellow("]"), gitInfo)
				fmt.Printf("%s%s%s %s",
					yellow("└"), yellow("─"), green("$"), green(commandInput))
			}
		}

		if promptText != "" && echo {
			fmt.Print("\n" + promptText)
		}

		var commandLine string
		var commandParts []string
		var commandArgs []string
		var commandLower string
		var command string

		if commandInput != "" {
			isWorking = false
			isPermission = false
			commandLine = strings.TrimSpace(commandInput)
			commandParts = utils.SplitCommandLine(commandLine)
			if len(commandParts) == 0 {
				continue
			}

			command = commandParts[0]
			commandArgs = commandParts[1:]
			commandLower = strings.ToLower(command)
		} else {
			commandLine = prompt.Input("", autoComplete)
			commandLine = strings.TrimSpace(commandLine)
			commandParts = utils.SplitCommandLine(commandLine)

			if len(commandParts) == 0 {
				continue
			}

			command = commandParts[0]
			commandArgs = commandParts[1:]
			commandLower = strings.ToLower(command)

			CommandHistory = append(CommandHistory, commandLine)
		}

		fmt.Println()

		num, err := strconv.Atoi(commandLower)

		if err == nil && len(commandArgs) == 0 {
			fmt.Println(num)
			continue
		}

		if commandLower == "cd" && len(commandArgs) < 1 {
			fmt.Println(dir)
			continue
		}

		if commandLower == "prompt" {
			handlePromptCommand(commandArgs, &promptText)
			continue
		}

		if commandLower == "help" {
			displayHelp()
			continue
		}

		// Handle the "cd" command and drive changes
		if commandLower == "cd" || len(command) == 2 && command[1] == ':' {
			targetDir := command
			if commandLower == "cd" && len(commandArgs) > 0 {
				targetDir = commandArgs[0]
			}

			err := CD.ChangeDirectory(targetDir)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

		isValid := utils.ValidCommand(commandLower, Commands)

		if !isValid {
			fullCommand := append([]string{command}, commandArgs...)
			err := utils.ExternalCommand(fullCommand)
			if err != nil {
				fullPath := filepath.Join(dir, command)
				fullCommand[0] = fullPath
				err = utils.ExternalCommand(fullCommand)
				if err != nil {
					suggestedCommand := suggestCommand(commandLower)
					fmt.Print(red(fmt.Sprintf("Error executing command '%s': %v\n", commandLine, err)))
					if suggestedCommand != "" {
						fmt.Print(yellow(fmt.Sprintf("Did you mean: %s?\n", suggestedCommand)))
					}
				}
			}
			continue
		}

		ExecuteCommand(commandLower, command, commandLine, dir, Commands, commandArgs, &isWorking, isPermission, username)

		if commandLower == "signout" {
			break
		}
	}

	// Восстановить оригинальные выводы
	os.Stdout = originalStdout
	os.Stderr = originalStderr
}
