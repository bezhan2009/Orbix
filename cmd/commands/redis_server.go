package commands

import (
	"fmt"
	"log"
	"os/exec"
)

func StartRedisServer() {
	// Команда для запуска Ubuntu и выполнения команды redis-server
	cmd := exec.Command("wsl", "-d", "Ubuntu", "redis-server")

	// Захват вывода команды
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Ошибка при запуске redis-server через WSL: %v", err)
	}

	// Вывод результата команды в консоль
	fmt.Printf("Output command: %s\n", string(output))
}
