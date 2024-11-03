package utils

import (
	"fmt"
	"goCmd/internal/OS"
	"goCmd/system"
)

func KillProcessUtil(commandArgs []string) {
	colors := system.GetColorsMap()

	if len(commandArgs) < 1 {
		fmt.Println(colors["yellow"]("kill <PID>"))
		return
	}

	PID := commandArgs[0]

	err := OS.KillProcess(PID)
	if err != nil {
		fmt.Println(colors["red"]("Error killing process: " + err.Error()))
		return
	}

	fmt.Println(colors["green"]("Process " + PID + " killed."))
}
