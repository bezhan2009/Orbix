package run

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/src"
	"goCmd/utils"
	"os"
)

// Init initializes CMD
func Init() {
	file, errOpen := os.Open("isRun.txt")
	if errOpen != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Println(red("you must run the program via run_main.bat or via main.sh if you are on Unix"))
		os.Exit(1)
	}
	defer file.Close()

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

	src.Orbix("", true)
}
