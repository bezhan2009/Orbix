package main

import (
	"fmt"
	"goCmd/src"
	"goCmd/utils"
)

func main() {
	if utils.IsHidden() {
		fmt.Println("You are BLOCKED!!!")
		return
	}

	src.CMD("")
}
