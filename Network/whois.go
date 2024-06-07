package Network

import (
	"fmt"
	"os/exec"
)

func Whois(domain string) {
	cmd := exec.Command("whois", domain)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing whois command:", err)
		return
	}
	fmt.Println(string(output))
}
