package FileIOStressTest

import (
	"crypto/rand"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"time"
)

func writeLargeFile(filename string, size int) error {
	data := make([]byte, size)
	rand.Read(data)
	return ioutil.WriteFile(filename, data, 0644)
}

func readLargeFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func FileIOCommand(filename string, size int) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	start := time.Now()
	err := writeLargeFile(filename, size)
	if err != nil {
		fmt.Println(red("Error writing file:", err))
		return
	}
	elapsed := time.Since(start)
	printResult := fmt.Sprintf("Wrote %d bytes to %s in %s\n", size, filename, elapsed)
	fmt.Printf(green(printResult))

	start = time.Now()
	data, err := readLargeFile(filename)
	if err != nil {
		fmt.Println(red("Error reading file:", err))
		return
	}
	elapsed = time.Since(start)
	printElapsedResult := fmt.Sprintf("Read %d bytes from %s in %s\n", len(data), filename, elapsed)
	fmt.Printf(green(printElapsedResult))

	err = os.Remove(filename)
	if err != nil {
		fmt.Println(red("Error removing file:", err))
	}
}
