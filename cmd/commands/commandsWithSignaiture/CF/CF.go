package CF

import (
	"fmt"
	"os"
)

func CreateFolder(commandArgs []string) (bool, error) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: cf <dir name>")
		return false, nil
	}
	err := os.Mkdir(commandArgs[0], 0755) // 0755 - права доступа к директории
	if err != nil {
		return false, err
	}
	return true, nil
}
