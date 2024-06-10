package WiFi

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func Scan() {
	cmd := exec.Command("netsh", "wlan", "show", "networks")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Ошибка при сканировании Wi-Fi:", err)
		return
	}
	fmt.Println(string(output))
}

// Connect пытается подключиться к Wi-Fi сети с заданным SSID и паролем.
func Connect(ssid, password string) bool {
	config := fmt.Sprintf(`<?xml version="1.0"?>
<WLANProfile xmlns="http://www.microsoft.com/networking/WLAN/profile/v1">
    <name>%s</name>
    <SSIDConfig>
        <SSID>
            <name>%s</name>
        </SSID>
    </SSIDConfig>
    <connectionType>ESS</connectionType>
    <connectionMode>manual</connectionMode>
    <MSM>
        <security>
            <authEncryption>
                <authentication>WPA2PSK</authentication>
                <encryption>AES</encryption>
                <useOneX>false</useOneX>
            </authEncryption>
            <sharedKey>
                <keyType>passPhrase</keyType>
                <protected>false</protected>
                <keyMaterial>%s</keyMaterial>
            </sharedKey>
        </security>
    </MSM>
</WLANProfile>`, ssid, ssid, password)

	filename := ssid + ".xml"
	err := ioutil.WriteFile(filename, []byte(config), 0644)
	if err != nil {
		fmt.Println("Ошибка при создании файла конфигурации Wi-Fi:", err)
		return false
	}
	defer os.Remove(filename)

	cmd := exec.Command("netsh", "wlan", "add", "profile", "filename="+filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Ошибка при добавлении профиля Wi-Fi: %v\n", err)
		return false
	}

	cmd = exec.Command("netsh", "wlan", "connect", "name="+ssid)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Ошибка при подключении к Wi-Fi: %v\n", err)
		return false
	}

	// Проверка подключения
	cmd = exec.Command("netsh", "wlan", "show", "interfaces")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Ошибка при проверке состояния подключения: %v\n", err)
		return false
	}
	outputStr := string(output)
	if strings.Contains(outputStr, ssid) && strings.Contains(outputStr, "Connected") {
		fmt.Println("Подключение успешно")
		return true
	}
	return false
}
