package run

import (
	"fmt"
	"goCmd/src"
	"goCmd/utils"
	"os"
	"path/filepath"
)

var Absdir, _ = filepath.Abs("")

// Init initializes CMD
func Init() {
	runFilePath := filepath.Join(Absdir, "isRun.txt")

	isRunSource, err := os.ReadFile(runFilePath)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", runFilePath, err)
		os.Exit(1)
	}

	fmt.Println(string(isRunSource))

	if string(isRunSource) == "false" {
		err := os.WriteFile(runFilePath, []byte("true"), 0777)
		if err != nil {
			fmt.Printf("Error writing to %s: %v\n", runFilePath, err)
			os.Exit(1)
		}
	} else if string(isRunSource) == "true" {
		fmt.Println("Program is already running.")
		os.Exit(1)
	}

	if utils.IsHidden() {
		fmt.Println("You are BLOCKED!!!")
		return
	}

	passwordsDir := "passwords"

	if _, err := os.Stat(passwordsDir); os.IsNotExist(err) {
		err := os.Mkdir(passwordsDir, 0755)
		if err != nil {
			fmt.Printf("Ошибка при создании папки %s: %v\n", passwordsDir, err)
			return
		}
		fmt.Printf("Папка %s успешно создана.\n", passwordsDir)
	}

	src.Orbix("")
}
