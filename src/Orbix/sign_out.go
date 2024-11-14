package Orbix

import (
	"fmt"
	"goCmd/src/user"
	"goCmd/structs"
	"goCmd/system"
	"os"
)

func SignOutUtil(username string, systemPath string, sd *system.AppState, sessionPrefix string) {
	sd.DeleteSession(sessionPrefix)
	err := os.Chdir(systemPath)
	if err != nil {
		fmt.Println(system.Red(fmt.Sprintf("Error when changing the path: %v", err)))
	}

	user.DeleteUserFromRunningFile(username)
	Orbix("",
		true,
		structs.RebootedData{},
		sd)
}
