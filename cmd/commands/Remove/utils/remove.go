package utils

import "os"

func RemoveFile(name string) (string, error) {
	err := os.Remove(name)
	if err != nil {
		return name, err
	} else {
		return name, nil
	}
}
