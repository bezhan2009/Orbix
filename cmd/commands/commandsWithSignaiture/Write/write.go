package Write

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/Write/utils"
	"strings"
)

func File(commandArgs []string) error {
	if len(commandArgs) < 2 {
		fmt.Println("Usage: write <file> <data>")
		return nil
	}

	nameFileForWrite := commandArgs[0]
	data := strings.Join(commandArgs[1:], " ")

	errWriting := utils.WriteFile(nameFileForWrite, data+"\n")

	if errWriting != nil {
		return errWriting
	} else {
		fmt.Printf("Successfully wrote data to file %s\n", nameFileForWrite)
		return nil
	}
}
