package Create

import (
	"fmt"
	"goCmd/commands/Create/utils"
)

func File() (error, string) {
	var name string
	fmt.Print("Введите названия для файла:")
	fmt.Scan(&name)

	if name == "debug.txt" {
		panic("PermissionDenied: You cannot write, delete or create a debug.txt file")
	}

	errExisting := IsExists(name)

	if errExisting == nil {
		return errExisting, name
	}

	name, err := utils.CreateFile(name)

	if err != nil {
		return err, name
	} else {
		return nil, name
	}
}
