package run

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/utils"
	"log"
	"os"
)

// Init initializes CMD
func Init() {
	red := color.New(color.FgRed).SprintFunc()

	_, err := os.Create("running.txt")
	if err != nil {
		fmt.Println(red("Error creating running.txt: "), err)
		os.Exit(1)
	}
	file, err := os.Open("running.txt")
	if err != nil {
		log.Fatal(err)
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
