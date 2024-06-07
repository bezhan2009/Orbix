package utils

import (
	"fmt"
	"time"
)

func AnimatedPrint(text string) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(10 * time.Millisecond) // Задержка 50 миллисекунд между символами
	}
}
