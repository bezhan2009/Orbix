package environment

import (
	"fmt"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/system/errs"
	utils2 "goCmd/utils"
	utils3 "goCmd/validators/utils"
	"strings"
)

func SetShortCutUtil(args []string) {
	colors := make(map[string]func(...interface{}) string)
	colors = system.GetColorsMap()

	if len(args) < 2 {
		fmt.Println(args)
		fmt.Println(colors["yellow"]("Usage: shortcut <shortcut_name> <value>"))
		return
	}

	var (
		varName string
		value   string
	)

	for iArg, arg := range args {
		if iArg == 0 {
			varName = args[0]
			continue
		}

		value += arg + " "
	}

	value = strings.TrimSpace(value)

	err := SetShortcut(strings.ToLower(strings.TrimSpace(varName)), value)
	if err != nil {
		fmt.Printf(colors["red"](fmt.Sprintf("Error: %s\n", err.Error())))
	} else {
		fmt.Printf(colors["green"](fmt.Sprintf("the values of the shortcut %s have been changed to %s successfully\n", varName, value)))
	}
}

// SetShortcut изменяет значение переменной по её имени с преобразованием типов
func SetShortcut(varName string, value string) error {
	if utils2.ValidCommandFast(varName, utils3.ValidateSymbols) {
		return errs.ValidationError
	}

	system.Shortcuts[varName] = value
	return nil
}

func DeleteShortcut(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println(system.Yellow("Usage: delshort <shortcut_name>"))
		return
	}

	shortcutName := commandArgs[0]
	if strings.TrimSpace(shortcutName) == "*" {
		system.Shortcuts = make(map[string]string)
		return
	}

	if _, exists := system.Shortcuts[strings.ToLower(strings.TrimSpace(shortcutName))]; !exists {
		fmt.Println(system.Red(fmt.Sprintf("the shortcut %s is invalid\n", shortcutName)))
		return
	}

	delete(system.Shortcuts, shortcutName)
}

func GetShortcutValueUtil(params *structs.ExecuteCommandFuncParams) string {
	args := params.CommandArgs

	if len(args) < 1 {
		fmt.Println(system.Yellow("Usage: getshort <shortcut_name>"))
		fmt.Println(system.Yellow("Or: getshort *"))
		return ""
	}

	shortName := args[0]

	if strings.TrimSpace(shortName) == "*" {
		for k, v := range system.Shortcuts {
			fmt.Println(system.GreenBold(fmt.Sprintf("%s: %s", k, v)))
		}
		return ""
	}

	fmt.Println(system.GreenBold(fmt.Sprintf("%s: %s", shortName, system.Shortcuts[shortName])))
	return system.Shortcuts[shortName]
}
