package utils

import "os"

func IsExists(name string) error {
	file, err := os.Open(name)
	defer file.Close()
	return err
}
