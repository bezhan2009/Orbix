package wifiUtils

import (
	"github.com/c-bata/go-prompt"
	"strings"
)

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
		{Text: "sendMSG", Description: "Отправка сообщения на все ПК с подключением к этой сети"},
		{Text: "clean", Description: "очищает экран"},
		{Text: "exit", Description: "Выход из программы"},
	}
	return prompt.FilterHasPrefix(suggestions, args[0], true)
}
