package ORPXI

import (
	"fmt"
	"goCmd/validators"
	"os"
)

func Password() {
	var password string
	var passwordByte *string
	fmt.Println("Enter Password: ")
	fmt.Scan(password)
	isValid := validators.Password(password)

	if isValid {
		passwordByte = &password
	} else {
		fmt.Println("Password is Invalid")
		return
	}

	Copy(password, passwordByte)
	bytes := []byte(password)
	password = string(bytes)
	fmt.Println("Your password(file):", password)
	os.Create(password)
	os.WriteFile(password, []byte(*passwordByte), os.ModePerm)
}

func Copy(origin string, dest *string) {
	*dest = origin
}
