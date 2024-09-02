package main

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/run"
	"goCmd/src"
	"log"
	"os"
	"time"
)

const maxRetryAttempts = 5         // Maximum number of restart attempts
const retryDelay = 1 * time.Second // Delay before restarting

// OrbixLoop executes the basic Orbix logic with panic handling.
func OrbixLoop(red func(a ...interface{}) string) error {
	defer func() {
		if r := recover(); r != nil {
			PanicText := fmt.Sprintf("Panic recovered: %v", r)
			fmt.Printf("\n%s\n", red(PanicText))
			log.Printf("Panic recovered: %v", r)
		}
	}()

	run.Init()
	src.Orbix("", true)
	return nil
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

	for {
		if err := OrbixLoop(red); err != nil {
			ErrorText := fmt.Sprintf("Error occurred: %v", err)
			fmt.Println(red(ErrorText))
			log.Printf("Error occurred: %v", err)
		}

		attempts++
		if attempts >= maxRetryAttempts {
			fmt.Println(red("Max retry attempts reached. Exiting..."))
			log.Println("Max retry attempts reached. Exiting...")
			break
		}
		RestartText := fmt.Sprintf("Restarting Orbix in %v", magenta(retryDelay.Seconds()))
		fmt.Println(green(RestartText), green("seconds..."))
		log.Printf("Restarting Orbix in %v seconds...", retryDelay.Seconds())
		time.Sleep(retryDelay)
	}
}
