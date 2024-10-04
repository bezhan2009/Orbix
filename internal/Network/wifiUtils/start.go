package wifiUtils

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	commands2 "goCmd/cmd/commands"
	"goCmd/internal/Network/wifiUtils/commands/NetworkScan"
	darwin2 "goCmd/internal/Network/wifiUtils/commands/Send/darwin"
	linux2 "goCmd/internal/Network/wifiUtils/commands/Send/linux"
	windows2 "goCmd/internal/Network/wifiUtils/commands/Send/windows"
	"goCmd/internal/Network/wifiUtils/commands/WiFi"
	"goCmd/internal/Network/wifiUtils/commands/WiFi/darwin"
	"goCmd/internal/Network/wifiUtils/commands/WiFi/linux"
	"goCmd/internal/Network/wifiUtils/commands/WiFi/windows"
	"goCmd/internal/OS"
	"strings"
)

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
			commands2.Screen()
		case "scanwifi":
			nameOS := OS.CheckOS()

			if nameOS == "windows" {
				windows.Scan()
			} else if nameOS == "linux" {
				linux.Scan()
			} else if nameOS == "darwin" {
				darwin.Scan()
			} else {
				fmt.Println("Unresolved OS")
			}
		case "connectwifi":
			if len(args) < 3 {
				fmt.Println("Использование: connectwifi <SSID> <password>")
				continue
			}
			nameOS := OS.CheckOS()
			if nameOS == "windows" {
				windows.Connect(args[1], args[2])
			} else if nameOS == "linux" {
				linux.Connect(args[1], args[2])
			} else if nameOS == "darwin" {
				darwin.Connect(args[1], args[2])
			} else {
				fmt.Println("Unresolved OS")
			}
		case "hackwifi":
			if len(args) < 3 {
				fmt.Println("Использования: hackwifi <SSID> <attempts>")
				continue
			}
			WiFi.AttemptConnectWithGeneratedPasswords(args[1], args[2])
		case "networkscan":
			NetworkScan.WiFi()
		case "sendMSG":
			if len(args) < 3 {
				fmt.Println("Использование: sendMSG <пользователь> <сообщение>")
				continue
			}

			nameOS := OS.CheckOS()
			if nameOS == "windows" {
				windows2.Message(args[1], strings.Join(args[2:], " "))
			} else if nameOS == "linux" {
				linux2.Message(args[1], strings.Join(args[2:], " "))
			} else if nameOS == "darwin" {
				darwin2.Message(args[1], strings.Join(args[2:], " "))
			} else {
				fmt.Println("Unresolved OS")
			}
		case "exit":
			fmt.Println("Выход из программы.")
			return
		default:
			fmt.Println("Неизвестная команда:", args[0])
		}
	}
}
