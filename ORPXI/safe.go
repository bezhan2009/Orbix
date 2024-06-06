package ORPXI

import (
	"bufio"
	"fmt"
	"goCmd/commands/Create/utils"
	"goCmd/validators"
	"os"
)

func Password() {
	reader := bufio.NewReader(os.Stdin)

	var password string
	var passwordByte *string

	fmt.Print("Enter Password: ")

	password, _ = reader.ReadString('\n')

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

	password, err := utils.CreateFile(password)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(password, []byte(*passwordByte), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}

func Copy(origin string, dest *string) {
	*dest = origin
}
