package utils

import "os"

func IsExists(name string) error {
	file, err := os.Open(name)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			return
		}
	}(file)
	return err
}
