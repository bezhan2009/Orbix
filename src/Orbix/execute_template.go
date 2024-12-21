package Orbix

import (
	"fmt"
	"goCmd/system"
	"path/filepath"
	"strings"
)

func TemplateUtil(commandArgs []string, SD *system.AppState) {
	if len(commandArgs) < 1 {
		fmt.Println(system.Yellow("Usage: template <template_name> echo=on"))
		fmt.Println(system.Yellow("Or: template <template_name> echo=off if you want without outputting the result"))
		return
	}

	if commandArgs[1] != "echo=on" && commandArgs[1] != "echo=off" {
		commandArgs[1] = "true"
	} else if commandArgs[1] == "echo=on" {
		commandArgs[1] = "true"
	} else if commandArgs[1] == "echo=off" {
		commandArgs[1] = "false"
	}

	extension := strings.ToLower(filepath.Ext(commandArgs[0]))
	if extension[1:] != system.OrbixTemplatesExtension {
		fmt.Println(system.Red(fmt.Sprintf("The template extension must be %s", system.OrbixTemplatesExtension)))
		return
	}

	if err := Start(commandArgs[0], commandArgs[1], SD); err != nil {
		fmt.Println(system.Red(err))
	}
}
