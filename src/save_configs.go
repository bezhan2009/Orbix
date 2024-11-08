package src

import (
	"fmt"
	"goCmd/cmd/commands"
	ReadEnvUtil "goCmd/cmd/commands/Read/utils"
	"goCmd/pkg/algorithms/PasswordAlgoritm"
	"goCmd/system"
	"os"
	"strings"
)

func silenceOutput() (func(), error) {
	// Определяем путь для "черной дыры" в зависимости от ОС
	var nullDevice string
	if system.OperationSystem == "windows" {
		nullDevice = "nul"
	} else {
		nullDevice = "/dev/null"
	}

	// Открываем файл nullDevice для перенаправления вывода
	nullFile, err := os.OpenFile(nullDevice, os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Сохраняем оригинальные os.Stdout и os.Stderr
	origStdout := os.Stdout
	origStderr := os.Stderr

	// Перенаправляем их на nullFile
	os.Stdout = nullFile
	os.Stderr = nullFile

	// Возвращаем функцию для восстановления
	return func() {
		os.Stdout = origStdout
		os.Stderr = origStderr
		nullFile.Close()
	}, nil
}

func SaveVars() {
	restoreOutput, err := silenceOutput() // Отключаем вывод
	if err != nil {
		fmt.Println("Error while disabling output:", err)
		return
	}
	defer restoreOutput() // Восстанавливаем вывод в конце

	err = commands.ChangeDirectory(Absdir)
	if err != nil {
		fmt.Println(red(err))
		return
	}

	execLtCommand("delete user.env")
	execLtCommand("create user.env")

	for key, value := range editableVars {
		var valueStr string
		switch v := value.(type) {
		case *int:
			valueStr = fmt.Sprintf("%d", *v)
		case *string:
			valueStr = fmt.Sprintf("%s", *v)
		case *bool:
			valueStr = fmt.Sprintf("%t", *v)
		// добавьте обработку других типов указателей, если нужно
		default:
			valueStr = fmt.Sprintf("%v", value) // для случаев, когда тип неизвестен
		}

		saveToEnv := fmt.Sprintf("write user.env %s=%s", PasswordAlgoritm.Usage(key, true), PasswordAlgoritm.Usage(valueStr, true))
		execLtCommand(saveToEnv)
	}
}

func LoadUserConfigs() error {
	restoreOutput, err := silenceOutput() // Отключаем вывод
	if err != nil {
		fmt.Println(red("Error while disabling output:", err))
		return err
	}
	defer restoreOutput() // Восстанавливаем вывод в конце
	err = commands.ChangeDirectory(Absdir)
	if err != nil {
		fmt.Println(red(err))
		return err
	}

	userConfigs, err := ReadEnvUtil.File("user.env")
	if err != nil {
		fmt.Println(red(err))
		return err
	}

	userConfigsStr := fmt.Sprintf("%v", string(userConfigs))
	userConfigsStrList := strings.Split(userConfigsStr, "\n")
	for _, userConfigStr := range userConfigsStrList {
		setVar := strings.Split(userConfigStr, "=")
		if len(setVar) != 2 {
			continue
		}

		saveToEnv := fmt.Sprintf("setvar %s %s", PasswordAlgoritm.Usage(setVar[0], false), PasswordAlgoritm.Usage(setVar[1], false))
		execLtCommand(saveToEnv)
	}

	return nil
}
