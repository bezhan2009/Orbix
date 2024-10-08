package src

import (
	"fmt"
	"os"
	"strings"
)

func handlePromptCommand(commandArgs []string, prompt *string) {
	if len(commandArgs) < 1 {
		printHint1 := fmt.Sprintf("prompt <name_prompt>\n")
		printHint2 := fmt.Sprintf("to delete prompt enter:\n")
		printHint3 := fmt.Sprintf("prompt delete\n")
		fmt.Print(yellow(printHint1))
		fmt.Print(yellow(printHint2))
		fmt.Print(yellow(printHint3))
		return
	}

	namePrompt := commandArgs[0]

	if namePrompt != "delete" {
		namePrompt = strings.TrimSpace(namePrompt)
		*prompt = namePrompt
		animatedPrint(fmt.Sprintf("Prompt set to: %s\n", *prompt), "green")
	} else {
		*prompt, _ = os.Getwd()
		animatedPrint(fmt.Sprintf("Prompt set to: %s\n", *prompt), "green")
		*prompt = ""
	}
}
