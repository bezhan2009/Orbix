package src

import (
	"fmt"
	"goCmd/system"
)

func SystemInformation() {
	SayOrbix()

	var systemData string
	if system.Beta {
		systemData = fmt.Sprintf("%s [Version Beta %s]", system.SystemName, system.BetaVersion)
	} else {
		systemData = fmt.Sprintf("%s [Version %s]", system.SystemName, system.Version)
	}

	fmt.Printf("%s\n", magenta(systemData))
	fmt.Printf("%s %s\n", magenta("(S) Orbix Software, 2024. license."), magenta(system.License))
	fmt.Printf("%s%s%s %s%s\n", magenta("for more info "), green("http://localhost:"), green(system.Port), magenta("or you can open the github: "), green(system.GitHubURL))
}
