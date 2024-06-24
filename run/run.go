package run

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/src"
	"goCmd/utils"
	"os"
	"path/filepath"
)

var Absdir, _ = filepath.Abs("")

// Init initializes CMD
func Init() {
	runFilePath := filepath.Join(Absdir, "isRun.txt")

	file, err := os.Open(runFilePath)
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Println(red("Запустите программу через run_main.bat либо если у вас Unix(Linux, MacOS) то запустите через main.sh"))
		os.Exit(1)
	}
	defer file.Close()

	os.Remove(runFilePath)

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
