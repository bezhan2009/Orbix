package src

import (
	"encoding/json"
	"fmt"
	"goCmd/cmd/commands"
	"goCmd/pkg/algorithms/PasswordAlgoritm"
	"goCmd/system"
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
		nullFile.Close()
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

	execLtCommand("delete user.json")
	execLtCommand("create user.json")

	err = commands.ChangeDirectory(Absdir)
	if err != nil {
		fmt.Println(red(err))
		return
	}

	values := dereferenceVariables(editableVars)

	for key, value := range values {
		valueStr := fmt.Sprintf("%v", value)
		values[key] = PasswordAlgoritm.Usage(valueStr, true)
	}

	data, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		fmt.Println("Ошибка при сериализации переменных:", err)
		return
	}

	err = ioutil.WriteFile("user.json", data, 0644)
	if err != nil {
		fmt.Println("Ошибка при записи файла:", err)
	}
}

// LoadUserConfigs Загрузка переменных из JSON и обновление указателей
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

	file, err := os.Open("user.json")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return err
	}

	loadedValues := map[string]interface{}{
		"location": &Location,
		"prompt":   &Prompt,
		"user":     &User,
		"empty":    &Empty,
	}

	err = json.Unmarshal(data, &loadedValues)
	if err != nil {
		fmt.Println("Ошибка при десериализации JSON:", err)
		return err
	}

	updatePointers(loadedValues, editableVars)

	// Установка переменных в окружение
	for key, value := range loadedValues {
		valueStr := fmt.Sprintf("%v", value)
		saveToEnv := fmt.Sprintf("setvar %s %v", PasswordAlgoritm.Usage(key, false), PasswordAlgoritm.Usage(valueStr, false))
		execLtCommand(saveToEnv)
	}

	return nil
}
