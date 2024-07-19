package darwin

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

func Scan() {
	cmd := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-s")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Ошибка при сканировании Wi-Fi:", err)
		return
	}
	fmt.Println(string(output))
}

func Connect(ssid, password string) bool {
	cmd := exec.Command("networksetup", "-setairportnetwork", "en0", ssid, password)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		fmt.Printf("Ошибка при подключении к Wi-Fi: %v\n", err)
		return false
	}

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
