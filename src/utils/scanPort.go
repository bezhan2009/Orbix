package utils

import (
	"fmt"
	"goCmd/internal/Network"
	"strconv"
)

func ScanPortUtil(commandArgs []string) {
	if len(commandArgs) < 2 {
		fmt.Println("Usage: scanport <host> <ports>")
		return
	}
	var ports []int
	for _, p := range commandArgs[1:] {
		port, err := strconv.Atoi(p)
		if err != nil {
			fmt.Printf("Invalid port: %s\n", p)
			return
		}
		ports = append(ports, port)
	}
	Network.ScanPort(commandArgs[0], ports)
}
