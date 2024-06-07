package ORPXI

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
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
func CheckUser(username string) bool {
	isIt, err := isPasswordDirectoryEmpty()

	if err != nil {
		fmt.Println(err)
	}

	if !isIt {
		var password string
		fmt.Print("Введите пароль:")
		fmt.Scan(&password)
		// Хешируем пароль
		hash := sha256.New()
		hash.Write([]byte(password))
		passwordHash := fmt.Sprintf("%x", hash.Sum(nil))

		// Ищем файл с именем, совпадающим с хешем пароля
		filePath := filepath.Join("passwords", passwordHash)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Println("Пользователь не найден или неверный пароль")
			return false
		}

		// Если файл существует, пользователь найден
		fmt.Println("Добро пожаловать,", username)
		return true
	} else {
		fmt.Println("Добро пожаловать,", username)
		return true
	}

}

// Функция для создания файла с хешем пароля для пользователя
func CreatePasswordFile(username, password string) error {
	// Хешируем пароль
	hash := sha256.New()
	hash.Write([]byte(password))
	passwordHash := fmt.Sprintf("%x", hash.Sum(nil))

	// Создаем файл с именем, совпадающим с хешем пароля пользователя
	file, err := os.Create(filepath.Join("passwords", passwordHash))
	if err != nil {
		return err
	}
	defer file.Close()

	// Записываем имя пользователя и его хеш в файл
	_, err = file.WriteString(fmt.Sprintf("%s:%s", username, passwordHash))
	if err != nil {
		return err
	}

	return nil
}
