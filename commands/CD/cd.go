package CD

import (
	"fmt"
	"os"
)

func ChangeDirectory(path string) error {
	err := os.Chdir(path)
	if err != nil {
		return fmt.Errorf("не удалось сменить директорию: %v", err)
	}
	return nil
}
