package src

import (
	"fmt"
	"goCmd/structs"
	"goCmd/utils"
)

func HandleUnknownCommandUtil(commandLower, commandLine string, commands []structs.Command) {
	if !utils.ValidCommand(commandLower, commands) {
		fmt.Printf("'%s' is not recognized as an internal or external command,\noperable program or batch file.\n", commandLine)
		if suggestedCommand := suggestCommand(commandLower); suggestedCommand != "" {
			fmt.Printf("Did you mean: %s?\n", suggestedCommand)
		}
	}
}
