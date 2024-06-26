package run

import (
	"fmt"
	"goCmd/src"
	"goCmd/utils"
	"os"
)

// Init initializes CMD
func Init() {
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
