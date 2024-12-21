package utils

import (
	"fmt"
	"goCmd/system"
	"time"
)

func PrintAnim(text string) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(1 * time.Millisecond)
	}
}

func AnimatedPrint(text string, color string) {
	colors := system.GetColorsMap()
	for _, char := range text {
		fmt.Print(colors[color](string(char)))
		time.Sleep(1 * time.Millisecond)
	}
}

func AnimatedPrintLong(text string, color string) {
	colors := system.GetColorsMap()
	for _, char := range text {
		fmt.Print(colors[color](string(char)))
		time.Sleep(1 * time.Second)
	}
}
