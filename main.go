package main

import (
	"fmt"
	"goCmd/Orbix"
	"goCmd/utils"
)

func main() {
	if utils.IsHidden() {
		fmt.Println("You are BLOCKED!!!")
		return
	}

	Orbix.CMD("")
}
