package src

import (
	"fmt"
)

func displayHelp() {
	fmt.Println(yellow("\tFor command information, type HELP\n"))
	const nameWidth = 20 // Задаем ширину для названий команд

	for _, command := range Commands {
		if command.Description != "" {
			// Используем форматирование с указанием ширины поля для команд
			fmt.Println(yellow(fmt.Sprintf("%-*s %s", nameWidth, command.Name, command.Description)))
		}
	}
}
