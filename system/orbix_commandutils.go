package system

import "sort"

// Инициализация: сортируем список по алфавиту по полю Name
func init() {
	sort.Slice(AdditionalCommands, func(i, j int) bool {
		return AdditionalCommands[i].Name < AdditionalCommands[j].Name
	})

	sort.Slice(Commands, func(i, j int) bool {
		return AdditionalCommands[i].Name < AdditionalCommands[j].Name
	})
}
