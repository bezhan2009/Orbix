package src

import (
	"fmt"
	"time"
)

// animatedPrint custom print
func animatedPrint(text string) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(1 * time.Millisecond)
	}
}
