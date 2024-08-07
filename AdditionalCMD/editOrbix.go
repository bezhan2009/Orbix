package AdditionalCMD

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"goCmd/system"
)

func SayOrbix() {
	// Создаем фигуру с текстом "Orbix" в доступном стиле
	say := fmt.Sprintf("%s %s", system.SystemName, system.Version)
	myFigure := figure.NewFigure(say, "larry3d", true)

	// Определяем цвет текста как магента
	magenta := color.New(color.FgMagenta).SprintFunc()

	// Выводим фигуру с текстом в магента цвете
	fmt.Println(magenta(myFigure.String()))
}
