package Ls

import (
	"fmt"
	"goCmd/utils"
	"io/ioutil"
	"os"
)

func PrintLS() {
	currentDir, err := os.Getwd()
	if err != nil {
		utils.AnimatedPrint(fmt.Sprint("Error getting current directory:", err))
		return
	}

	fmt.Println("\tDirectory:", currentDir)
	fmt.Println()

	files, err := ioutil.ReadDir(".")
	if err != nil {
		utils.AnimatedPrint(fmt.Sprint("Error reading directory:", err))
		return
	}

	for _, file := range files {
		mode := file.Mode().String()
		modTime := file.ModTime().Format("02.01.2006 15:04")
		name := file.Name()

		if file.IsDir() {
			name += "/"
		}

		fmt.Printf("%-20s %-20s %10d %s\n", mode, modTime, file.Size(), name)
	}
}
