package linux

import (
	"fmt"
	"os/exec"
)

func Message(username, message string) {
	cmd := exec.Command("notify-send", username, message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Ошибка при отправке сообщения: %v\n", err)
		return
	}

	fmt.Println("Сообщение успешно отправлено.")
	fmt.Println(string(output))
}
