package environment

import (
	"encoding/json"
	"fmt"
	"goCmd/cmd/commands"
	"goCmd/cmd/commands/Remove"
	"goCmd/cmd/commands/fcommands"
	"goCmd/pkg/algorithms/PasswordAlgoritm"
	"goCmd/system"
	"goCmd/utils"
	"io/ioutil"
	"os"
	"reflect"
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
		err := nullFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}, nil
}

// Вспомогательная функция для получения значений из указателей
func dereferenceVariables(vars map[string]interface{}) map[string]interface{} {
	values := make(map[string]interface{})
	for key, ptr := range vars {
		val := reflect.ValueOf(ptr)
		if val.Kind() == reflect.Ptr && !val.IsNil() {
			values[key] = val.Elem().Interface()
		} else {
			values[key] = nil // Если указатель нулевой, значение сохраняется как nil
		}
	}
	return values
}

// Вспомогательная функция для обновления значений по указателям
func updatePointers(vars map[string]interface{}, data map[string]interface{}) {
	for key, ptr := range vars {
		if newValue, exists := data[key]; exists {
			val := reflect.ValueOf(ptr)
			if val.Kind() == reflect.Ptr && val.Elem().CanSet() {
				val.Elem().Set(reflect.ValueOf(newValue))
			}
		}
	}
}

// SaveVars Сохранение переменных в формате JSON
func SaveVars() {
	restoreOutput, err := silenceOutput() // Отключаем вывод
	if err != nil {
		fmt.Println("Error while disabling output:", err)
		return
	}
	defer restoreOutput() // Восстанавливаем вывод в конце

	system.UserDir, _ = os.Getwd()
	system.EditableVars["user_dir"] = &system.UserDir
	err = commands.ChangeDirectory(system.Absdir)
	if err != nil {
		fmt.Println(system.Red(err))
		return
	}
	defer func() {
		err = commands.ChangeDirectory(system.UserDir)
		if err != nil {
			fmt.Println(system.Red(err))
		}
	}()

	_, err = Remove.File("rem", []string{"user.json"})
	if err != nil {
		fmt.Println("Error while removing file user.json:", err)
	}
	_, err = fcommands.CreateFile("user.json")
	if err != nil {
		fmt.Println("Error while creating file user.json:", err)
	}

	values := dereferenceVariables(system.EditableVars)

	for key, value := range values {
		valueStr := fmt.Sprintf("%v", value)
		values[key] = PasswordAlgoritm.Usage(valueStr, true)
	}

	data, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		fmt.Println("Error serializing variables:", err)
		return
	}

	err = ioutil.WriteFile("user.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}

// LoadUserConfigs Загрузка переменных из JSON и обновление указателей
func LoadUserConfigs() error {
	restoreOutput, err := silenceOutput() // Отключаем вывод
	if err != nil {
		fmt.Println(system.Red("Error while disabling output:", err))
		return err
	}
	defer restoreOutput() // Восстанавливаем вывод в конце

	err = commands.ChangeDirectory(system.Absdir)
	if err != nil {
		fmt.Println(system.Red(err))
		return err
	}
	defer func() {
		err = commands.ChangeDirectory(system.UserDir)
		if err != nil {
			fmt.Println(system.Red(err))
		}
	}()

	file, err := os.Open("user.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	loadedValues := map[string]interface{}{
		"location": &system.Location,
		"prompt":   &system.Prompt,
		"user":     &system.User,
		"empty":    &system.Empty,
	}

	err = json.Unmarshal(data, &loadedValues)
	if err != nil {
		fmt.Println("Error deserializing JSON:", err)
		return err
	}

	updatePointers(loadedValues, system.EditableVars)

	// Установка переменных в окружение
	for key, value := range loadedValues {
		valueStr := fmt.Sprintf("%v", value)
		saveToEnv := fmt.Sprintf("%s %v", PasswordAlgoritm.Usage(key, false), PasswordAlgoritm.Usage(valueStr, false))
		SetVariableUtil(utils.SplitCommandLine(saveToEnv))
	}

	err = commands.ChangeDirectory(system.UserDir)
	if err != nil {
		return err
	}

	return nil
}
