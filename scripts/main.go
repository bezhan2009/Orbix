package main

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/run"
	"goCmd/src"
	"goCmd/structs"
	"log"
	"os"
	"time"
)

const maxRetryAttempts = 5         // Максимальное количество попыток перезапуска
const retryDelay = 1 * time.Second // Задержка перед перезапуском

// OrbixLoop запускает основную логику Orbix с обработкой паники.
func OrbixLoop(red func(a ...interface{}) string, panicChan chan any) {
	defer func() {
		if r := recover(); r != nil {
			PanicText := fmt.Sprintf("Panic recovered: %v", r)
			fmt.Printf("\n%s\n", red(PanicText))
			log.Printf("Panic recovered: %v", r)
			panicChan <- r // Отправляем панику в канал
		} else {
			panicChan <- nil // Если нет паники, отправляем nil
		}
	}()

	run.Init()
	src.Orbix("", true, structs.RebootedData{})
	panicChan <- nil // Если все прошло успешно
}

func main() {
	logFile, err := os.OpenFile("orbix.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	attempts := 0
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	// Канал для получения уведомлений о панике
	panicChan := make(chan any)

	for {
		isPanic := false

		// Запускаем OrbixLoop в отдельной горутине
		go OrbixLoop(red, panicChan)

		// Ожидаем результат работы OrbixLoop
		err := <-panicChan
		if err != nil {
			ErrorText := fmt.Sprintf("Error occurred: %v", err)
			fmt.Println(red(ErrorText))
			log.Printf("Error occurred: %v", err)
			isPanic = true
		}

		attempts++
		if attempts >= maxRetryAttempts {
			fmt.Println(red("Max retry attempts reached. Exiting..."))
			log.Println("Max retry attempts reached. Exiting...")
			break
		}

		if isPanic {
			RestartText := fmt.Sprintf("Restarting Orbix in %v", magenta(retryDelay.Seconds()))
			fmt.Println(green(RestartText), green("seconds..."))
			log.Printf("Restarting Orbix in %v seconds...", retryDelay.Seconds())
			time.Sleep(retryDelay)
		} else {
			break
		}
	}
}
