package utils

import (
	"fmt"
	"goCmd/Network"
)

func OpenLinkUtil(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: open_link <url>")
		return
	}

	Network.OpenBrowser(commandArgs[0])
}
