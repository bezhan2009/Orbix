package commandsWithSignaiture

import (
	"io"
	"os"
	"strings"

	"github.com/atotto/clipboard"
)

func File(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Если целевая директория - "buffer", используем буфер обмена
	if strings.ToLower(dst) == "buffer" {
		// Читаем содержимое файла в память
		data, err := io.ReadAll(sourceFile)
		if err != nil {
			return err
		}

		// Помещаем содержимое в буфер обмена
		if err := writeToClipboard(data); err != nil {
			return err
		}

		return nil
	}

	// Иначе, создаем целевой файл и копируем содержимое
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	err = destinationFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func writeToClipboard(data []byte) error {
	// Записываем данные в буфер обмена
	return clipboard.WriteAll(string(data))
}
