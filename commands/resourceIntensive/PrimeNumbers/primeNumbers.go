package PrimeNumbers

import (
	"fmt"
	"math"
	"time"
)

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func generatePrimes(limit int) []int {
	primes := []int{}
	for i := 2; i < limit; i++ {
		if isPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}

func PrimeCommand(limit int) {
	start := time.Now()
	primes := generatePrimes(limit)
	elapsed := time.Since(start)

	fmt.Printf("Generated %d primes up to %d in %s\n", len(primes), limit, elapsed)
}
