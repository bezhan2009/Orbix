package mycmd

import (
	"fmt"
	"os"
)

func Start(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: mycmd <script.mycmd>")
		return
	}

	fileName := os.Args[0]
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading script file: %v\n", err)
		return
	}

	Interpreter(string(content))
}
