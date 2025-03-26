package OS

import (
	"fmt"
	"goCmd/system/errs"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
)

// Получение списка процессов (кроссплатформенно)
func getProcesses() {
	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Ошибка при получении списка процессов: %v", err)
	}
	fmt.Println("Активные процессы:")
	for _, p := range processes {
		name, _ := p.Name()
		fmt.Printf("PID: %d, Имя: %s\n", p.Pid, name)
	}
}

// GetEnvVariable Чтение переменной окружения
func GetEnvVariable(varName string) (string, error) {
	value, exists := os.LookupEnv(varName)
	if exists {
		return value, nil
	} else {
		return "", errs.VariableDoesNotExist
	}
}

// SetEnvVariable Установка переменной окружения
func SetEnvVariable(varName, value string) error {
	err := os.Setenv(varName, value)
	if err != nil {
		log.Printf("Ошибка при установке переменной окружения: %v", err)

		return err
	}

	return nil
}

// GetFileInfo Получение информации о файле
func GetFileInfo(filePath string) (string, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Error while getting information about file: %v", err)

		return "", err
	}

	return fmt.Sprintf("File: %s\nSize: %d байт\nMod time: %s\n",
		info.Name(), info.Size(), info.ModTime()), nil
}

// IsPortOpen Проверка, занят ли порт
func IsPortOpen(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return true // Порт занят
	}
	listener.Close()
	return false // Порт свободен
}

// Получение информации о системе
func getSystemInfo() {
	fmt.Println("Операционная система:", CheckOS())
	fmt.Println("Архитектура:", runtime.GOARCH)
	fmt.Println("Количество CPU:", runtime.NumCPU())
}

// Проверка загрузки CPU
func getCPUUsage() {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Fatalf("Ошибка при получении загрузки CPU: %v", err)
	}
	fmt.Printf("Загрузка CPU: %.2f%%\n", percent[0])
}

// Чтение текущего каталога
func getCurrentDirectory() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Ошибка при получении текущего каталога: %v", err)
	}
	fmt.Println("Текущий каталог:", dir)
}

// Завершение процесса по PID
func killProcess(pid int) {
	p, err := os.FindProcess(pid)
	if err != nil {
		log.Fatalf("Ошибка при завершении процесса: %v", err)
	}
	err = p.Kill()
	if err != nil {
		log.Fatalf("Не удалось завершить процесс с PID %d: %v", pid, err)
	}
	fmt.Printf("Процесс с PID %d завершён\n", pid)
}

// Запуск команды в терминале (кроссплатформенно)
func runCommand(command string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Ошибка при выполнении команды: %v", err)
	}
	fmt.Printf("Вывод команды:\n%s\n", string(output))
}
