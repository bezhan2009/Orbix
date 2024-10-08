package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/internal/Network"
	"strconv"
)

func ScanPortUtil(commandArgs []string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	if len(commandArgs) < 2 {
		fmt.Println(yellow("Usage: scanport <host> <ports>"))
		return
	}

	var ports []int

	for _, p := range commandArgs[1:] {
		port, err := strconv.Atoi(p)
		if err != nil {
			printErr := fmt.Sprintf("Invalid port: %s\n", p)
			fmt.Printf(red(printErr))
			return
		}
		ports = append(ports, port)
	}
	Network.ScanPort(commandArgs[0], ports)
}
