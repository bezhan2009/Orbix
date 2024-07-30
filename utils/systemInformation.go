package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/editCMD"
	"goCmd/system"
)

func SystemInformation() {
	editCMD.SayOrbix()
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	systemData := fmt.Sprintf("%s [Version %s]", system.SystemName, system.Version)
	fmt.Printf("%s\n", magenta(systemData))
	fmt.Printf("%s %s\n", magenta("(S) Orbix Software, 2024. license"), magenta(system.License))
}
