package bun

import (
	"goCmd/security/block/bun/utils"
	"os"
)

func BunCommand(command string) {
	isFile := utils.BunGoCMD(command)

	if isFile {
		_, err := os.Open("bunnedCommand.json")

		if err != nil {
			os.Create("bunnedCommand.json")
		}
	}
}
