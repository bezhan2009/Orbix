package PasswordAlgoritm

import (
	"encoding/base64"
	"goCmd/validators/utils"
)

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

func XorEncryptDecrypt(data string, key []byte, encode bool) string {
	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ key[i%len(key)]
	}

	if encode {
		// Для шифрования: возвращаем base64
		return base64.StdEncoding.EncodeToString(result)
	}

	// Для дешифрования: вход - base64, выход - строка
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return data // fallback для обратной совместимости
	}

	result = make([]byte, len(decoded))
	for i := 0; i < len(decoded); i++ {
		result[i] = decoded[i] ^ key[i%len(key)]
	}
	return string(result)
}
