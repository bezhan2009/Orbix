package src

import (
	"fmt"
	"goCmd/structs"
	"goCmd/utils"
)

func HandleUnknownCommandUtil(commandLower string, commands []structs.Command) {
	if !utils.ValidCommand(commandLower, commands) {
		printUnknown := fmt.Sprintf("'%s' is not recognized as an internal or external command,\noperable program or batch file.\n", commandLower)
		fmt.Printf(red(printUnknown))
		if suggestedCommand := suggestCommand(commandLower); suggestedCommand != "" {
			printSuggest := fmt.Sprintf("Did you mean: %s?\n", suggestedCommand)
			fmt.Printf(yellow(printSuggest))
		}
	}
}
