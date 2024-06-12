package wifiUtils

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"goCmd/Network/wifiUtils/commands/NetworkScan"
	"goCmd/Network/wifiUtils/commands/Send"
	"goCmd/Network/wifiUtils/commands/WiFi"
	"goCmd/commands/commandsWithoutSignature/Clean"
	"strings"
)

var commands = []string{
	"scanwifi", "connectwifi", "hackwifi", "networkscan", "sendSMS", "exit",
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
			Clean.Screen()
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
