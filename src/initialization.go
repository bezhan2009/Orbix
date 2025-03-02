package src

import (
	"fmt"
	"goCmd/cmd/commands"
	"goCmd/cmd/dirInfo"
	"goCmd/src/environment"
	"goCmd/src/user"
	"goCmd/structs"
	"goCmd/system"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
)

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
					fmt.Printf("ORB %s> %s", dir, system.Green(*commandInput))
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
