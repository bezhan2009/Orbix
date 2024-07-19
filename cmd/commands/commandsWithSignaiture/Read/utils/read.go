package utils

import (
	"os"
)

func File(name string) ([]byte, error) {
	errOpening := IsExists(name)

	if errOpening != nil {
		return nil, errOpening
	}

	data, err := os.ReadFile(name)

	return data, err
}
