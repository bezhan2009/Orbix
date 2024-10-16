package Read

import (
	"fmt"
	"github.com/fatih/color"
	utils2 "goCmd/cmd/commands/Read/utils"
)

func File(commandArgs []string) error {
	yellow := color.New(color.FgYellow).SprintFunc()
	if len(commandArgs) < 1 {
		fmt.Println(yellow(fmt.Sprint("Usage: read <file>")))
		return nil
	}

	nameFileForRead := commandArgs[0]

	dataRead, errReading := utils2.File(nameFileForRead)
	if errReading != nil {
		return errReading
	} else {
		fmt.Println(yellow(string(dataRead)))
		return nil
	}
}
