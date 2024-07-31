package DF

import (
	"fmt"
	"os"
	"path/filepath"
)

// DeleteFolder удаляет директорию и все ее содержимое
func DeleteFolder(commandArgs []string) (bool, error) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: df <foldername>")
		return false, nil
	}
	folderName := commandArgs[0]
	// Проверяем, существует ли директория
	info, err := os.Stat(folderName)
	if os.IsNotExist(err) {
		return false, fmt.Errorf("folder does not exist: %s", folderName)
	}

	// Проверяем, что это директория
	if !info.IsDir() {
		return false, fmt.Errorf("%s is not a directory", folderName)
	}

	// Рекурсивное удаление всех файлов и поддиректорий
	err = filepath.Walk(folderName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Удаляем файлы и пустые директории
		if info.IsDir() {
			return os.Remove(path)
		}

		return os.Remove(path)
	})

	if err != nil {
		return false, err
	}

	// Удаляем саму директорию
	return true, os.Remove(folderName)
}
