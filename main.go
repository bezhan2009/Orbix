package main

import (
	"fmt"
	"goCmd/cmd/commands"
	"goCmd/system"
	"goCmd/utils"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Канал для получения уведомлений о сигнале
	sigChan := make(chan os.Signal, 1)
	// Отлавливаем SIGINT (Ctrl+C) и SIGTERM
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем обработку сигналов в отдельной горутине
	go func() {
		for sig := range sigChan {
			// Обрабатываем сигнал SIGINT (Ctrl+C)
			if sig == syscall.SIGINT {
				//fmt.Println("\nCtrl+C detected, but the program won't exit.")
			}
		}
	}()

	// Основная логика программы
	colors := system.GetColorsMap()

	fmt.Println(colors["cyan"]("Checking for updates..."))

	err := utils.CheckForUpdates()
	if err != nil {
		fmt.Println("Error checking for updates:", err)
	}

	err = commands.ChangeDirectory("scripts")
	if err != nil {
		fmt.Println("Error changing directory:", err)
		return
	}

	fmt.Print(colors["cyan"]("Preparing for launch"))
	utils.AnimatedPrintLong("...", "cyan")

	fmt.Println()

	command := []string{"go", "run", "orbix.go"}
	err = utils.ExternalCommand(command)
	if err != nil {
		fmt.Println("Error executing external command:", err)
	}
}
