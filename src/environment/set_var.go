package environment

import (
	"errors"
	"fmt"
	"goCmd/src/utils"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/system/errs"
	utils2 "goCmd/utils"
	utils3 "goCmd/validators/utils"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrNotFoundAndCreated = fmt.Errorf("variable not found, created new one with this name")
)

func SetVariableUtil(args []string) {
	colors := make(map[string]func(...interface{}) string)
	colors = system.GetColorsMap()

	if len(args) < 2 {
		fmt.Println(args)
		fmt.Println(colors["yellow"]("Usage: setvar <variable_name> <value>"))
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

	err := SetVariable(strings.ToLower(strings.TrimSpace(varName)), value)
	if err != nil {
		fmt.Printf(colors["red"](fmt.Sprintf("Error: %s\n", err.Error())))
	} else {
		fmt.Printf(colors["green"](fmt.Sprintf("the values of the variable %s have been changed to %s successfully\n", varName, value)))
	}
}

// SetVariable изменяет значение переменной по её имени с преобразованием типов
func SetVariable(varName string, value string) error {
	if utils2.ValidCommandFast(varName, utils3.ValidateSymbols) {
		return errs.ValidationError
	}

	if strings.TrimSpace(strings.ToLower(varName)) == "empty" {
		return fmt.Errorf(fmt.Sprintf("the variable %s is nil\n", varName))
	}

	if strings.TrimSpace(value) == "current_user" || strings.TrimSpace(value) == "$current_user" {
		value = system.User
	}

	// Проверяем, есть ли такая переменная в нашем списке
	if variable, exists := system.EditableVars[varName]; exists {
		v := reflect.ValueOf(variable).Elem()
		newValue := reflect.ValueOf(value)

		// Преобразование типов в зависимости от типа переменной
		switch v.Kind() {
		case reflect.String:
			// Преобразуем в строку, если это возможно
			if newValue.Kind() != reflect.String {
				convertedValue, err := utils.ConvertToString(newValue)
				if err != nil {
					return err
				}
				v.SetString(convertedValue)
			} else {
				v.SetString(newValue.String())
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// Преобразуем в int
			convertedValue, err := utils.ConvertToInt(newValue)
			if err != nil {
				return err
			}
			v.SetInt(convertedValue)

		case reflect.Bool:
			// Преобразуем в bool
			convertedValue, err := utils.ConvertToBool(newValue)
			if err != nil {
				return err
			}
			v.SetBool(convertedValue)

		default:
			return fmt.Errorf("the %s variable type is not supported", varName)
		}
		return nil
	}

	if strings.TrimSpace(varName) == "user" {
		system.User = value
	}

	if strings.TrimSpace(varName) == "*" {
		return fmt.Errorf(fmt.Sprintf("the variable %s is nil\n", varName))
	}

	_, err := strconv.Atoi(string(varName[0]))
	if err == nil {
		return errors.New(fmt.Sprintf("Variable cannot starts with number: %s", varName))
	}

	// Если переменная не найдена, добавляем её в список с переданным значением
	system.AvailableEditableVars = append(system.AvailableEditableVars, varName)
	system.CustomEditableVars = append(system.CustomEditableVars, varName)
	system.EditableVars[varName] = &value
	return ErrNotFoundAndCreated
}

func DeleteVariable(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println(system.Yellow("Usage: delvar <variable_name>"))
		return
	}

	varname := commandArgs[0]
	if strings.TrimSpace(varname) == "*" {
		for i := 0; i < len(system.CustomEditableVars); i++ {
			delete(system.EditableVars, system.CustomEditableVars[i])
		}

		return
	}

	if !utils2.IsValid(varname, system.CustomEditableVars) {
		fmt.Println(system.Red(fmt.Sprintf("the variable %s is invalid\n", varname)))
		return
	}

	delete(system.EditableVars, varname)
}

func GetVariableValueUtil(params *structs.ExecuteCommandFuncParams) {
	args := params.CommandArgs

	if len(args) < 1 {
		fmt.Println(system.Yellow("Usage: getvar <variable_name>"))
		fmt.Println(system.Yellow("Or: getvar *"))
		return
	}

	varName := args[0]

	if strings.TrimSpace(varName) == "user" {
		if strings.TrimSpace(system.User) == "" {
			fmt.Println(system.Green("user:"), system.Green(params.LoopData.Username))
			return
		} else {
			fmt.Println(system.Green("user:"), system.User)
			return
		}
	}

	if strings.TrimSpace(varName) == "current_user" {
		fmt.Println(system.Green("current_user:"), system.Green(params.LoopData.Username))
		return
	}

	if strings.TrimSpace(varName) == "*" {
		fmt.Println(system.Green("current_user:"), system.Green(params.LoopData.Username))

		for _, v := range system.AvailableEditableVars {
			value, err := GetVariableValue(v)
			if err != nil {
				continue
			}

			fmt.Println(system.Green(fmt.Sprintf("%s: %s", v, value)))
		}

		return
	}

	value, err := GetVariableValue(varName)
	if err != nil {
		fmt.Println(system.Red(err.Error()))
		return
	}

	fmt.Println(system.Green(fmt.Sprintf("%s: %s", varName, value)))
}

func GetVariableValue(varName string) (interface{}, error) {
	variable, exists := system.EditableVars[strings.TrimSpace(strings.ToLower(varName))]
	if !exists {
		return nil, errors.New(fmt.Sprintf("The %s variable was not found or cannot be changed", varName))
	}

	// Разыменование указателя на значение
	switch v := variable.(type) {
	case *int:
		return *v, nil
	case *[]int:
		return *v, nil
	case *string:
		return *v, nil
	case *[]string:
		return *v, nil
	case *bool:
		return *v, nil
	case *[]bool:
		return *v, nil
	case *float64:
		return *v, nil
	case *[]float64:
		return *v, nil
	default:
		return v, errors.New(fmt.Sprintf("Unsupported variable type for %s", varName))
	}
}
