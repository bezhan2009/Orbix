package utils

import (
	"goCmd/commands/resourceIntensive/PrimeNumbers"
	"strconv"
)

func CalculatePrimesUtil(commandArgs []string) {
	limit := 100000
	if len(commandArgs) > 0 {
		if l, err := strconv.Atoi(commandArgs[0]); err == nil {
			limit = l
		}
	}
	PrimeNumbers.PrimeCommand(limit)
}
