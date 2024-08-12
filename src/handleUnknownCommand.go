package src

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/structs"
	"goCmd/utils"
)

func HandleUnknownCommandUtil(commandLower, commandLine string, commands []structs.Command) {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	if !utils.ValidCommand(commandLower, commands) {
		printUnknown := fmt.Sprintf("'%s' is not recognized as an internal or external command,\noperable program or batch file.\n", commandLine)
		fmt.Printf(red(printUnknown))
		if suggestedCommand := suggestCommand(commandLower); suggestedCommand != "" {
			printSuggest := fmt.Sprintf("Did you mean: %s?\n", suggestedCommand)
			fmt.Printf(yellow(printSuggest))
		}
	}
}
