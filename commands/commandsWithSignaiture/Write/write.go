package Write

import (
	"fmt"
	"goCmd/commands/Write/utils"
	"goCmd/debug"
	"strings"
)

func File(command string, commandArgs []string) {
	if len(commandArgs) < 2 {
		fmt.Println("Использование: write <файл> <данные>")
		return
	}

	nameFileForWrite := commandArgs[0]

	data := strings.Join(commandArgs[1:], " ")

	if nameFileForWrite == "debug.txt" {
		debug.Commands(command, false)
		fmt.Println("PermissionDenied: You cannot write, delete or create a debug.txt file")
		return
	}

	errWriting := utils.WriteFile(nameFileForWrite, data+"\n")

	if errWriting != nil {
		debug.Commands(command, false)
		fmt.Println(errWriting)
	} else {
		debug.Commands(command, true)
		fmt.Printf("Мы успешно записали данные в файл %s\n", nameFileForWrite)
	}
}
