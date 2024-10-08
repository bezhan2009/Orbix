package src

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"goCmd/system"
)

func SayOrbix() {
	// Creating a shape with the text "Orbix" in an accessible style
	var say string
	if system.BetaVersion {
		say = fmt.Sprintf("%s Beta", system.SystemName)
	} else {
		say = fmt.Sprintf("%s %s", system.SystemName, system.Version)
	}

	myFigure := figure.NewFigure(say, "larry3d", true)

	// We output a shape with text in magenta color
	fmt.Println(magenta(myFigure.String()))
}
