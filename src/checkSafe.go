package src

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"goCmd/cmd/commands/commandsWithoutSignature/Clean"
	"goCmd/pkg/algorithms/PasswordAlgoritm"
	"goCmd/system"
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

// isPasswordDirectoryEmpty checks if there are any files in the passwords directory.
func isPasswordDirectoryEmpty() (bool, error) {
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
func CheckUser(usernameFromDir string) (string, bool) {
	isEmpty, err := isPasswordDirectoryEmpty()
	if err != nil {
		Clean.Screen()
		fmt.Printf("%s\n", red("Error checking the password directory:", err))
		return "", false
	}

	if isEmpty {
		Clean.Screen()
		fmt.Printf("%s\n", green("Welcome,", usernameFromDir))
		return usernameFromDir, true
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(magenta("enable secure[Y/N]: "))
		var enable string
		enable, _ = reader.ReadString('\n')
		enable = strings.TrimSpace(enable)

		if enable == "" {
			fmt.Println(red("You entered a blank value!"))
			continue
		}

		if strings.ToLower(strings.TrimSpace(enable)) != "y" {
			return usernameFromDir, true
		} else {
			system.IsAdmin = false
			break
		}
	}

	for {
		fmt.Printf("%s", magenta("Enter username: "))
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		if username == "" {
			fmt.Println(red("You entered a blank value!"))
			continue
		}

		runningPath := AbsDirRun
		runningPath += "\\running.txt"
		sourceRunning, errReading := os.ReadFile(runningPath)

		var dataRunning string

		if errReading == nil {
			dataRunning = string(sourceRunning)
			lines := strings.Split(dataRunning, "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) == strings.TrimSpace(username) {
					fmt.Println(red("This user already exists!"))
					return "", false
				}
			}
		}

		fmt.Printf("%s", magenta("Enter password: "))
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Error:", err)
			return "", false
		}

		password := string(bytePassword)
		password = strings.TrimSpace(password)

		password = PasswordAlgoritm.Usage(password, true)
		hashedPassword := hashPasswordFromUser(password)

		passwordsDir, err := getPasswordsDir()
		if err != nil {
			Clean.Screen()
			fmt.Printf("%s\n", magenta("Ошибка при получении пути директории паролей"))
			return "", false
		}

		passwordDir := filepath.Join(passwordsDir, username)
		filePath := filepath.Join(passwordDir, hashedPassword)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			Clean.Screen()
			fmt.Printf("%s\n", red("User not found or password is incorrect!"))
			return usernameFromDir, false
		}

		Clean.Screen()
		fmt.Printf("%s\n", magenta("Welcome, ", username))
		system.User = username
		return username, true
	}
}

// hashPasswordFromUser hashes the user's password.
func hashPasswordFromUser(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
