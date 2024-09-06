package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithSignaiture/Edit"
	"goCmd/utils"
	"os"
	"path/filepath"
)

// EditFileUtil - функция для редактирования файла или запуска beta версии
func EditFileUtil(commandArgs []string) {
	yellow := color.New(color.FgYellow).SprintFunc()

	if len(commandArgs) > 0 && commandArgs[0] == "beta" {
		// Попробуем найти и запустить "editBeta.exe"
		//executablePath, err := findExecutable("editBeta.exe")
		//if err != nil {
		//	fmt.Println(yellow("editBeta.exe not found in the current directory or path"))
		//	return
		//}

		// Запускаем editBeta.exe
		//err = exec.Command(executablePath).Run()
		fullCommand := []string{"py", "editBeta.py"}
		err := utils.ExternalCommand(fullCommand)
		if err != nil {
			fmt.Printf("Error running editBeta.exe: %v\n", err)
		}
		return
	}

	// Стандартная работа команды edit
	fmt.Println(yellow("to use beta version of command edit:"))
	fmt.Println(yellow("Usage: edit beta"))

	if len(commandArgs) < 1 {
		fmt.Println(yellow("Usage: edit <file>"))
		return
	}

	if err := Edit.File(commandArgs[0]); err != nil {
		fmt.Println(err)
		return
	}
}

// findExecutable - функция для нахождения пути к исполняемому файлу
func findExecutable(executableName string) (string, error) {
	// Получаем текущий рабочий каталог
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("unable to find executable path: %v", err)
	}

	// Определяем путь к директории, где находится запущенный файл
	dir := filepath.Dir(execPath)

	// Строим полный путь к "editBeta.exe"
	executablePath := filepath.Join(dir, executableName)

	// Проверяем, существует ли файл по этому пути
	if _, err := os.Stat(executablePath); os.IsNotExist(err) {
		return "", fmt.Errorf("executable not found")
	}

	return executablePath, nil
}
