package ORPXI

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Функция, которая проверяет, есть ли файлы в директории passwords.
func isPasswordDirectoryEmpty() (bool, error) {
	files, err := os.ReadDir("passwords")
	if err != nil {
		return false, err
	}

	return len(files) == 0, nil
}

// Функция, которая проверяет пользователя и его пароль.
func CheckUser(usernameFromDir string) bool {
	isEmpty, err := isPasswordDirectoryEmpty()
	if err != nil {
		fmt.Println("Ошибка при проверке директории с паролями:", err)
		return false
	}

	if isEmpty {
		fmt.Println("Добро пожаловать,", usernameFromDir)
		return true
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите имя пользователя: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	fmt.Print("Введите пароль: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Хешируем пароль
	hashedPassword := hashPasswordFromUser(password)
	passwordDir := filepath.Join("passwords", username)

	// Ищем файл с именем, совпадающим с хешем пароля
	filePath := filepath.Join(passwordDir, hashedPassword)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Пользователь не найден или неверный пароль")
		return false
	}

	// Если файл существует, пользователь найден
	fmt.Println("Добро пожаловать,", username)
	return true
}

// Функция для хеширования пароля
func hashPasswordFromUser(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
