package Network

import (
	"fmt"
	"os/exec"
)

func Ping(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: pingview <hostname>")
		return
	}
	hostname := args[0]
	cmd := exec.Command("ping", hostname)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing ping:", err)
		return
	}
	fmt.Println(string(output))
}
