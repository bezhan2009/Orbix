package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/AdditionalCMD"
	"goCmd/system"
)

func SystemInformation() {
	AdditionalCMD.SayOrbix()
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	systemData := fmt.Sprintf("%s [Version %s]", system.SystemName, system.Version)
	fmt.Printf("%s\n", magenta(systemData))
	fmt.Printf("%s %s\n", magenta("(S) Orbix Software, 2024. license."), magenta(system.License))
	fmt.Printf("%s:%s %s\n", magenta("for more info http://localhost"), system.Port, magenta("or you can open the github: https://github.com/bezhan2009/Orbix"))
}
