package Orbix

import (
	"fmt"
	"time"
)

func animatedPrint(text string) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(1 * time.Millisecond)
	}
}
