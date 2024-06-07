package ORPXI

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"goCmd/validators"
	"os"
	"path/filepath"
	"strings"
)

func NewUser() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	fmt.Print("Enter NewUser: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	isValid := validators.Password(password)

	if !isValid {
		fmt.Println("NewUser is Invalid")
		return
	}

	// Create the 'passwords' directory if it doesn't exist
	passwordDir := filepath.Join("passwords", username)
	if _, err := os.Stat(passwordDir); os.IsNotExist(err) {
		err = os.MkdirAll(passwordDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating passwords directory:", err)
			return
		}
	}

	// Hash the password
	hashedPassword := hashPassword(password)

	// Use the hash of the password as the filename
	passwordFilePath := filepath.Join(passwordDir, hashedPassword)
	err := os.WriteFile(passwordFilePath, []byte(hashedPassword), os.ModePerm)
	if err != nil {
		fmt.Println("Error writing to password file:", err)
		return
	}

	fmt.Println("Your password (file):", passwordFilePath)
	fmt.Println("Hashed password saved to file.")
	os.Exit(1)
}

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
