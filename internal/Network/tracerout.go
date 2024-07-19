package Network

import (
	"fmt"
	"os/exec"
)

func Traceroute(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: traceroute <hostname>")
		return
	}
	hostname := args[0]
	cmd := exec.Command("traceroute", hostname)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing traceroute:", err)
		return
	}
	fmt.Println(string(output))
}
