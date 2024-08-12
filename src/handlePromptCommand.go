package src

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

func handlePromptCommand(commandArgs []string, prompt *string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	if len(commandArgs) < 1 {
		printHint1 := fmt.Sprintf("prompt <name_prompt>\n")
		printHint2 := fmt.Sprintf("to delete prompt enter:\n")
		printHint3 := fmt.Sprintf("prompt delete\n")
		animatedPrint(yellow(printHint1))
		animatedPrint(yellow(printHint2))
		animatedPrint(yellow(printHint3))
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
