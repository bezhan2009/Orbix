package environment

import (
	"fmt"
	"goCmd/system"
	"goCmd/utils"
)

func SystemInformation() {
	utils.SayOrbix()

	var systemInfo string
	if system.Beta {
		systemInfo = fmt.Sprintf("%s [Version Beta %s]", system.SystemName, system.BetaVersion)
	} else {
		systemInfo = fmt.Sprintf("%s [Version %s]", system.SystemName, system.Version)
	}

	fmt.Printf("%s\n", system.Magenta(systemInfo))
	fmt.Printf("%s %s\n", system.Magenta("(S) CMD Software, 2024. license."), system.Magenta(system.License))
	fmt.Printf("%s%s%s %s%s\n", system.Magenta("for more info "), system.Green("http://localhost:"), system.Green(system.Port), system.Magenta("or you can open the github: "), system.Green(system.GitHubURL))
}
