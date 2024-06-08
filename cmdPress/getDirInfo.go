package cmdPress

func CmdDir(dir string) string {
	var count uint8
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
	var count uint8
	var user string

	for i := 0; i < len(dir); i++ {

		if dir[i] == '\\' {
			count += 1
			continue
		}

		if count > 1 && count < 3 {
			user += string(dir[i])
		}
	}

	return user
}

func CMDGetDirWithOutApp(dir string) string {
	var count uint8
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
