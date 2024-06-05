package Edit

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// File function for editing a file
func File(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR, 0644) // Открываем файл для чтения и записи
	if err != nil {
		return fmt.Errorf("error opening file: %v", err) // Если возникла ошибка при открытии файла, возвращаем ошибку
	}
	defer file.Close() // Гарантируем закрытие файла после завершения функции

	// Read current file contents
	scanner := bufio.NewScanner(file) // Создаем сканер для чтения файла построчно
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text()) // Считываем каждую строку файла и добавляем в массив lines
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err) // Если возникла ошибка при чтении файла, возвращаем ошибку
	}

	// Display current contents
	fmt.Println("Current file contents:") // Выводим текущие содержимое файла
	for i, line := range lines {
		fmt.Printf("%d: %s\n", i+1, line) // Печатаем строки с номерами
	}

	// User editing
	fmt.Println("\nEnter new content (type 'exit()' on a new line to save and exit):") // Просим пользователя ввести новое содержимое
	newScanner := bufio.NewScanner(os.Stdin)                                           // Создаем сканер для чтения пользовательского ввода
	var newContent []string
	for newScanner.Scan() {
		text := newScanner.Text()
		if strings.TrimSpace(text) == "exit()" { // Если пользователь ввел exit(), выходим из цикла
			break
		}
		newContent = append(newContent, text) // Добавляем введенные строки в массив newContent
	}
	if err := newScanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %v", err) // Если возникла ошибка при чтении ввода, возвращаем ошибку
	}

	// Write new contents to file
	file.Truncate(0)                // Очищаем файл
	file.Seek(0, 0)                 // Устанавливаем указатель в начало файла
	writer := bufio.NewWriter(file) // Создаем буферизированный писатель
	for _, line := range newContent {
		_, err = writer.WriteString(line + "\n") // Записываем новые строки в файл
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err) // Если возникла ошибка при записи, возвращаем ошибку
		}
	}
	writer.Flush() // Сбрасываем буферизированные данные в файл

	return nil
}
