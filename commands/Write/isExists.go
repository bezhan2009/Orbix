package Write

import "os"

func IsExists(name string) error {
	_, err := os.Open(name)
	return err
}
