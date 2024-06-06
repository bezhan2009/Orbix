package Create

import (
	"fmt"
	"goCmd/commands/commandsWithSignaiture/Create/utils"
)

func File() (string, error) {
	var name string
	fmt.Print("Введите названия для файла:")
	fmt.Scan(&name)

	if name == "debug.txt" {
		fmt.Println("PermissionDenied: You cannot write, delete or create a debug.txt file")
		return name, nil
	}

	errExisting := utils.IsExists(name)

	if errExisting == nil {
		return name, errExisting
	}

	name, err := utils.CreateFile(name)

	if err != nil {
		return name, err
	} else {
		return name, err
	}
}
