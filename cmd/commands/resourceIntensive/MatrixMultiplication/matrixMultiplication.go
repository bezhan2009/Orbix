package MatrixMultiplication

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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

// ParseMatrixString принимает строку и возвращает двумерный массив int
func ParseMatrixString(matrixStr string) ([][]int, error) {
	// Убираем пробелы и квадратные скобки в начале и в конце
	matrixStr = strings.TrimSpace(matrixStr)
	matrixStr = strings.TrimPrefix(matrixStr, "[")
	matrixStr = strings.TrimSuffix(matrixStr, "]")

	// Разделяем строки по закрывающим и открывающим квадратным скобкам
	rows := strings.Split(matrixStr, "],[")

	// Для хранения результирующего двумерного массива
	matrix := [][]int{}

	for _, row := range rows {
		// Убираем пробелы и оставшиеся скобки
		row = strings.TrimSpace(row)
		row = strings.TrimPrefix(row, "[")
		row = strings.TrimSuffix(row, "]")

		// Разделяем строку по запятым
		strNumbers := strings.Split(row, ",")

		// Для хранения чисел текущей строки
		numRow := []int{}

		for _, strNum := range strNumbers {
			// Преобразуем строку в число
			num, err := strconv.Atoi(strings.TrimSpace(strNum))
			if err != nil {
				return nil, err
			}
			numRow = append(numRow, num)
		}

		// Добавляем текущую строку в матрицу
		matrix = append(matrix, numRow)
	}

	return matrix, nil
}

func MatrixMulCommand(commandArgs []string) {
	var a, b [][]int
	var err error
	echo := false

	// Проверяем наличие аргумента echo
	for _, arg := range commandArgs {
		if arg == "echo=on" {
			echo = true
		} else if arg == "echo=off" {
			echo = false
		}
	}

	// Удаляем аргументы echo из списка команд
	nonEchoArgs := []string{}
	for _, arg := range commandArgs {
		if arg != "echo=on" && arg != "echo=off" {
			nonEchoArgs = append(nonEchoArgs, arg)
		}
	}

	if len(nonEchoArgs) >= 2 {
		// Преобразуем строки в матрицы
		a, err = ParseMatrixString(nonEchoArgs[0])
		if err != nil {
			fmt.Println("Ошибка при разборе строки первой матрицы:", err)
			return
		}

		b, err = ParseMatrixString(nonEchoArgs[1])
		if err != nil {
			fmt.Println("Ошибка при разборе строки второй матрицы:", err)
			return
		}

		// Проверка, что матрицы квадратные и имеют одинаковый размер
		if len(a) != len(b) || len(a[0]) != len(b[0]) || len(a) != len(a[0]) {
			fmt.Println("Ошибка: матрицы должны быть квадратными и иметь одинаковый размер.")
			return
		}
	} else {
		size := 500 // размер матрицы по умолчанию
		a = generateMatrix(size)
		b = generateMatrix(size)
		fmt.Println("Недостаточно аргументов, генерируются случайные матрицы размером", size)
	}

	start := time.Now()
	c := multiplyMatrices(a, b)
	elapsed := time.Since(start)

	fmt.Printf("Matrix multiplication completed in %s\n", elapsed)

	if echo {
		fmt.Println("Matrix a:", a)
		fmt.Println("Matrix b:", b)
	}

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
