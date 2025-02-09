package handlers

import (
	"fmt"
	"goCmd/src"
	"goCmd/src/service"
	utils2 "goCmd/system"
	"goCmd/utils"
)

func HandleUnknownCommandUtil(command, commandLower string,
	commands map[string]struct{}) {
	if !utils.ValidCommandFast(commandLower, commands) && src.UnknownCommandsCounter == 0 && !utils2.CheckPackageExists(command) {
		printUnknown := fmt.Sprintf("'%s' is not recognized as an internal or external command,\noperable program or batch file.\n", command)
		fmt.Printf(utils2.Red(printUnknown))
		if suggestedCommand := service.SuggestCommand(commandLower); suggestedCommand != "" {
			printSuggest := fmt.Sprintf("Did you mean: %s?\n", suggestedCommand)
			fmt.Printf(utils2.Yellow(printSuggest))
		}
	}

	src.UnknownCommandsCounter = src.UnknownCommandsCounter + 1
}
