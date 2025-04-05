package utils

import (
	"goCmd/cmd/commands/Read/utils"
	"log"
	"os"
)

func WriteFile(name string, data string) error {
	const op = "utils.WriteFile"

	errOpening := IsExists(name)
	if errOpening != nil {
		return errOpening
	}

	oldData, errReadFile := utils.File(name)

	var err error
	if errReadFile == nil {
		oldData = append(oldData, []byte(data)...)
		err = os.WriteFile(name, oldData, 0666)
		if err != nil {
			log.Printf("[%s] Error writing to file: %s", op, err)
		}
	} else {
		return err
	}

	return err
}
