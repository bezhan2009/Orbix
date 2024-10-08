package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/internal/Network"
)

func DnsLookupUtil(commandArgs []string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	if len(commandArgs) < 1 {
		fmt.Println(yellow("Usage: dnslookup <domain>"))
		return
	}
	Network.DNSLookup(commandArgs[0])
}
