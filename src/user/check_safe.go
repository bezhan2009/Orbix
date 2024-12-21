package user

import (
	"bufio"
	"fmt"
	_chan "goCmd/chan"
	"goCmd/cmd/commands"
	"goCmd/pkg/algorithms/PasswordAlgoritm"
	"goCmd/system"
	"goCmd/utils"
	"golang.org/x/term"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var AbsDirRun, _ = filepath.Abs("")

// getPasswordsDir returns the absolute path of the passwords' directory.
func getPasswordsDir() (string, error) {
	return filepath.Abs("passwords")
}

// IsPasswordDirectoryEmpty checks if there are any files in the passwords directory.
func IsPasswordDirectoryEmpty() (bool, error) {
	passwordsDir, err := getPasswordsDir()
	if err != nil {
		return false, err
	}

	files, err := os.ReadDir(passwordsDir)
	if err != nil {
		return false, err
	}

	return len(files) == 0, nil
}

// CheckUser checks the username and password.
func CheckUser(usernameFromDir string, sd *system.AppState) (string, bool, error) {
	if !system.Unauthorized {
		_chan.UserStatusAuth = true
	}

	system.Unauthorized = true
	currentPath, _ := os.Getwd()
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			fmt.Println(system.Red("Error changing directory:", err))
		}
	}(currentPath)

	err := commands.ChangeDirectory(system.SourcePath)
	if err != nil {
		fmt.Println(system.Red("Error changing directory:", err))
		return "", false, err
	}

	if sd == nil {
		fmt.Println(system.Red("Fatal: Session is nil!!!"))
		os.Exit(1)
	}

	isEmpty, err := IsPasswordDirectoryEmpty()
	if err != nil {
		commands.Screen()
		fmt.Printf("%s\n", system.Red("Error checking the password directory:", err))
		sd.IsAdmin = true
		return "", false, err
	}

	if isEmpty {
		commands.Screen()
		fmt.Printf("%s\n", system.Green("Welcome,", usernameFromDir))
		_chan.EnableSecure = true

		sd.IsAdmin = true
		return usernameFromDir, true, err
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(system.Magenta("enable secure[Y/N]: "))

		var enable string
		enable, _ = reader.ReadString('\n')
		enable = strings.TrimSpace(enable)

		if enable == "" {
			fmt.Println(system.Red("You entered a blank value!"))
			continue
		}

		if strings.ToLower(strings.TrimSpace(enable)) != "y" {
			sd.IsAdmin = true
			_chan.UseOldPrompt = true
			_chan.EnableSecure = false
			return usernameFromDir, true, nil
		} else {
			break
		}
	}

	var attempts uint = 0

	for {
		fmt.Printf("%s", system.Magenta("Enter username: "))
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		if username == "" {
			fmt.Println(system.Red("You entered a blank value!"))
			continue
		}

		runningPath := AbsDirRun
		runningPath += fmt.Sprintf("\\%s", system.OrbixRunningUsersFileName)
		sourceRunning, errReading := os.ReadFile(runningPath)

		var dataRunning string

		if errReading == nil {
			dataRunning = string(sourceRunning)
			lines := strings.Split(dataRunning, "\n")
			exactMatchFound := false

			for _, line := range lines {
				if strings.TrimSpace(line) == strings.TrimSpace(username) {
					system.Unauthorized = true
					fmt.Println(system.Red("This user already exists!"))
					return "", false, system.UserAlreadyExists
				}
			}

			if !exactMatchFound {
				for _, line := range lines {
					if strings.Contains(strings.TrimSpace(line), strings.TrimSpace(username)) {
						fmt.Println(system.Red("Partial match found with: " + line))
						return "", false, system.ExactMatchNotFound
					}
				}
			}
		}

		fmt.Printf("%s", system.Magenta("Enter password: "))
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Error:", err)
			return "", false, err
		}

		password := string(bytePassword)
		password = strings.TrimSpace(password)

		password = PasswordAlgoritm.Usage(password, true)
		hashedPassword := utils.HashPasswordFromUser(password)

		passwordsDir, err := getPasswordsDir()
		if err != nil {
			commands.Screen()
			fmt.Printf("%s\n", system.Red("Error getting password directory:", err))
			return "", false, err
		}

		passwordDir := filepath.Join(passwordsDir, username)
		filePath := filepath.Join(passwordDir, hashedPassword)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			commands.Screen()
			system.Unauthorized = true
			fmt.Printf("%s\n", system.Red("User not found or password is incorrect!"))
			attempts += 1
			if !(attempts > system.MaxUserAuthAttempts) {
				fmt.Println(system.Red("Try one more time!"))
				continue
			}

			return usernameFromDir, false, err
		}

		commands.Screen()
		fmt.Printf("%s\n", system.Magenta("Welcome, ", username))

		_chan.EnableSecure = true
		sd.IsAdmin = false
		sd.User = username
		return username, true, nil
	}
}
