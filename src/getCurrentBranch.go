package src

import (
	"bytes"
	"os/exec"
)

func GetCurrentGitBranch() (string, error) {
	// Создаем команду для выполнения
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	// Создаем буфер для захвата вывода
	var out bytes.Buffer
	cmd.Stdout = &out

	// Выполняем команду
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Получаем результат и убираем лишние пробелы
	branch := out.String()
	return branch[:len(branch)-1], nil
}
