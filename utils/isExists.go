package utils

import "os"

func IsExists(name string) (*os.File, error) {
	f, err := os.Open(name)
	return f, err
}
