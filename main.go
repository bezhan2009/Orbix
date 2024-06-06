package main

import (
	"fmt"
	"goCmd/ORPXI"
	"goCmd/editCMD"
	"goCmd/utils"
)

func main() {
	if utils.IsHidden() {
		fmt.Println("You are BLOCKED!!!")
		return
	}

	editCMD.StartEditing()

	ORPXI.CMD()
}
