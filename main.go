package main

import (
	"fmt"
	"goCmd/editCMD"
	"goCmd/goCmd"
	"goCmd/utils"
)

func main() {
	if utils.IsHidden() {
		fmt.Println("You are BLOCKED!!!")
		return
	}

	editCMD.StartEditing()

	goCmd.GoCmd()
}
