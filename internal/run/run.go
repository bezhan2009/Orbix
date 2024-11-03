package run

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"goCmd/utils"
	"log"
	"os"
)

// Init initializes CMD
func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, creating a new one:", err)
		// Открываем или создаем .env файл
		file, err := os.Create(".env")
		if err != nil {
			log.Fatal("Error creating .env file:", err)
		}
		defer file.Close()

		// Записываем содержимое в .env файл
		content := `BETA: N
DEFAULT_USER: orbix
USERS_DEFAULT_PASSWORD: 12345678
USE_NEW_PROMPT: Y
USERS_LOCATION: USA
PROMPT: _>`
		_, err = file.WriteString(content)
		if err != nil {
			log.Fatal("Error writing to .env file:", err)
		}

		fmt.Println(".env file created successfully")
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}

	red := color.New(color.FgRed).SprintFunc()

	file, err := os.Open("running.txt")
	if err != nil {
		file, err = os.Create("running.txt")
		if err != nil {
			fmt.Println(red("Error creating running.txt: "), err)
			os.Exit(1)
		}
	}
	defer func() {
		err = file.Close()
		if err != nil {
			return
		}
	}()

	if utils.IsHidden() {
		fmt.Println(red("You are BLOCKED!!!"))
		os.Exit(1)
	}

	passwordsDir := "passwords"

	if _, err = os.Stat(passwordsDir); os.IsNotExist(err) {
		err = os.Mkdir(passwordsDir, 0755)
		if err != nil {
			printErr := fmt.Sprintf("Error creating folder %s: %v\n", passwordsDir, err)
			fmt.Println(red(printErr))
			os.Exit(1)
		}
	}
}
