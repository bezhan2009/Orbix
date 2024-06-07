package utils

import (
	"fmt"
	"time"
)

func AnimatedPrint(text string) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(1 * time.Millisecond) // Задержка 50 миллисекунд между символами
	}
}
