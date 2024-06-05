package main

import (
	"fmt"
	"goCmd/goCmd"
	"goCmd/utils"
)

func main() {
	if utils.IsHidden() {
		fmt.Println("You are BLOCKED!!!")
		return
	}

	goCmd.GoCmd()
}
