package bun

import (
	"goCmd/safe/block/bun/utils"
	"os"
)

func UserGoCMD(command string, intentionallyBan bool) bool {
	isFile := utils.BunGoCMD(command)

	if isFile || intentionallyBan {
		os.Create("safe/block/bun/component/bunnedUser.json")

		isFile = true
	}

	return isFile
}
