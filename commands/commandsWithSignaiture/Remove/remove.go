package Remove

import (
	"fmt"
	"goCmd/commands/commandsWithSignaiture/Remove/utils"
)

func File(commandArgs []string) (string, error) {
	if len(commandArgs) < 1 {
		fmt.Println("Использования: remove <файл>")
	}

	var name string

	name = commandArgs[0]

	if name == "debug.txt" {
		fmt.Println("PermissionDenied: You cannot write, delete or create a debug.txt file")
		return name, nil
	}

	errExisting := utils.IsExists(name)

	if errExisting != nil {
		return name, errExisting
	}

	return utils.RemoveFile(name)
}
