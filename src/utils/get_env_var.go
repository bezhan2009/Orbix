package utils

import (
	"fmt"
	"goCmd/cmd/commands"
	"goCmd/system"
)

func GetEnvVarUtil(commandArgs []string) {
	value := commands.GetEnvVar(commandArgs)
	if value != "" {
		fmt.Println(system.GreenBold(fmt.Sprintf("%s=%s ", commandArgs[0], value)))
	}
}
