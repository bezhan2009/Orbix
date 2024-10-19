package src

import (
	"fmt"
	"goCmd/utils"
	"os"
	"strings"
)

func handlePromptCommand(commandArgs []string, prompt *string) {
	if len(commandArgs) < 1 {
		printHint1 := fmt.Sprintf("prompt <name_prompt> <color(default: green)>\n")
		printHint2 := fmt.Sprintf("to delete prompt enter:\n")
		printHint3 := fmt.Sprintf("prompt delete\n")

		fmt.Print(yellow(printHint1))
		fmt.Print(yellow(printHint2))
		fmt.Print(yellow(printHint3))

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
			animatedPrint(fmt.Sprintf("Prompt set to: %s\n", commandArgs[0]), commandArgs[1])
		} else {
			animatedPrint(fmt.Sprintf("Prompt set to: %s\n", commandArgs[0]), "green")
		}
	} else {
		*prompt, _ = os.Getwd()
		animatedPrint(fmt.Sprintf("Prompt set to: %s\n", *prompt), "green")
		*prompt = ""
	}
}
