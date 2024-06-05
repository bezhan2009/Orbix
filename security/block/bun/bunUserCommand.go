package bun

import (
	"goCmd/security/block/bun/utils"
	"os"
)

func UserGoCMD(command string, intentionallyBan bool) bool {
	isFile := utils.BunGoCMD(command)

	if isFile || intentionallyBan {
		os.Create("security/block/bun/component/bunnedUser.json")

		isFile = true
	}

	return isFile
}
