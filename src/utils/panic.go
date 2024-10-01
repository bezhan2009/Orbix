package utils

import "goCmd/structs"

func PanicOnError(executeCommand structs.ExecuteCommandFuncParams) {
	if len(executeCommand.CommandArgs) > 0 {
		panic(executeCommand.CommandArgs[0])
	} else {
		panic("executeCommand: command panic!!!")
	}
}
