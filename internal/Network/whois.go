package Network

import (
	"fmt"
	"goCmd/utils"
	"os/exec"
)

func Whois(domain string) {
	cmd := exec.Command("whois", domain)
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.AnimatedPrint(fmt.Sprint("Error executing whois command:", err), "red")
		return
	}
	utils.AnimatedPrint(fmt.Sprint(string(output)), "red")
}
