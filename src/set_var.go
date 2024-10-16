package src

import (
	"errors"
	"fmt"
	"goCmd/src/utils"
	"goCmd/structs"
	"goCmd/system"
	"reflect"
	"strings"
)

func SetVariableUtil(args []string) {
	colors := make(map[string]func(...interface{}) string)
	colors = system.GetColorsMap()

	if len(args) < 2 {
		fmt.Println(colors["yellow"]("Usage: setvar <var name> <value>"))
		return
	}

	varName := args[0]
	value := args[1]
	err := SetVariable(strings.ToLower(strings.TrimSpace(varName)), value)
	if err != nil {
		fmt.Printf(colors["red"](fmt.Sprintf("Error: %s\n", err.Error())))
	} else {
		fmt.Printf(colors["green"](fmt.Sprintf("the values of the variable %s have been changed to %s successfully\n", varName, value)))
	}
}

// SetVariable изменяет значение переменной по её имени с преобразованием типов
func SetVariable(varName string, value interface{}) error {
	// Проверяем, есть ли такая переменная в нашем списке
	if variable, exists := editableVars[varName]; exists {
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
	return errors.New(fmt.Sprintf("The %s variable was not found or cannot be changed", varName))
}

func GetVariableValueUtil(params structs.ExecuteCommandFuncParams) {
	args := params.CommandArgs

	if len(args) < 1 {
		fmt.Println(yellow("Usage: getvar <variable_name>"))
		fmt.Println(yellow("Or: getvar *"))
		return
	}

	varName := args[0]

	if strings.TrimSpace(varName) == "user" {
		if strings.TrimSpace(User) == "" {
			fmt.Println(green("user:"), green(params.Username))
			return
		} else {
			fmt.Println(green("user:"), User)
			return
		}
	}

	if strings.TrimSpace(varName) == "current_user" {
		fmt.Println(green("current_user:"), green(params.Username))
		return
	}

	if strings.TrimSpace(varName) == "*" {
		fmt.Println(green("current_user:"), green(params.Username))

		for _, v := range availableEditableVars {
			value, err := GetVariableValue(v)
			if err != nil {
				fmt.Println(red(err.Error()))
			}

			fmt.Println(green(fmt.Sprintf("%s: %s", v, value)))
		}

		return
	}

	value, err := GetVariableValue(varName)
	if err != nil {
		fmt.Println(red(err.Error()))
	}

	fmt.Println(green(fmt.Sprintf("%s: %s", varName, value)))
}

func GetVariableValue(varName string) (interface{}, error) {
	variable, exists := editableVars[strings.TrimSpace(strings.ToLower(varName))]
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
		return nil, errors.New(fmt.Sprintf("Unsupported variable type for %s", varName))
	}
}
