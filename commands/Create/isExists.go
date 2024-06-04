package Create

import "os"

func IsExists(name string) error {
	_, err := os.Open(name)
	return err
}