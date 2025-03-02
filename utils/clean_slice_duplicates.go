package utils

func CleanSliceDuplicates(slice []string) []string {
	// Проверяем, что входной слайс не nil
	if slice == nil {
		return []string{}
	}

	// Используем map для хранения уникальных значений
	mapSlice := make(map[string]struct{})
	for _, obj := range slice {
		if obj != "" { // Проверка на пустую строку (если необходимо)
			mapSlice[obj] = struct{}{}
		}
	}

	// Создаем результирующий слайс с уникальными значениями
	unicSlice := make([]string, 0, len(mapSlice))
	for key := range mapSlice {
		unicSlice = append(unicSlice, key)
	}

	return unicSlice
}
