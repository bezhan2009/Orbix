package Ls

import (
	"fmt"
	"io/ioutil"
)

func PrintLS() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			fmt.Printf("%s/\t", file.Name())
		} else {
			fmt.Print(file.Name(), "\t")
		}
	}
}
