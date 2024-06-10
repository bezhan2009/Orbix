package wifiUtils

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"goCmd/Network/wifiUtils/commands/NetworkScan"
	"goCmd/Network/wifiUtils/commands/Send"
	"goCmd/Network/wifiUtils/commands/WiFi"
	"os"
	"os/exec"
	"strings"
)

var commands = []string{
	"scanwifi", "connectwifi", "hackwifi", "networkscan", "sendSMS", "exit",
}

func Screen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func autoComplete(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()
	if text == "" {
		return []prompt.Suggest{}
	}

	args := strings.Fields(text)
	if len(args) < 2 && args[0] == "connectwifi" || len(args) < 2 && args[0] == "hackwifi" {
		networks := getAvailableNetworks()
		suggestions := make([]prompt.Suggest, len(networks))
		for i, network := range networks {
			suggestions[i] = prompt.Suggest{Text: network, Description: "Доступная Wi-Fi сеть"}
		}
		return suggestions
	}

	suggestions := []prompt.Suggest{
		{Text: "scanwifi", Description: "Сканирование доступных Wi-Fi сетей"},
		{Text: "connectwifi", Description: "Подключение к Wi-Fi сети"},
		{Text: "hackwifi", Description: "Попытка взлома сети Wi-Fi"},
		{Text: "networkscan", Description: "Сканирование сети и получение информации об устройствах"},
		{Text: "sendSMS", Description: "Отправка SMS сообщения"},
		{Text: "sendMSG", Description: "Отправка сообщения на все ПК с подключением к этой сети"},
		{Text: "clean", Description: "очищает экран"},
		{Text: "exit", Description: "Выход из программы"},
	}
	return prompt.FilterHasPrefix(suggestions, args[0], true)
}

func Start() {
	fmt.Println("Добро пожаловать в утилиту для сетевого взаимодействия!")
	for {
		t := prompt.Input("> ", autoComplete)
		args := strings.Fields(t)
		if len(args) == 0 {
			continue
		}
		switch args[0] {
		case "help":
			showHelp()
		case "clean":
			Screen()
		case "scanwifi":
			WiFi.Scan()
		case "connectwifi":
			if len(args) < 3 {
				fmt.Println("Использование: connectwifi <SSID> <password>")
				continue
			}
			WiFi.Connect(args[1], args[2])
		case "hackwifi":
			if len(args) < 3 {
				fmt.Println("Использования: hackwifi <SSID> <attempts>")
				continue
			}
			WiFi.AttemptConnectWithGeneratedPasswords(args[1], args[2])
		case "networkscan":
			NetworkScan.WiFi()
		case "sendsms":
			if len(args) < 3 {
				fmt.Println("Использование: sendSMS <номер> <сообщение>")
				continue
			}
			Send.SMS(args[1], strings.Join(args[2:], " "))
		case "sendMSG":
			if len(args) < 3 {
				fmt.Println("Использование: sendMSG <пользователь> <сообщение>")
				continue
			}
			Send.Message(args[1], strings.Join(args[2:], " "))
		case "exit":
			fmt.Println("Выход из программы.")
			return
		default:
			fmt.Println("Неизвестная команда:", args[0])
		}
	}
}

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

func showHelp() {
	helpText := `
Доступные команды:
- scanwifi: Сканирование доступных Wi-Fi сетей
- connectwifi: Подключение к Wi-Fi сети
- hackwifi: Попытка взлома сети Wi-Fi
- networkscan: Сканирование сети и получение информации об устройствах
- sendSMS: Отправка SMS сообщения
- help: Показать справку по использованию программы
- sendMSG: Отправка сообщения на все ПК с подключением к этой сети
- exit: Выход из программы
`
	fmt.Println(helpText)
}
