package ORPXI

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Start(shablonName string) error {
	shablonName = strings.TrimSpace(shablonName)

	file, err := os.OpenFile(shablonName, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // Игнорировать пустые строки
		}
		CMD(line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return nil
}
