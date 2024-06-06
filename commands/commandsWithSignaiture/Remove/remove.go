package Remove

import (
	"fmt"
	"os"
)

func File(commandArgs []string) (error, string) {
	if len(commandArgs) < 1 {
		fmt.Println("Использования: remove <файл>")
	}

	var name string

	name = commandArgs[0]

	if name == "debug.txt" {
		fmt.Println("PermissionDenied: You cannot write, delete or create a debug.txt file")
	}

	errExisting := IsExists(name)

	if errExisting != nil {
		return errExisting, name
	}

	err := os.Remove(name)
	if err != nil {
		return err, name
	} else {
		return nil, name
	}
}
