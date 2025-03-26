package src

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	_chan "goCmd/chan"
	"goCmd/src/service"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
	"strings"
)

func commandVar(commandLower string) bool {
	// Если набор команд будет расширяться, можно использовать map для O(1) поиска.
	return commandLower == "setvar" ||
		commandLower == "delvar" ||
		commandLower == "getvar" ||
		commandLower == "cd" ||
		commandLower == "cd.."
}

func ReadCommandLine(commandInput string) (string, string, []string, string) {
	var commandLine string
	if commandInput != "" {
		commandLine = strings.TrimSpace(commandInput)
	} else {
		// Чтение ввода
		commandLine = strings.TrimSpace(prompt.Input("", service.AutoComplete))
	}

	// Если ввода нет, сразу возвращаем пустые значения
	if commandLine == "" {
		return "", "", nil, ""
	}

	// Используем Fields для разбиения по любым пробельным символам
	tokens := strings.Fields(commandLine)
	if len(tokens) == 0 {
		return "", "", nil, ""
	}

	// Вычисляем нижний регистр первого токена один раз
	firstTokenLower := strings.ToLower(tokens[0])
	if commandVar(firstTokenLower) {
		_chan.IsVarsFnUpd = true
		// Возвращаем: исходную строку, первый токен, оставшиеся токены и первый токен в нижнем регистре
		return commandLine, tokens[0], tokens[1:], firstTokenLower
	}

	// Если команда не является setvar/delvar/getvar, используем разбиение по более сложной логике
	commandParts := utils.SplitCommandLine(commandLine)
	if len(commandParts) == 0 {
		return "", "", nil, ""
	}

	// Вычисляем нижний регистр первого элемента и возвращаем результат
	firstPartLower := strings.ToLower(commandParts[0])
	return commandLine, commandParts[0], commandParts[1:], firstPartLower
}

func ProcessCommand(commandLower string) bool {
	if strings.TrimSpace(commandLower) == "cd" && system.GitExists {
		return true
	}

	if strings.TrimSpace(commandLower) == "git" && system.GitExists {
		return true
	}

	return false
}

func CatchSyntaxErrs(execCommandCatchErrs structs.ExecuteCommandCatchErrs) (findErr bool) {
	if *execCommandCatchErrs.EchoTime && *execCommandCatchErrs.RunOnNewThread && !(execCommandCatchErrs.CommandLower == "orbix") {
		fmt.Println(system.Red("You cannot take timing and running on new thread at the same time"))
		return true
	}

	return false
}

// RemoveFlags удаляет части строки, если они содержатся в OrbixFlags
func RemoveFlags(input string) string {
	// Разделяем строку на части
	parts := strings.Fields(input)
	var result []string

	// Проходим по всем частям
	for _, part := range parts {
		// Проверяем, есть ли текущая часть в OrbixFlags
		if !contains(system.Flags, part) {
			// Если часть не является флагом, добавляем её в результат
			result = append(result, part)
		}
	}

	// Соединяем оставшиеся части в строку и возвращаем
	return strings.Join(result, " ")
}

// contains проверяет, находится ли элемент в списке
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func CommandFile(command string) bool {
	return command == "py" ||
		command == "read" ||
		command == "edit" ||
		command == "create" ||
		command == "rem" ||
		command == "delvar" ||
		command == "format" ||
		command == "del" ||
		command == "gocode" ||
		command == "delete" ||
		command == "cf" ||
		command == "df" ||
		command == "rustc" ||
		command == "cl"
}

func FullFileName(commandArgs *[]string) {
	if len(*commandArgs) <= 1 {
		return
	}

	var fileName string

	if len(*commandArgs) > 1 {
		for _, arg := range *commandArgs {
			fileName += arg + " "
		}

		fileName = strings.TrimSpace(fileName)

		resultSlice := []string{fileName}
		*commandArgs = resultSlice
	} else {
		return
	}
}
