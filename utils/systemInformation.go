package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/editCMD"
)

func SystemInformation() {
	editCMD.StartEditing()
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	fmt.Printf("%s\n", magenta("src [Версия 0.94]"))
	fmt.Printf("%s\n", magenta("(c) src Software, 2024. Все права защищены."))
}
