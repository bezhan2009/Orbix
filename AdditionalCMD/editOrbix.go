package AdditionalCMD

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"goCmd/system"
)

func SayOrbix() {
	// Creating a shape with the text "Orbix" in an accessible style
	say := fmt.Sprintf("%s %s", system.SystemName, system.Version)
	myFigure := figure.NewFigure(say, "larry3d", true)

	// Defining the color of the text as a magenta
	magenta := color.New(color.FgMagenta).SprintFunc()

	// We output a shape with text in magenta color
	fmt.Println(magenta(myFigure.String()))
}
