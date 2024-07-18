package run

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/src"
	"goCmd/utils"
	"log"
	"os"
)

// Init initializes CMD
func Init() {
	os.Create("running.txt")
	file, err := os.Open("running.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if utils.IsHidden() {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Println(red("You are BLOCKED!!!"))
		return
	}

	passwordsDir := "passwords"

	if _, err := os.Stat(passwordsDir); os.IsNotExist(err) {
		err := os.Mkdir(passwordsDir, 0755)
		if err != nil {
			fmt.Printf("Ошибка при создании папки %s: %v\n", passwordsDir, err)
			return
		}
	}

	src.Orbix("", true)
}
