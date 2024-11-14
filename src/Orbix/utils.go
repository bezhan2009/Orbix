package Orbix

import (
	"fmt"
	"goCmd/structs"
	"goCmd/system"
	"log"
	"os"
	"strings"
)

func RecoverFromThePanic(commandInput string,
	r any,
	echo bool,
	SD *system.AppState) {
	PanicText := fmt.Sprintf("Panic recovered: %v", r)
	fmt.Printf("\n%s\n", system.Red(PanicText))

	if system.RebootAttempts > system.MaxRetryAttempts {
		fmt.Println(system.Red("Max retry attempts reached. Exiting..."))
		log.Println("Max retry attempts reached. Exiting...")
		os.Exit(1)
	}

	system.RebootAttempts += 1

	fmt.Println(system.Yellow("Recovering from panic"))

	log.Printf("Panic recovered: %v", r)

	var reboot = structs.RebootedData{
		Username: system.UserName,
		Recover:  r,
		Prefix:   system.Prefix,
	}

	Orbix(strings.TrimSpace(commandInput),
		echo,
		reboot,
		SD)
}

func RestartAfterInitFn(SD *system.AppState,
	sessionData *system.AppState,
	rebooted structs.RebootedData,
	prefix,
	username string,
	echo bool) {
	SD.User = username
	SD.IsAdmin = sessionData.IsAdmin
	rebooted.Prefix = prefix
	if len(os.Args) > 1 {
		return
	}

	Orbix("",
		echo,
		rebooted,
		SD)
}
