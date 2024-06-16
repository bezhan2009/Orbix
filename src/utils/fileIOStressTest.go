package utils

import (
	"goCmd/commands/resourceIntensive/FileIOStressTest"
	"strconv"
)

func FileIOStressTestUtil(commandArgs []string) {
	filename := "largefile.dat"
	size := 100 * 1024 * 1024
	if len(commandArgs) > 0 {
		if s, err := strconv.Atoi(commandArgs[0]); err == nil {
			size = s
		}
	}
	FileIOStressTest.FileIOCommand(filename, size)
}
