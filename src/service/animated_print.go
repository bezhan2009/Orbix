package service

import (
	"fmt"
	"goCmd/system"
	"goCmd/utils"
	"time"
)

// AnimatedPrint custom print
func AnimatedPrint(text string, color string) {
	colors := system.GetColorsMap()

	validColors := []string{"yellow", "green", "blue", "magenta", "cyan", "red"}
	isValid := utils.IsValid(color, validColors)

	if !isValid {
		color = "green"
	}

	for _, char := range text {
		fmt.Print(colors[color](string(char)))
		time.Sleep(1 * time.Millisecond)
	}
}
