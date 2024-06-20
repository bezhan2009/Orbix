package darwin

import (
	"fmt"
	"os/exec"
)

func Message(username, message string) {
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "%s"`, message, username))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Ошибка при отправке сообщения: %v\n", err)
		return
	}
	fmt.Println("Сообщение успешно отправлено.")
	fmt.Println(string(output))
}
