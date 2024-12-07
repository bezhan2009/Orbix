package user

import (
	"bufio"
	"fmt"
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
func CheckUser(usernameFromDir string, sd *system.AppState) (string, bool) {
	currentPath, _ := os.Getwd()
	defer os.Chdir(currentPath)

	err := commands.ChangeDirectory(system.SourcePath)
	if err != nil {
		fmt.Println(system.Red("Error changing directory:", err))
		return "", false
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
		return "", false
	}

	if isEmpty {
		commands.Screen()
		fmt.Printf("%s\n", system.Green("Welcome,", usernameFromDir))
		sd.IsAdmin = true
		return usernameFromDir, true
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
			return usernameFromDir, true
		} else {
			break
		}
	}

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
					fmt.Println(system.Red("This user already exists!"))
					return "", false
				}
			}

			if !exactMatchFound {
				for _, line := range lines {
					if strings.Contains(strings.TrimSpace(line), strings.TrimSpace(username)) {
						fmt.Println(system.Red("Partial match found with: " + line))
						return "", false
					}
				}
			}
		}

		fmt.Printf("%s", system.Magenta("Enter password: "))
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Error:", err)
			return "", false
		}

		password := string(bytePassword)
		password = strings.TrimSpace(password)

		password = PasswordAlgoritm.Usage(password, true)
		hashedPassword := utils.HashPasswordFromUser(password)

		passwordsDir, err := getPasswordsDir()
		if err != nil {
			commands.Screen()
			fmt.Printf("%s\n", system.Magenta("Ошибка при получении пути директории паролей"))
			return "", false
		}

		passwordDir := filepath.Join(passwordsDir, username)
		filePath := filepath.Join(passwordDir, hashedPassword)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			commands.Screen()
			fmt.Printf("%s\n", system.Red("User not found or password is incorrect!"))
			return usernameFromDir, false
		}

		commands.Screen()
		fmt.Printf("%s\n", system.Magenta("Welcome, ", username))

		sd.IsAdmin = false
		sd.User = username
		return username, true
	}
}
