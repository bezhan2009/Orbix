package utils

import "os"

func BunGoCMD(command string) bool {
	bunCommand := false
	for i := 0; i < len(command); i++ {
		if command[i] == '.' {
			bunCommand = true
		}
	}

	if bunCommand == true {
		_, err := os.Open(command)
		if err != nil {
			bunCommand = false
		}
	}

	return bunCommand
}
