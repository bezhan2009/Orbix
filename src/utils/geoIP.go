package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/internal/Network"
)

func GeoIPUtil(commandArgs []string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	if len(commandArgs) < 1 {
		fmt.Println(yellow("Usage: geoip <ip>"))
		return
	}
	Network.GeoIP(commandArgs[0])
}
