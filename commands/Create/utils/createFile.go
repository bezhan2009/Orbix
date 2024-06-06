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
		_, err := os.Create(name)
		if err != nil {
			return name, err
		} else {
			return name, nil
		}
	}
}
