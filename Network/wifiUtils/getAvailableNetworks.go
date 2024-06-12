package wifiUtils

import (
	"fmt"
	"os/exec"
	"strings"
)

func getAvailableNetworks() []string {
	cmd := exec.Command("netsh", "wlan", "show", "networks")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Ошибка при сканировании Wi-Fi:", err)
		return nil
	}

	lines := strings.Split(string(output), "\n")
	var networks []string
	for _, line := range lines {
		if strings.Contains(line, "SSID") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				networks = append(networks, strings.TrimSpace(parts[1]))
			}
		}
	}
	return networks
}
