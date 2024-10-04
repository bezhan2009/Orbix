package utils

import (
	"goCmd/cmd/commands/commandsWithoutSignature/neofetch"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
	"path/filepath"
)

func NeofetchUtil(executeCommand structs.ExecuteCommandFuncParams, session *system.Session, Commands []structs.Command) {
	if system.OperationSystem == "windows" {
		neofetch.FetchNeofetch(session)
	} else {
		isValid := utils.ValidCommand(executeCommand.CommandLower, Commands)

		if !isValid {
			fullCommand := append([]string{executeCommand.Command}, executeCommand.CommandArgs...)
			err := utils.ExternalCommand(fullCommand)
			if err != nil {
				fullPath := filepath.Join(executeCommand.Dir, executeCommand.Command)
				fullCommand[0] = fullPath
				_ = utils.ExternalCommand(fullCommand)
			}
		}
	}
}
