package Ls

import (
	"fmt"
	"goCmd/utils"
	"io/ioutil"
)

func PrintLS() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		utils.AnimatedPrint(fmt.Sprint("Error reading directory:", err))
		return
	}
	for _, file := range files {
		if file.IsDir() {
			utils.AnimatedPrint(fmt.Sprintf("%s/\t", file.Name()))
		} else {
			utils.AnimatedPrint(fmt.Sprint(file.Name(), "\t"))
		}
	}
}
