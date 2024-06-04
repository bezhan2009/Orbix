package Rename

import (
	"fmt"
	"os"
)

func Rename() error {
	var (
		oldPath string
		newPath string
	)

	fmt.Println("Введите старое имя файла:")
	fmt.Scan(&oldPath)
	fmt.Println("Введите новое имя файла:")
	fmt.Scan(&newPath)
	err := os.Rename(oldPath, newPath)
	if err == nil {
		return err
	} else {
		return err
	}
}
