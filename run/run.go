package run

import (
	"fmt"
	"goCmd/src"
	"goCmd/utils"
)

func CMD() {
	if utils.IsHidden() {
		fmt.Println("You are BLOCKED!!!")
		return
	}

	src.CMD("")
}
