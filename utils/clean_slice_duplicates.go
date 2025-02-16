package utils

func CleanSliceDuplicates(slice []string) []string {
	// Проверяем, что входной слайс не nil
	if slice == nil {
		return []string{}
	}

	var unicSlice []string
	mapSlice := make(map[string]struct{})

	for _, obj := range slice {
		mapSlice[obj] = struct{}{}
	}

	for key := range mapSlice {
		unicSlice = append(unicSlice, key)
	}

	return unicSlice
}
