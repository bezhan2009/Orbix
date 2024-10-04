package src

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"goCmd/system"
	"goCmd/utils"
	"os"
	"path/filepath"
	"runtime"
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
				if _, err := file.WriteString(username + "\n"); err != nil {
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

func printPromptInfo(location, user, dirC string, green, cyan, yellow, magenta func(...interface{}) string, sd *system.Session, commandInput string) {
	fmt.Printf("\n%s%s%s%s%s%s%s%s %s%s%s%s%s%s%s%s%s%s%s\n",
		yellow("┌"), yellow("─"), yellow("("), cyan("Orbix@"+user), yellow(")"), yellow("─"), yellow("["),
		yellow(location), magenta(time.Now().Format("15:04")), yellow("]"), yellow("─"), yellow("["),
		cyan("~"), cyan(dirC), yellow("]"), yellow(" git:"), green("["), green(sd.GitBranch), green("]"))
	fmt.Printf("%s%s%s %s",
		yellow("└"), yellow("─"), green("$"), green(commandInput))
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

func processCommand(commandLower string, commandArgs []string, sd *system.Session) error {
	if commandLower == "cd" {
		SetGitBranch(sd)
		return nil
	}

	if commandLower == "git" && len(commandArgs) > 2 {
		if commandArgs[0] == "switch" {
			SetGitBranch(sd)
		}
	}

	if commandLower == "signout" {
		return fmt.Errorf("signout")
	}

	return nil
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

func fetchNeofetch() {
	// Получаем информацию о пользователе
	username := GlobalSession.User

	// Получаем информацию о системе
	osName := fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
	hostStat, _ := host.Info()
	uptime := fmt.Sprintf("%v", time.Duration(hostStat.Uptime)*time.Second)
	cpuInfo, _ := cpu.Info()
	cpuModel := cpuInfo[0].ModelName
	memStat, _ := mem.VirtualMemory()
	memory := fmt.Sprintf("%.2fMiB / %.2fMiB", float64(memStat.Used)/1024/1024, float64(memStat.Total)/1024/1024)

	// Проверяем переменные окружения и устанавливаем значения по умолчанию
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash" // Или другой шела по умолчанию
	}
	terminal := os.Getenv("TERM")
	if terminal == "" {
		terminal = "unknown" // Или определить другой терминал по умолчанию
	}

	// Цветное оформление
	title := color.New(color.FgCyan, color.Bold).SprintFunc()
	info := color.New(color.FgWhite).SprintFunc()
	value := color.New(color.FgYellow).SprintFunc()

	neofetch := fmt.Sprintf(`%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s`,
		cyan(fmt.Sprintf("       /0d/+00sssooyPONyssss/   /0-/+oosss0oyMMNyssssD           %s@%s", title(username), title("user"))),
		cyan(fmt.Sprintf("       /ef/+oossso00ygSNySss/   /0-/+F0sssooyMMNyssssG           %s", info("-----------"))),
		cyan(fmt.Sprintf("       /9-/+oosss0yGGNyssss-+   /0-/+oo00sooyMMNyssss+           OS: %s", value(osName))),
		cyan(fmt.Sprintf("       /P=/+00sss0oySFNyssss/   /0-/+D0sssooyMMNyssss\\           Kernel: %s", value(hostStat.KernelVersion))),
		cyan(fmt.Sprintf("       /6=/+oosso0oyADNyssgs/   /0-/+oSOssoayDFNyssss-           Terminal: %s", value(terminal))),
		cyan(fmt.Sprintf("       /+Y/+oosssooyLDNyssss/   /0-/+ooFssooyMMNyssss/           Shell: %s", value(shell))),
		cyan(fmt.Sprintf("                                                                 Uptime: %s", value(uptime))),
		cyan(fmt.Sprintf("       /0-/+ooss=aoyMMNyssasa   /0-/+oossso+yOOOyssss/           CPU: %s", value(cpuModel))),
		cyan(fmt.Sprintf("       /6=/+oosso0oyADNyssgs/   /0-/+oSOssoayDFNyssss-           Memory: %s", value(memory))),
		cyan("       /0-/+ooss--oyMMNysfss/   /0-/+odsssfsyMMNyssss   "),
		cyan("       /0-/+oosssooyMMNydsss+   /0-/+ofsss+-yMMNyssss/   "),
		cyan("       /0-/+oosssooyMMNyssa-+   /0-/+oosss(0yMMNyssss/   "),
		cyan("       /0-/+oosssooyMMNyssds=   /0-/+oosssooyMMNyssss/"))

	fmt.Println(neofetch)
}
