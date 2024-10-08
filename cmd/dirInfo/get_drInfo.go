package dirInfo

import (
	"goCmd/system"
	"os"
)

func CmdDir(dir string) string {
	if system.OperationSystem != "windows" {
		return dir
	}

	var count uint16
	var dirC string

	for i := 0; i < len(dir); i++ {

		if dir[i] == '\\' {
			count += 1
		}

		if count > 2 {
			dirC += string(dir[i])
		}
	}

	return dirC
}

func CmdUser(dir string) string {
	if system.OperationSystem == "linux" {
		// В Linux просто возвращаем текущего пользователя
		return os.Getenv("USER")
	} else {
		var count uint16
		var user string

		for i := 0; i < len(dir); i++ {
			if dir[i] == '\\' {
				count++
				continue
			}

			if count > 1 && count < 3 {
				user += string(dir[i])
			}
		}

		return user
	}
}

func CMDGetDirWithOutApp(dir string) string {
	var count uint16
	var dirC string

	for i := 0; i < len(dir); i++ {

		if dir[i] == '\\' {
			count += 1
		}

		if count > 2 {
			dirC += string(dir[i])
		}
	}

	return dirC
}
