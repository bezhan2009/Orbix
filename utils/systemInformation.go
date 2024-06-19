package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/editCMD"
)

func SystemInformation() {
	editCMD.SayOrbix()
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	fmt.Printf("%s\n", magenta("Orbix [Версия 0.94]"))
	fmt.Printf("%s\n", magenta("(S) Orbix Software, 2024. No license yet."))
}
