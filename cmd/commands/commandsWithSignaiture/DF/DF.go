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

	// Создаем срез для хранения директорий, чтобы удалить их позже
	var dirs []string

	// Рекурсивное удаление всех файлов и поддиректорий
	err = filepath.Walk(folderName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Сохраняем директории для удаления позже
		if info.IsDir() {
			dirs = append(dirs, path)
		} else {
			// Удаляем файлы сразу
			if err := os.Remove(path); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	// Удаляем директории в обратном порядке
	for i := len(dirs) - 1; i >= 0; i-- {
		if err := os.Remove(dirs[i]); err != nil {
			return false, err
		}
	}

	// Удаляем саму директорию
	return true, os.Remove(folderName)
}
