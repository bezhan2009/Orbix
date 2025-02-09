package Orbix

import (
	"fmt"
	"goCmd/system"
	"path/filepath"
	"strings"
)

func TemplateUtil(commandArgs []string, SD *system.AppState) {
	// Проверяем, передано ли имя шаблона
	if len(commandArgs) < 1 {
		fmt.Println(system.Yellow("Usage: template <template_name> [echo=on|echo=off]"))
		return
	}

	// Если флаг echo не передан, добавляем значение по умолчанию
	if len(commandArgs) < 2 {
		commandArgs = append(commandArgs, "echo=on")
	}

	// Приводим флаг к булевому значению
	switch commandArgs[1] {
	case "echo=on":
		commandArgs[1] = "true"
	case "echo=off":
		commandArgs[1] = "false"
	default:
		commandArgs[1] = "true" // значение по умолчанию, если аргумент не распознан
	}

	// Проверяем расширение файла шаблона
	extension := strings.ToLower(filepath.Ext(commandArgs[0]))
	// Если расширение пустое или его длина меньше 2 (например, "" или ".")
	if len(extension) < 2 || extension[1:] != system.OrbixTemplatesExtension {
		fmt.Println(system.Red(fmt.Sprintf("The template extension must be .%s", system.OrbixTemplatesExtension)))
		return
	}

	// Если всё в порядке, запускаем обработку шаблона
	if err := Start(commandArgs[0], commandArgs[1], SD); err != nil {
		fmt.Println(system.Red(err))
	}
}
