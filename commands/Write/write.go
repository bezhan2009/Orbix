package Write

import (
	"goCmd/commands/Read"
	"os"
)

func File(name string, data string) error {
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
