package Orbix

import (
	"fmt"
)

func ExecuteShablonUtil(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: shablon <template_name>")
		return
	}
	if err := Start(commandArgs[0]); err != nil {
		fmt.Println(err)
	}
}
