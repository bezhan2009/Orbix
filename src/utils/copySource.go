package utils

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/copySource"
)

func CommandCopySourceUtil(commandArgs []string) {
	if len(commandArgs) < 2 {
		fmt.Println("Usage: copysource <srcDir> <dstDir>")
		fmt.Println("Example: copysource example.txt destination_directory.txt:")
		fmt.Println("copysource example.txt bufer")
		fmt.Println("Arguments:")
		fmt.Println("copysource example.txt bufer")
		return
	}
	copySource.File(commandArgs[0], commandArgs[1])
}
