package src

import (
	"goCmd/structs"
	"log"
	"os"
)

func SignOutUtil(username string, systemPath string) {
	err := os.Chdir(systemPath)
	if err != nil {
		log.Fatalf("Error when changing the path: %v", err)
	}

	removeUserFromRunningFile(username)
	Orbix("", true, structs.RebootedData{})
}
