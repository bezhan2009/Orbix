package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/internal/Network"
)

func IPInfoUtil(commandArgs []string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	if len(commandArgs) < 1 {
		fmt.Println(yellow("Usage: ipinfo <ip>"))
		return
	}
	Network.IPInfo(commandArgs[0])
}
