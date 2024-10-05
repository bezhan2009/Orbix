package src

import (
	"fmt"
	"goCmd/cmd/commands"
	"goCmd/cmd/dirInfo"
	ExCommUtils "goCmd/src/utils"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

var (
	Absdir, _             = filepath.Abs("")
	DirUser, _            = filepath.Abs("")
	UsernameFromDir       = dirInfo.CmdUser(DirUser)
	GlobalSession         = system.Session{}
	PreviousSessionPath   = ""
	PreviousSessionPrefix = ""
)

func Orbix(commandInput string, echo bool, rebooted structs.RebootedData, SD *system.AppState) {
	system.OrbixWorking = true
	RestartAfterInit := false
	if strings.TrimSpace(commandInput) == "restart" {
		RestartAfterInit = true
	}

	if SD == nil {
		fmt.Println(red("Fatal: Session is nil!!!"))
		os.Exit(1)
	}

	if err := commands.ChangeDirectory(Absdir); err != nil {
		fmt.Println(red(err))
	}

	sessionData := SD

	// Initialize colors
	InitColors()

	if !echo && commandInput == "" {
		fmt.Println(red("You cannot enable echo with an empty Input command!"))
		return
	}

	if echo && rebooted.Username == "" && commandInput == "" {
		utils.SystemInformation()
	}

	isWorking := true
	isPermission := true
	if commandInput != "" {
		isPermission = false
	}

	// Check if password directory is empty once and handle errors here
	isEmpty, err := isPasswordDirectoryEmpty()
	if err != nil {
		animatedPrint("Error checking password directory: " + err.Error() + "\n")
		return
	}

	username := ""

	if strings.TrimSpace(rebooted.Username) != "" {
		username = strings.TrimSpace(rebooted.Username)
	} else if !isEmpty && commandInput == "" {
		dir, _ := os.Getwd()
		user := dirInfo.CmdUser(dir)

		nameUser, isSuccess := CheckUser(user, sessionData)
		if !isSuccess {
			return
		}

		username = nameUser
		initializeRunningFile(username) // New helper function

		if user == username {
			sessionData.IsAdmin = true
			sessionData.User = user
		} else {
			sessionData.IsAdmin = false
			sessionData.User = username
		}
	}

	location := os.Getenv("CITY")
	if location == "" {
		location = "Unknown City"
	}

	// Signal handling setup (outside the loop)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var signalReceived bool

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			<-signalChan
			signalReceived = true
			fmt.Print(red("^C"))
			if !sessionData.IsAdmin {
				dir, _ := os.Getwd()

				dirC := dirInfo.CmdDir(dir)
				user := sessionData.User
				if user == "" {
					user = dirInfo.CmdUser(dir)
				}

				var printUserDir string

				if username != "" {
					user = username
					printUserDir = user
				}
				fmt.Println()
				gitBranch, _ := GetCurrentGitBranch()
				printPromptInfo(location, printUserDir, dirC, green, cyan, yellow, magenta, &system.Session{Path: dir, GitBranch: gitBranch}, commandInput)
			} else {
				dir, _ := os.Getwd()
				fmt.Println()
				fmt.Printf("\nORB %s>%s", dir, green(commandInput))
			}
		}
	}()

	originalStdout, originalStderr := os.Stdout, os.Stderr
	devNull, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0666)
	defer func() {
		err = devNull.Close()
		if err != nil {
			return
		}
	}()

	var prompt string
	var prefix string

	prefix = sessionData.NewSessionData(sessionData.Path, sessionData.User, sessionData.GitBranch, sessionData.IsAdmin)

	if rebooted.Prefix != "" {
		prefix = rebooted.Prefix
	} else {
		prefix = sessionData.NewSessionData(sessionData.Path, sessionData.User, sessionData.GitBranch, sessionData.IsAdmin)
	}

	session, exists := sessionData.GetSession(prefix)
	if !exists {
		fmt.Println(red("Session does not exist!"))
		return
	}

	if session == nil {
		fmt.Println(red("Session is nil!"))
		return
	}

	// Initialize Global Vars
	Init(session)

	session.PreviousPath = PreviousSessionPath
	fmt.Println(green(session.PreviousPath))
	if PreviousSessionPrefix != "" {
		session, _ = sessionData.GetSession(PreviousSessionPrefix)
	}

	GlobalSession = *session

	dir, _ := os.Getwd()
	system.Path = dir

	for isWorking {
		system.OrbixWorking = true

		if len(session.CommandHistory) < 10 {
			Init(session)
		}

		// Check if signal was received and reset flag after handling it
		if signalReceived {
			signalReceived = false
			continue // Continue the loop after signal
		}

		// Redirect output based on the echo setting
		if echo {
			os.Stdout, os.Stderr = originalStdout, originalStderr
		} else {
			os.Stdout, os.Stderr = devNull, devNull
		}

		// Directory and user context setup (execute only when necessary)
		dir, _ = os.Getwd()
		printUserDir := UsernameFromDir // Use cached username for printing

		if RestartAfterInit {
			SD.User = username
			SD.IsAdmin = sessionData.IsAdmin
			rebooted.Prefix = prefix
			Orbix(commandInput, echo, rebooted, SD)
			return
		}

		if echo && session.IsAdmin {
			if prompt == "" {
				fmt.Printf("ORB %s>%s", dir, green(commandInput))
			} else {
				fmt.Print(green(prompt))
			}
		}

		if !session.IsAdmin {
			dirC := dirInfo.CmdDir(dir)
			user := session.User
			if user == "" {
				user = dirInfo.CmdUser(dir)
			}

			if username != "" {
				user = username
				printUserDir = user
			}

			// Single user check outside repeated prompt formatting
			if !checkUserInRunningFile(username) {
				fmt.Println(red("User not authorized."))
				isWorking = false
				isPermission = false
				continue
			}

			if echo {
				if prompt == "" {
					printPromptInfo(location, printUserDir, dirC, green, cyan, yellow, magenta, session, commandInput) // New helper function
				} else {
					fmt.Print(green(prompt))
				}
			}
		}

		// Command processing
		commandLine, command, commandArgs, commandLower := readCommandLine(commandInput) // Refactored input handling
		if commandLine == "" {
			continue
		}

		if commandInt, err := strconv.Atoi(command); err == nil && len(commandArgs) == 0 {
			fmt.Println(magenta(commandInt))
			continue
		}

		session.CommandHistory = append(session.CommandHistory, commandLine)

		session.Path = dir
		GlobalSession.CommandHistory = session.CommandHistory

		isValid := utils.ValidCommand(commandLower, Commands)

		if !isValid {
			fullCommand := append([]string{command}, commandArgs...)
			err = utils.ExternalCommand(fullCommand)
			if err != nil {
				fullPath := filepath.Join(dir, command)
				fullCommand[0] = fullPath
				err = utils.ExternalCommand(fullCommand)
				if err != nil {
					isValid = utils.ValidCommand(commandLower, AdditionalCommands)
					if !isValid {
						HandleUnknownCommandUtil(commandLower, commandLine, Commands)
						continue
					}
				}
			}
			continue
		}

		if strings.TrimSpace(commandLower) == "prompt" {
			handlePromptCommand(commandArgs, &prompt)
			continue
		}

		// Process command
		if err = processCommand(commandLower, commandArgs, session); err != nil {
			fmt.Println(red(err.Error()))
		}

		execCommand := structs.ExecuteCommandFuncParams{
			Command:       command,
			CommandLower:  commandLower,
			CommandArgs:   commandArgs,
			Dir:           dir,
			IsWorking:     &isWorking,
			IsPermission:  isPermission,
			Username:      username,
			SD:            sessionData,
			SessionPrefix: prefix,
		}

		if strings.TrimSpace(commandLower) == "orbix" {
			PreviousSessionPrefix = prefix
		}

		if strings.TrimSpace(commandLower) == "neofetch" {
			ExCommUtils.NeofetchUtil(execCommand, session, Commands)
			continue
		}

		ExecuteCommand(execCommand)
	}

	// Restore original outputs
	os.Stdout, os.Stderr = originalStdout, originalStderr
	PreviousSessionPath = session.Path
	session, _ = sessionData.GetSession(PreviousSessionPrefix)

	if err = commands.ChangeDirectory(session.Path); err != nil {
		fmt.Println(red(err))
	}

	system.UserName = session.User

	sessionData.DeleteSession(prefix)

	system.OrbixWorking = false
}
