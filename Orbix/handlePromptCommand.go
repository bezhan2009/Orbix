package Orbix

import (
	"fmt"
	"os"
	"strings"
)

func handlePromptCommand(commandArgs []string, prompt *string) {
	if len(commandArgs) < 1 {
		animatedPrint("prompt <name_prompt>\n")
		animatedPrint("to delete prompt enter:\n")
		animatedPrint("prompt delete\n")
		return
	}

	namePrompt := commandArgs[0]

	if namePrompt != "delete" {
		namePrompt = strings.TrimSpace(namePrompt)
		*prompt = namePrompt
		animatedPrint(fmt.Sprintf("Prompt set to: %s\n", *prompt))
	} else {
		*prompt, _ = os.Getwd()
		animatedPrint(fmt.Sprintf("Prompt set to: %s\n", *prompt))
		*prompt = ""
	}
}
