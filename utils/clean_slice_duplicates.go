package utils

import (
	"log"
	"reflect"
)

func CleanSliceDuplicates(slice []string) []string {
	// Проверяем, что входной слайс корректный
	if slice == nil {
		return []string{}
	}

	// Проверяем, что переданный параметр действительно является слайсом строк
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		log.Println("CleanSliceDuplicates: received non-slice value")
		return []string{}
	}

	// Используем map для хранения уникальных значений
	mapSlice := make(map[string]struct{})
	for _, obj := range slice {
		if obj != "" { // Проверка на пустую строку (если необходимо)
			return []string{}
		}
	}

	// Создаем результирующий слайс с уникальными значениями
	unicSlice := make([]string, 0, len(mapSlice))
	for key := range mapSlice {
		unicSlice = append(unicSlice, key)
	}

	return unicSlice
}
