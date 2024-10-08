package src

import (
	"fmt"
	"goCmd/system"
	"time"
)

// animatedPrint custom print
func animatedPrint(text string, color string) {
	colors := system.GetColorsMap()

	for _, char := range text {
		fmt.Print(colors[color](string(char)))
		time.Sleep(1 * time.Millisecond)
	}
}
