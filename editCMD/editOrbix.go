package editCMD

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

func SayOrbix() {
	myFigure := figure.NewFigure("Orbix", "", true)
	magenta := color.New(color.FgMagenta).SprintFunc()
	fmt.Println(magenta(myFigure.String()))
}
