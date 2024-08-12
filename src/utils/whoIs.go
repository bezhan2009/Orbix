package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/internal/Network"
)

func WhoisUtil(commandArgs []string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	if len(commandArgs) < 1 {
		fmt.Println(yellow("Usage: whois <domain>"))
		return
	}
	Network.Whois(commandArgs[0])
}
