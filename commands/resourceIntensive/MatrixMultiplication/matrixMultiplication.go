package MatrixMultiplication

import (
	"fmt"
	"math/rand"
	"time"
)

func multiplyMatrices(a, b [][]int) [][]int {
	n := len(a)
	c := make([][]int, n)
	for i := range c {
		c[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return c
}

func generateMatrix(size int) [][]int {
	matrix := make([][]int, size)
	for i := range matrix {
		matrix[i] = make([]int, size)
		for j := range matrix[i] {
			matrix[i][j] = rand.Intn(10)
		}
	}
	return matrix
}

func MatrixMulCommand() {
	size := 500
	a := generateMatrix(size)
	b := generateMatrix(size)

	start := time.Now()
	c := multiplyMatrices(a, b)
	elapsed := time.Since(start)

	fmt.Printf("Matrix multiplication of size %d completed in %s\n", size, elapsed)
	fmt.Println("Result matrix checksum:", checksum(c))
}

func checksum(matrix [][]int) int {
	sum := 0
	for i := range matrix {
		for j := range matrix[i] {
			sum += matrix[i][j]
		}
	}
	return sum
}
