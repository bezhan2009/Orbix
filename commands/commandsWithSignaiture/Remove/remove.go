package Remove

import (
	"fmt"
	"os"
)

func File() (error, string) {
	var name string
	fmt.Print("Введите названия файла для удаления:")
	fmt.Scan(&name)

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
