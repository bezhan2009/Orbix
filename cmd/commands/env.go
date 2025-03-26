package commands

import (
	"errors"
	"fmt"
	"goCmd/internal/OS"
	"goCmd/system"
	"goCmd/system/errs"
)

func GetEnvVar(commandArgs []string) string {
	if len(commandArgs) < 1 {
		fmt.Println(system.Yellow("Usage: getenv <var name>"))
		return ""
	}

	value, err := OS.GetEnvVariable(commandArgs[0])
	if err != nil {
		if errors.Is(err, errs.VariableDoesNotExist) {
			fmt.Println(system.Red("The env variable " + commandArgs[0] + " does not exist"))
			return ""
		}

		fmt.Println(system.Red("Error getting env variable" + commandArgs[0] + ":" + err.Error()))
		return ""
	}

	return value
}

func SetEnvVar(commandArgs []string) {
	if len(commandArgs) < 2 {
		fmt.Println(system.Yellow("Usage: setenv <var name> <new value>"))
		return
	}

	err := OS.SetEnvVariable(commandArgs[0], commandArgs[1])
	if err != nil {
		fmt.Println(system.Red("Error setting env variable" + commandArgs[0] + ":" + err.Error()))
		return
	}

	fmt.Println("Successfully set env variable" + commandArgs[0] + ":" + commandArgs[1])
}
