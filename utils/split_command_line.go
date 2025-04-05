package utils

import "strings"

func SplitCommandLine(input string) []string {
	var result []string
	var current strings.Builder
	inDoubleQuotes := false
	inSingleQuotes := false

	for _, r := range input {
		switch {
		case r == ' ' && !inDoubleQuotes && !inSingleQuotes:
			// Разделяем по пробелу, если не внутри кавычек
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		case r == '=' && !inDoubleQuotes && !inSingleQuotes:
			// Разделяем по пробелу, если не внутри кавычек
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		case r == '"':
			// Если встречена двойная кавычка, переключаем состояние
			if !inSingleQuotes {
				inDoubleQuotes = !inDoubleQuotes
			} else {
				current.WriteRune(r) // Если внутри одинарных кавычек, добавляем символ
			}
		case r == '(':
			// Если встречена скобка, переключаем состояние
			if !inSingleQuotes {
				inDoubleQuotes = !inDoubleQuotes
			} else {
				current.WriteRune(r) // Если внутри одинарных кавычек, добавляем символ
			}
		case r == '\'':
			// Если встречена вот такая палка \, переключаем состояние
			if !inDoubleQuotes {
				inSingleQuotes = !inSingleQuotes
			} else {
				current.WriteRune(r) // Если внутри двойных кавычек, добавляем символ
			}
		default:
			// Добавляем символ в текущую строку
			current.WriteRune(r)
		}
	}

	// Добавляем последний накопленный элемент
	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
