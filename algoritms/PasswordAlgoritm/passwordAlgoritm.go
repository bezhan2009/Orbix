package PasswordAlgoritm

func Usage(password string, encrypt bool) string {
	if encrypt {
		encryptedPassword := EncryptPassword(password)
		return encryptedPassword
	} else {
		decrypt := DecryptPassword(password)
		return decrypt
	}
}
