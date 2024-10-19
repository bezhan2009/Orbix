package utils

import "os"

func IsExists(name string) error {
	file, err := os.Open(name)
	defer func() {
		err = file.Close()
		if err != nil {
			return
		}
	}()
	return err
}
