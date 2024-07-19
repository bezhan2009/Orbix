package debug

import (
	"fmt"
	"os"
	"time"
)

func Commands(command string, isSuccess bool, args []string, user string, dir string) error {
	file, err := os.OpenFile("debug.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	timestamp := time.Now().Format(time.RFC3339)
	status := "False"
	if isSuccess {
		status = "True"
	}

	data := fmt.Sprintf("Timestamp: %s\nUser: %s\nDirectory: %s\nCommand: %s\nArguments: %v\nSuccess: %s\n\n",
		timestamp, user, dir, command, args, status)

	if _, err := file.WriteString(data); err != nil {
		return err
	}

	return nil
}
