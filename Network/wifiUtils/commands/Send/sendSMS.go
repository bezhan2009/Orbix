package Send

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"os/exec"
)

func SMS(number, message string) {
	client := resty.New()
	resp, err := client.R().
		SetFormData(map[string]string{
			"To":      number,
			"Message": message,
			// Здесь нужно указать необходимые данные для API
		}).
		Post("https://api.sms-provider.com/send")

	if err != nil {
		fmt.Println("Ошибка при отправке SMS:", err)
		return
	}
	fmt.Println("Ответ от сервера:", resp)
}

func Message(username, message string) {
	cmd := exec.Command("msg", username, message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Ошибка при отправке сообщения: %v\n", err)
		return
	}
	fmt.Println("Сообщение успешно отправлено.")
	fmt.Println(string(output))
}
