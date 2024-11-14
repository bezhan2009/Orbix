package utils

import (
	"goCmd/cmd/commands"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
	"path/filepath"
)

func NeofetchUtil(executeCommand structs.ExecuteCommandFuncParams, user string, Commands []system.Command) {
	if system.OperationSystem == "windows" {
		commands.FetchNeofetch(user)
	} else {
		isValid := utils.ValidCommand(executeCommand.CommandLower, Commands)

		if !isValid {
			fullCommand := append([]string{executeCommand.Command}, executeCommand.CommandArgs...)
			err := utils.ExternalCommand(fullCommand)
			if err != nil {
				fullPath := filepath.Join(system.UserDir, executeCommand.Command)
				fullCommand[0] = fullPath
				_ = utils.ExternalCommand(fullCommand)
			}
		}
	}
}
