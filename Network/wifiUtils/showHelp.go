package wifiUtils

import "fmt"

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
