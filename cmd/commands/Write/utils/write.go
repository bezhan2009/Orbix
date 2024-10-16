package utils

import (
	"goCmd/cmd/commands/Read/utils"
	"os"
)

func WriteFile(name string, data string) error {
	errOpening := IsExists(name)
	if errOpening != nil {
		return errOpening
	}

	oldData, errReadFile := utils.File(name)

	var err error
	if errReadFile == nil {
		oldData = append(oldData, []byte(data)...)
		err = os.WriteFile(name, oldData, 0666)
	} else {
		return err
	}

	return err
}
