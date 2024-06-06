package validators

import "goCmd/validators/utils"

func Password(password string) bool {
	var specialSymbols []string

	specialSymbols = utils.GetValidateSymbols()

	for i := 0; i < len(specialSymbols); i++ {
		for j := range password {
			if specialSymbols[i] == string(rune(j)) {
				return false
			}
		}
	}

	return true
}
