package commands

import (
	"fmt"
	"goCmd/system"
	"goCmd/utils"
	"regexp"
	"strings"

	"github.com/common-nighthawk/go-figure"
)

// Print Функция для печати текста с поддержкой шрифтов и цветов
func Print(commandArgs []string) {
	var (
		font          string
		animatedPrint bool
		colorFuncs    = system.GetColorsMap()
	)

	// Поиск параметра font
	for i, arg := range commandArgs {
		if strings.HasPrefix(arg, "font=") {
			font = strings.Split(arg, "=")[1]
			commandArgs = append(commandArgs[:i], commandArgs[i+1:]...)
		}

		if strings.HasPrefix(arg, "animate") {
			animatedPrint = true
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
				if len(colorText) == 2 {
					colorName := strings.TrimSpace(colorText[0])
					text := strings.TrimSpace(colorText[1])

					if colorFunc, ok := colorFuncs[colorName]; ok {
						printWithFont(text, font, colorFunc, animatedPrint)
					} else {
						printWithFont(text, font, fmt.Sprint, animatedPrint)
					}
				} else {
					fmt.Println("Ошибка: Неправильный формат для цвета и текста.")
				}
			} else {
				// Вывести текст без цвета
				text := part
				printWithFont(text, font, fmt.Sprint, animatedPrint)
			}
		}
	}

	fmt.Println()
}

// Вспомогательная функция для вывода текста со шрифтом и цветом
func printWithFont(text, font string, colorFunc func(a ...interface{}) string, animate bool) {
	// Проверка текста, если используется 2D или 3D шрифт
	if font == "2d" || font == "3d" {
		// Регулярное выражение для фильтрации: только английские буквы и цифры
		re := regexp.MustCompile(`[^a-zA-Z0-9 !@#+$%^&*()_]`)
		text = re.ReplaceAllString(text, "")

		// Если текст пустой после фильтрации
		if text == "" {
			fmt.Println("Недопустимые символы для выбранного шрифта.")
			return
		}
	}

	if font == "3d" {
		myFigure := figure.NewFigure(text, "larry3d", true)
		if !animate {
			fmt.Println(colorFunc(myFigure.String()))
		} else {
			utils.PrintAnim(colorFunc(myFigure.String()))
		} // Выводим текст в 3D с цветом
	} else if font == "2d" {
		myFigure := figure.NewFigure(text, "", true)
		if !animate {
			fmt.Println(colorFunc(myFigure.String()))
		} else {
			utils.PrintAnim(colorFunc(myFigure.String()))
		} // Выводим текст в 2D с цветом
	} else {
		// Обычный вывод без 3D/2D эффекта
		if !animate {
			fmt.Print(colorFunc(text), " ")
		} else {
			utils.PrintAnim(colorFunc(text))
		}
	}
}
