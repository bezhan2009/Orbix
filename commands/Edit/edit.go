package Edit

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// File function for editing a file
func File(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	fmt.Println("Current file contents:")
	for i, line := range lines {
		fmt.Printf("%d: %s\n", i+1, line)
	}

	fmt.Println("\nEnter new content (type 'exit()' on a new line to save and exit):")
	newScanner := bufio.NewScanner(os.Stdin)
	var newContent []string
	for newScanner.Scan() {
		text := newScanner.Text()
		if strings.TrimSpace(text) == "exit()" {
			break
		}

		newContent = append(newContent, text)
	}
	if err := newScanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}

	file.Truncate(0)
	file.Seek(0, 0)
	writer := bufio.NewWriter(file)
	for _, line := range newContent {
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}
	writer.Flush()

	return nil
}
