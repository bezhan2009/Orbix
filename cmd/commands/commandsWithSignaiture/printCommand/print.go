package printCommand

import (
	"fmt"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

// Print Функция для печати текста с поддержкой шрифтов и цветов
func Print(commandArgs []string) {
	var (
		font       string
		colorFuncs = map[string]func(a ...interface{}) string{
			"red":     color.New(color.FgRed).SprintFunc(),
			"yellow":  color.New(color.FgYellow).SprintFunc(),
			"cyan":    color.New(color.FgCyan).SprintFunc(),
			"green":   color.New(color.FgGreen).SprintFunc(),
			"magenta": color.New(color.FgMagenta).SprintFunc(),
			"blue":    color.New(color.FgBlue).SprintFunc(),
		}
	)

	// Поиск параметра font
	for i, arg := range commandArgs {
		if strings.HasPrefix(arg, "font=") {
			font = strings.Split(arg, "=")[1]
			commandArgs = append(commandArgs[:i], commandArgs[i+1:]...)
		}
	}

	// Вывод текста в зависимости от шрифта и цвета
	for _, arg := range commandArgs {
		parts := strings.Split(arg, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.Contains(part, ":") {
				colorText := strings.SplitN(part, ":", 2)
				colorName := strings.TrimSpace(colorText[0])
				text := strings.TrimSpace(colorText[1])

				if colorFunc, ok := colorFuncs[colorName]; ok {
					printWithFont(text, font, colorFunc)
				} else {
					// Если цвет не найден, вывести обычный текст
					printWithFont(text, font, fmt.Sprint)
				}
			} else {
				// Вывести текст без цвета
				text := part
				printWithFont(text, font, fmt.Sprint)
			}
		}
	}

	fmt.Println()
}

// Вспомогательная функция для вывода текста со шрифтом и цветом
func printWithFont(text, font string, colorFunc func(a ...interface{}) string) {
	if font == "3d" {
		myFigure := figure.NewFigure(text, "larry3d", true)
		fmt.Println(colorFunc(myFigure.String())) // Выводим текст в 3D с цветом
	} else if font == "2d" {
		myFigure := figure.NewFigure(text, "", true)
		fmt.Println(colorFunc(myFigure.String())) // Выводим текст в 2D с цветом
	} else {
		// Обычный вывод без 3D/2D эффекта
		fmt.Print(colorFunc(text), " ")
	}
}
