package utils

import (
	"fmt"
	"goCmd/internal/Network"
)

func DnsLookupUtil(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: dnslookup <domain>")
		return
	}
	Network.DNSLookup(commandArgs[0])
}
