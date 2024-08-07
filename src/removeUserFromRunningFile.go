package src

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func removeUserFromRunningFile(username string) {
	runningPath := filepath.Join(Absdir, "running.txt")

	sourceRunning, err := os.ReadFile(runningPath)
	if err != nil {
		fmt.Printf("Ошибка чтения файла running.txt: %v\n", err)
		return
	}

	lines := strings.Split(string(sourceRunning), "\n")
	var updatedLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != username {
			updatedLines = append(updatedLines, line)
		}
	}

	newContent := strings.Join(updatedLines, "\n")
	err = os.WriteFile(runningPath, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("Ошибка записи в файл running.txt: %v\n", err)
		return
	}
}
