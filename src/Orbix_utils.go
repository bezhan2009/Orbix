package src

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/fsnotify/fsnotify"
	"goCmd/system"
	"goCmd/utils"
	"log"
	"os"
	"path/filepath"
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

func printPromptInfo(location, user, dirC string, green, cyan, yellow, magenta func(...interface{}) string, sd *system.Session, commandInput string) {
	fmt.Printf("\n%s%s%s%s%s%s%s%s %s%s%s%s%s%s%s%s%s%s%s\n",
		yellow("┌"), yellow("─"), yellow("("), cyan("Orbix@"+getUser(user)), yellow(")"), yellow("─"), yellow("["),
		yellow(location), magenta(time.Now().Format("15:04")), yellow("]"), yellow("─"), yellow("["),
		cyan("~"), cyan(dirC), yellow("]"), yellow(" git:"), green("["), green(sd.GitBranch), green("]"))
	fmt.Printf("%s%s%s %s",
		yellow("└"), yellow("─"), green(strings.TrimSpace(Prompt)), green(commandInput))
}

func readCommandLine(commandInput string) (string, string, []string, string) {
	var commandLine string
	if commandInput != "" {
		commandLine = strings.TrimSpace(commandInput)
	} else {
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
	if strings.TrimSpace(commandLower) == "cd" {
		return true, nil
	}

	if strings.TrimSpace(commandLower) == "git" {
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
