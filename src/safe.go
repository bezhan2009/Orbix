package src

import (
	"bufio"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"goCmd/pkg/algorithms/PasswordAlgoritm"
	"goCmd/utils"
	"goCmd/validators"
	"os"
	"path/filepath"
	"strings"
)

// PrintNewUser displays the message NewUser!!! ASCII-style art in green
func PrintNewUser() {
	myFigure := figure.NewFigure("New User!!!", "", true)
	greenText := color.New(color.FgGreen).SprintFunc()
	fmt.Println(greenText(myFigure.String()))
}

func NewUser(systemPath string) {
	fmt.Println(systemPath)
	err := os.Chdir(systemPath)
	if err != nil {
		fmt.Println(red("Error when changing the path:", err))
	}

	reader := bufio.NewReader(os.Stdin)
	PrintNewUser()

	fmt.Printf("%s", magenta("Enter username: "))
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Printf("%s", magenta("Enter password: "))
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	isValid := validators.Password(password)

	if !isValid {
		fmt.Println(red("NewUser is Invalid"))
		return
	}

	passwordDir := filepath.Join("passwords", username)
	if _, err = os.Stat(passwordDir); os.IsNotExist(err) {
		err = os.MkdirAll(passwordDir, os.ModePerm)
		if err != nil {
			fmt.Println(red("Error creating passwords directory:", err))
			return
		}
	}

	// Encrypt and hash the password
	encryptedPassword := PasswordAlgoritm.Usage(password, true)
	hashedPassword := utils.HashPasswordFromUser(encryptedPassword)

	// Use the hash of the password as the filename
	passwordFilePath := filepath.Join(passwordDir, hashedPassword)
	err = os.WriteFile(passwordFilePath, []byte(hashedPassword), os.ModePerm)
	if err != nil {
		fmt.Println(red("Error writing to password file:", err))
		return
	}

	fmt.Println(green("Your password (file):", passwordFilePath))
	fmt.Println(green("Hashed password saved to file."))
}
