package validators

import (
	"goCmd/validators/utils"
	"unicode"
)

func Password(password string) bool {
	specialSymbols := utils.GetValidateSymbols()
	alphabetSymbols := utils.GetAlphabetSymbols()

	for _, char := range password {
		for _, sym := range specialSymbols {
			if string(char) == sym {
				return false
			}
		}

		if !unicode.Is(unicode.Latin, char) {
			return false
		}

		for _, sym := range alphabetSymbols {
			if string(char) == sym {
				return false
			}
		}
	}

	return true
}
