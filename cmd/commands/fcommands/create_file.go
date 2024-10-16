package fcommands

import (
	"errors"
	"os"
)

func CreateFile(name string) (string, error) {
	if name == "" {
		return name, errors.New("nameError: Name cannot be empty")
	} else {
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
}
