package Read

import (
	"os"
)

func File(name string) ([]byte, error) {
	data, err := os.ReadFile(name)

	return data, err
}
