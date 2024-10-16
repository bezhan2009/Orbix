package commands

import (
	"fmt"
	"goCmd/pkg/algorithms/PasswordAlgoritm"
	"goCmd/system"
	"goCmd/utils"
	"os"
	"path/filepath"
	"strings"
)

func SetUserFromENV(systemPath string) {
	colors := system.GetColorsMap()

	fmt.Println(systemPath)
	err := os.Chdir(systemPath)
	if err != nil {
		fmt.Println(colors["red"]("Error when changing the path:", err))
	}

	username := string(strings.TrimSpace(os.Getenv("DEFAULT_USER")))

	password := string(strings.TrimSpace(os.Getenv("USERS_DEFAULT_PASSWORD")))
	password = strings.TrimSpace(password)

	password = PasswordAlgoritm.Usage(password, true)
	hashedPassword := utils.HashPasswordFromUser(password)

	passwordDir := filepath.Join("passwords", username)
	if _, err = os.Stat(passwordDir); os.IsNotExist(err) {
		err = os.MkdirAll(passwordDir, os.ModePerm)
		if err != nil {
			fmt.Println(colors["red"]("Error creating passwords directory:", err))
			return
		}
	}

	// Use the hash of the password as the filename
	passwordFilePath := filepath.Join(passwordDir, hashedPassword)
	err = os.WriteFile(passwordFilePath, []byte(hashedPassword), os.ModePerm)
	if err != nil {
		fmt.Println(colors["red"]("Error writing to password file:", err))
		return
	}

	fmt.Println(colors["cyan"]("Username:", username))
	fmt.Println(colors["cyan"]("Password:", password))
	fmt.Println()
	fmt.Println(colors["green"]("Your password (file):", passwordFilePath))
	fmt.Println(colors["green"]("Hashed password saved to file."))
}
