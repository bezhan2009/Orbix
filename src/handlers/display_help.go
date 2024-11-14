package handlers

import (
	"fmt"
	"goCmd/system"
)

func DisplayHelp() {
	fmt.Println(system.Yellow("\tFor command information, type HELP\n"))
	const nameWidth = 20 // Задаем ширину для названий команд

	for _, command := range system.Commands {
		if command.Description != "" {
			// Используем форматирование с указанием ширины поля для команд
			fmt.Println(system.Yellow(fmt.Sprintf("%-*s %s", nameWidth, command.Name, command.Description)))
		}
	}
}
