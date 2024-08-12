package run

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/src"
	"goCmd/system"
	"goCmd/utils"
	"log"
	"os"
)

// Init initializes CMD
func Init() {
	red := color.New(color.FgRed).SprintFunc()

	os.Create("running.txt")
	file, err := os.Open("running.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	system.Path = utils.Getwd()

	if utils.IsHidden() {
		fmt.Println(red("You are BLOCKED!!!"))
		return
	}

	passwordsDir := "passwords"

	if _, err := os.Stat(passwordsDir); os.IsNotExist(err) {
		err := os.Mkdir(passwordsDir, 0755)
		if err != nil {
			printErr := fmt.Sprintf("Error creating folder %s: %v\n", passwordsDir, err)
			fmt.Println(red(printErr))
			return
		}
	}

	src.Orbix("", true)
}
