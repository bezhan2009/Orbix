package utils

import (
	"errors"
	"fmt"
	"os"
)

func CreateFile(name string) (string, error) {
	if name == "" {
		fmt.Println("nameError: Name cannot be empty!!!")
		return name, errors.New("nameError: Name cannot be empty")
	} else {
		file, err := os.Create(name)
		if err != nil {
			return name, err
		}
		defer file.Close()
		return name, nil
	}
}
