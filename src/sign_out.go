package src

import (
	"goCmd/structs"
	"goCmd/system"
	"log"
	"os"
)

func SignOutUtil(username string, systemPath string, sd *system.AppState, sessionPrefix string) {
	sd.DeleteSession(sessionPrefix)
	err := os.Chdir(systemPath)
	if err != nil {
		log.Fatalf("Error when changing the path: %v", err)
	}

	RemoveUserFromRunningFile(username)
	Orbix("", true, structs.RebootedData{}, sd)
}
