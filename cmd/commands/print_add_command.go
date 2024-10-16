package commands

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

func PrintAddCommand() string {
	myFigure := figure.NewFigure("Add Command!!!", "", true)
	greenText := color.New(color.FgGreen).SprintFunc()
	return greenText(myFigure.String())
}
