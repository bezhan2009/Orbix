package PiCalculation

import (
	"fmt"
	"math/big"
	"time"
)

func calculatePi(precision int) *big.Float {
	pi := new(big.Float).SetPrec(uint(precision)).SetFloat64(0)
	k := new(big.Float).SetPrec(uint(precision))
	n := new(big.Float).SetPrec(uint(precision))
	one := new(big.Float).SetPrec(uint(precision)).SetInt64(1)
	four := new(big.Float).SetPrec(uint(precision)).SetInt64(4)

	for i := int64(0); i < int64(precision); i++ {
		k.SetInt64(2 * i)
		n.Quo(one, new(big.Float).SetPrec(uint(precision)).SetInt64(2*i+1))
		if i%2 == 0 {
			pi.Add(pi, n)
		} else {
			pi.Sub(pi, n)
		}
	}

	pi.Mul(pi, four)
	return pi
}

func PiCalcCommand(precision int) {
	start := time.Now()
	pi := calculatePi(precision)
	elapsed := time.Since(start)

	fmt.Printf("Calculated pi to %d digits in %s\n", precision, elapsed)
	fmt.Println("Pi:", pi.Text('f', 6))
}
