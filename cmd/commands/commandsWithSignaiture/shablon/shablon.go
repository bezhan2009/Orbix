package shablon

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/Edit"
	"os"
)

func Make() {
	name := ""

	fmt.Println("Названия шаблона:")
	fmt.Scan(&name)

	os.Create(name)
	Edit.File(name)
}
