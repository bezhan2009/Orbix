package Orbix

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
	"goCmd/algorithms/PasswordAlgoritm"
	"goCmd/commands/commandsWithoutSignature/Clean"
	"os"
	"path/filepath"
	"strings"
)

// isPasswordDirectoryEmpty Функция, которая проверяет, есть ли файлы в директории passwords.
func isPasswordDirectoryEmpty() (bool, error) {
	files, err := os.ReadDir("passwords")
	if err != nil {
		return false, err
	}

	return len(files) == 0, nil
}

// CheckUser Функция, которая проверяет пользователя и его пароль.
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
		fmt.Printf("%s\n", green("Добро пожаловать,", usernameFromDir))
		return usernameFromDir, true
	}

	reader := bufio.NewReader(os.Stdin)
	magenta := color.New(color.FgMagenta).SprintFunc()
	fmt.Printf("%s", magenta("Введите имя пользователя: "))
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	fmt.Printf("%s", magenta("Введите пароль: "))
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	password = PasswordAlgoritm.Usage(password, true)
	hashedPassword := hashPasswordFromUser(password)
	passwordDir := filepath.Join("passwords", username)

	filePath := filepath.Join(passwordDir, hashedPassword)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		Clean.Screen()

		fmt.Printf("%s\n", magenta("Пользователь не найден или неверный пароль"))
		return usernameFromDir, false
	}

	Clean.Screen()

	fmt.Printf("%s\n", magenta("Добро пожаловать, ", username))
	return username, true
}

// hashPasswordFromUser Функция для хеширования пароля
func hashPasswordFromUser(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
