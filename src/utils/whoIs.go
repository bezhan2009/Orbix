package utils

import (
	"fmt"
	"goCmd/internal/Network"
)

func WhoisUtil(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: whois <domain>")
		return
	}
	Network.Whois(commandArgs[0])
}
