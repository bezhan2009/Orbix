package WiFi

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

func Connect(ssid, password string) {
	// Создаем файл конфигурации Wi-Fi профиля
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
		return
	}
	defer os.Remove(filename)

	// Добавляем профиль Wi-Fi
	cmd := exec.Command("netsh", "wlan", "add", "profile", "filename="+filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Ошибка при добавлении профиля Wi-Fi:", err)
		return
	}

	// Подключаемся к Wi-Fi
	cmd = exec.Command("netsh", "wlan", "connect", "name="+ssid)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Ошибка при подключении к Wi-Fi:", err)
		return
	}
	fmt.Println(string(output))
}
