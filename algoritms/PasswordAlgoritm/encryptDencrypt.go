package PasswordAlgoritm

import "goCmd/validators/utils"

func EncryptPassword(password string) string {
	alphabetSymbols := utils.GetAlphabetSymbols()
	encrypted := ""

	for _, char := range password {
		if symbol, exists := alphabetSymbols[char]; exists {
			encrypted += symbol
		} else {
			encrypted += string(char)
		}
	}

	return encrypted
}

func DecryptPassword(encrypted string) string {
	alphabetSymbols := utils.GetAlphabetSymbols()
	decrypted := ""

	// Создание обратного маппинга для дешифрования
	reverseMapping := make(map[string]rune)
	for k, v := range alphabetSymbols {
		reverseMapping[v] = k
	}

	for _, char := range encrypted {
		if original, exists := reverseMapping[string(char)]; exists {
			decrypted += string(original)
		} else {
			decrypted += string(char)
		}
	}

	return decrypted
}
