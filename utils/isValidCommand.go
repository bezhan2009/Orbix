package utils

func ValidCommand(command string, commands []string) bool {
	for i := 0; i < len(commands); i++ {
		if command == commands[i] {
			return true
		}
	}

	return false
}
