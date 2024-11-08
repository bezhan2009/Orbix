package fcommands

import (
	"errors"
	"goCmd/system"
	"os"
	"strings"
)

func CreateFile(name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return name, errors.New("NameError: Name cannot be empty")
	}

	// Проверка условия
	if system.OrbixFileNames[strings.TrimSpace(name)] == 1 {
		// Создание файла с правами доступа 0600 (только для владельца)
		file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			return name, err
		}
		defer func() {
			err = file.Close()
			if err != nil {
				return
			}
		}()

		return name, nil
	}

	// Если условие не выполнено, создаем файл с обычными правами
	file, err := os.Create(name)
	if err != nil {
		return name, err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			return
		}
	}()

	return name, nil
}
