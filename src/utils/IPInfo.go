package utils

import (
	"fmt"
	"goCmd/internal/Network"
)

func IPInfoUtil(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: ipinfo <ip>")
		return
	}
	Network.IPInfo(commandArgs[0])
}
