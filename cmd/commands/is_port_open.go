package commands

import (
	"fmt"
	"goCmd/internal/OS"
	"goCmd/system"
	"strconv"
	"strings"
)

func IsPortOpen(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println(system.Yellow("Usage: chport <port>"))
		fmt.Println(system.Yellow("Flags: --without-host\n\tDesc: In default we using localhost for checking port availability. \n\tIf that flag is turn on than we checking port availability without host"))
		return
	}

	host := "localhost"

	if len(commandArgs) > 1 {
		if strings.TrimSpace(commandArgs[1]) == "--without-host" {
			host = ""
		}
	}

	portInt, err := strconv.Atoi(commandArgs[0])
	if err != nil {
		fmt.Println(system.Red("Invalid port number"))
		return
	}

	if OS.IsPortOpen(host, portInt) {
		fmt.Println(system.GreenBold("Port " + strconv.Itoa(portInt) + " is open"))
	} else {
		fmt.Println(system.RedBold("Port " + strconv.Itoa(portInt) + " is not open"))
	}
}
