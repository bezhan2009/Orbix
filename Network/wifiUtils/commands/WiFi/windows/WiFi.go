package windows

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
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
        <safe>
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
        </safe>
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

	// Проверка доступа к интернету
	for i := 0; i < 3; i++ { // Повторяем проверку несколько раз
		time.Sleep(1 * time.Second) // Пауза перед каждой проверкой
		resp, err := http.Get("http://clients3.google.com/generate_204")
		if err == nil && resp.StatusCode == http.StatusNoContent {
			fmt.Println("Подключение успешно")
			return true
		}
	}

	fmt.Println("Нет доступа к интернету")
	return false
}
