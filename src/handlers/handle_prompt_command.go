package handlers

import (
	"fmt"
	"goCmd/src/service"
	utils2 "goCmd/system"
	"goCmd/utils"
	"os"
	"strings"
)

func HandlePromptCommand(commandArgs []string, prompt *string) {
	if len(commandArgs) < 1 {
		printHint1 := fmt.Sprintf("prompt <name_prompt> <color(default: green)>\n")
		printHint2 := fmt.Sprintf("to delete prompt enter:\n")
		printHint3 := fmt.Sprintf("prompt delete\n")

		fmt.Print(utils2.Yellow(printHint1))
		fmt.Print(utils2.Yellow(printHint2))
		fmt.Print(utils2.Yellow(printHint3))

		return
	}

	namePrompt := commandArgs[0]

	var namePromptWithColor string
	var isValid bool

	if len(commandArgs) > 1 {
		validColors := []string{"yellow", "green", "blue", "magenta", "cyan", "red"}
		isValid = utils.IsValid(commandArgs[1], validColors)
	}

	if len(commandArgs) > 1 && isValid {
		namePromptWithColor = fmt.Sprintf("%s, %s", namePrompt, commandArgs[1])
	} else {
		namePromptWithColor = fmt.Sprintf("%s, %s", namePrompt, "green")
	}

	if namePrompt != "delete" {
		namePrompt = strings.TrimSpace(namePromptWithColor)
		*prompt = namePromptWithColor
		if len(commandArgs) > 1 {
			service.AnimatedPrint(fmt.Sprintf("Prompt set to: %s\n", commandArgs[0]), commandArgs[1])
		} else {
			service.AnimatedPrint(fmt.Sprintf("Prompt set to: %s\n", commandArgs[0]), "green")
		}
	} else {
		*prompt, _ = os.Getwd()
		service.AnimatedPrint(fmt.Sprintf("Prompt set to: %s\n", *prompt), "green")
		*prompt = ""
	}
}
