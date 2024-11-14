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

	if len(os.Args) > 1 {
		err := commands.ChangeDirectory("scripts")
		if err != nil {
			fmt.Println("Error changing directory:", err)
			return
		}

		args := os.Args[1:]

		command := []string{"go", "run", "orbix.go"}
		command = append(command, args...)

		err = utils.ExternalCommand(command)
		if err != nil {
			fmt.Println("Error executing external command:", err)
		}

		return
	}

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

	fmt.Println(colors["cyan"](""))

	command := []string{"go", "run", "orbix.go"}
	go func() {
		err = utils.ExternalCommand(command)
		if err != nil {
			fmt.Println("Error executing external command:", err)
		}
	}()
}
