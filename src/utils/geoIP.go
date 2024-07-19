package utils

import (
	"fmt"
	"goCmd/internal/Network"
)

func GeoIPUtil(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: geoip <ip>")
		return
	}
	Network.GeoIP(commandArgs[0])
}
