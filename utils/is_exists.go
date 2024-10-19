package utils

import "os"

func IsExists(name string) (*os.File, error) {
	f, err := os.Open(name)
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			return
		}
	}(f)

	return f, err
}
