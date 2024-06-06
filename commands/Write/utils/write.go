package utils

import (
	"goCmd/commands/Read"
	"os"
)

func WriteFile(name string, data string) error {
	errOpening := IsExists(name)
	if errOpening != nil {
		return errOpening
	}

	oldData, errReadFile := Read.File(name)

	var err error
	if errReadFile == nil {
		oldData = append(oldData, []byte(data)...)
		err = os.WriteFile(name, oldData, 0666)
	} else {
		return err
	}

	return err
}
