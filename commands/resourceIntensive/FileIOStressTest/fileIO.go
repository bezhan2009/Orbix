package FileIOStressTest

import (
	"crypto/rand"
	"fmt"
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
	start := time.Now()
	err := writeLargeFile(filename, size)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	elapsed := time.Since(start)
	fmt.Printf("Wrote %d bytes to %s in %s\n", size, filename, elapsed)

	start = time.Now()
	data, err := readLargeFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	elapsed = time.Since(start)
	fmt.Printf("Read %d bytes from %s in %s\n", len(data), filename, elapsed)

	os.Remove(filename)
}
