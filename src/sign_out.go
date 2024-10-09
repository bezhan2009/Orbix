package src

import (
	"fmt"
	"goCmd/structs"
	"goCmd/system"
	"os"
)

func SignOutUtil(username string, systemPath string, sd *system.AppState, sessionPrefix string) {
	sd.DeleteSession(sessionPrefix)
	err := os.Chdir(systemPath)
	if err != nil {
		fmt.Println(red(fmt.Sprintf("Error when changing the path: %v", err)))
	}

	RemoveUserFromRunningFile(username)
	Orbix("", true, structs.RebootedData{}, sd)
}
