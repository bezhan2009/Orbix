package utils

import (
	"goCmd/cmd/commands/resourceIntensive/PiCalculation"
	"strconv"
)

func CalculatePiUtil(commandArgs []string) {
	precision := 10001
	if len(commandArgs) > 0 {
		if p, err := strconv.Atoi(commandArgs[0]); err == nil {
			precision = p
		}
	}

	PiCalculation.PiCalcCommand(precision)
}
