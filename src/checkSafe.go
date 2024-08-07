package src

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithoutSignature/Clean"
	"goCmd/pkg/algorithms/PasswordAlgoritm"
	"os"
	"path/filepath"
	"strings"
)

var AbsdirRun, _ = filepath.Abs("")

// getPasswordsDir returns the absolute path of the passwords directory.
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
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("%s\n", red("Ошибка при проверке директории с паролями:", err))
		return "", false
	}

	if isEmpty {
		Clean.Screen()
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s\n", green("Welcome,", usernameFromDir))
		return usernameFromDir, true
	}

	reader := bufio.NewReader(os.Stdin)
	magenta := color.New(color.FgMagenta).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	for {
		fmt.Print(magenta("enable secure[Y/N]: "))
		var enable string
		enable, _ = reader.ReadString('\n')
		enable = strings.TrimSpace(enable)

		if enable == "" {
			fmt.Println(red("You entered a blank value!"))
			continue
		}

		if strings.ToLower(enable) != "y" {
			return usernameFromDir, true
		} else {
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
		runningPath := AbsdirRun
		runningPath += "\\running.txt"
		sourceRunning, errReading := os.ReadFile(runningPath)
		var dataRunning string

		if errReading == nil {
			dataRunning = string(sourceRunning)
			lines := strings.Split(dataRunning, "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) == username {
					fmt.Println(red("This user already exists!"))
					return "", false
				}
			}
		}

		fmt.Printf("%s", magenta("Enter password: "))
		password, _ := reader.ReadString('\n')
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
		return username, true
	}
}

// hashPasswordFromUser hashes the user's password.
func hashPasswordFromUser(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
