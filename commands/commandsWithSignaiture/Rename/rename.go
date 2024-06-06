package Rename

import (
	"fmt"
	"os"
)

func Rename(commandArgs []string) error {
	if len(commandArgs) < 2 {
		fmt.Println("Использование: rename <файл> <новое имя для файла>")
		return nil
	}

	err := os.Rename(commandArgs[0], commandArgs[1])
	if err == nil {
		return err
	} else {
		return err
	}
}
