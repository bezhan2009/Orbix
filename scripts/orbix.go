package main

import (
	"fmt"
	_chan "goCmd/chan"
	"goCmd/internal/run"
	"goCmd/scripts/handlers"
	"goCmd/src/Orbix"
	"goCmd/src/service"
	"goCmd/src/user"
	"goCmd/structs"
	"goCmd/system"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

const maxRetryAttempts = system.MaxRetryAttempts // Maximum number of restart attempts
const retryDelay = system.RetryDelay             // Delay before restart

func OrbixLoop(panicChan chan any,
	appState *system.AppState) {
	colorsMap := system.GetColorsMap()
	red := colorsMap["red"]

	go func() {
		for {
			time.Sleep(retryDelay)
			if _chan.IsVarsFnUpd {
				time.Sleep(retryDelay)
				_chan.SaveVarsFn()
				_chan.IsVarsFnUpd = false
			}
		}
	}()

	defer func() {
		if system.Debug {
			return
		}

		if r := recover(); r != nil {
			user.DeleteUserFromRunningFile(system.UserName)
			PanicText := fmt.Sprintf("Panic recovered: %v", r)
			fmt.Printf("\n%s\n", red(PanicText))
			log.Printf("Panic recovered: %v", r)
			panicChan <- r
		} else {
			panicChan <- nil
		}
	}()

	Orbix.Orbix("",
		true,
		structs.RebootedData{},
		appState)

	panicChan <- nil
}

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Printf("Error creating CPU profile: %v", err)
	}
	defer f.Close()

	// Запуск профилирования CPU.
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Printf("Error starting CPU profile: %v", err)
	}
	defer pprof.StopCPUProfile()

	// Initialization Orbix
	run.Init()

	// Initialization system vars
	appState := system.Init()

	colors := system.GetColorsMap()

	green := colors["green"]
	red := colors["red"]
	magenta := colors["magenta"]

	if len(os.Args) > 1 {
		args := os.Args[1:]
		command := ""
		for i, arg := range args {
			if i == len(args)-1 && arg == "beta" {
				continue
			}

			command += arg + " "
		}

		Orbix.Orbix(command,
			true,
			structs.RebootedData{},
			appState)
		return
	}

	// Initialization cache files
	run.CacheInit()

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/web-cmd", handlers.HandlerWebOrbix)
	go func() {
		var err error
		var portInt int
		for {
			URL := fmt.Sprintf("localhost:%s", system.Port)
			err = http.ListenAndServe(URL, nil)
			if err != nil {
				fmt.Println()
				fmt.Println(red(fmt.Sprintf("Server failed to start: %v",
					err)))
				portInt, err = strconv.Atoi(system.Port)
				if err != nil {
					fmt.Println(red(fmt.Sprintf("PortError: %v",
						err)))
					system.Port = "6060"
					system.ErrorStartingServer = true
					continue
				}

				portInt += 1
				port := fmt.Sprint(portInt)
				system.Port = port

				system.ErrorStartingServer = true
				continue
			}
			break
		}
	}()

	go func() {
		for {
			time.Sleep(retryDelay)
			if system.ErrorStartingServer {
				fmt.Println(green("The server was able to resolve the error, and now server is listening on port " + system.Port))
				fmt.Print(green(" > "))
				break
			}
		}
	}()

	logFile, err := os.OpenFile("orbix.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer func() {
		err = logFile.Close()
		if err != nil {
			return
		}
	}()

	log.SetOutput(logFile)

	panicChan := make(chan any)

	system.Localhost = fmt.Sprintf("http://localhost:%s", system.Port)

	for {
		isPanic := false

		go OrbixLoop(panicChan, appState)

		var orbixErr any
		orbixRunning := true

		for orbixRunning {
			select {
			case request := <-system.TurtleStartChan:
				service.StartTurtleWindow(request)

				// If we reach here, the user CLOSED the Turtle window.
				//
				// Ebiten cannot be started again in this process.
				system.TurtleWindowState.Store(2)

			case orbixErr = <-panicChan:
				orbixRunning = false
			}
		}

		if orbixErr != nil {
			errorText := fmt.Sprintf(
				"Error occurred: %v",
				orbixErr,
			)

			fmt.Println(red(errorText))
			log.Printf(
				"Error occurred: %v",
				orbixErr,
			)

			isPanic = true
		}

		system.Attempts++

		if system.Attempts > maxRetryAttempts {
			fmt.Println(
				red(
					"Max retry attempts reached. Exiting...",
				),
			)

			log.Println(
				"Max retry attempts reached. Exiting...",
			)

			break
		}

		if isPanic {
			restartText := fmt.Sprintf(
				"Restarting Orbix in %v",
				magenta(retryDelay.Seconds()),
			)

			fmt.Println(
				green(restartText),
				green("seconds..."),
			)

			log.Printf(
				"Restarting Orbix in %v seconds...",
				retryDelay.Seconds(),
			)

			time.Sleep(retryDelay)
		} else {
			break
		}
	}

	go func() {
		time.Sleep(retryDelay)
		if !system.OrbixWorking {
			user.DeleteUserFromRunningFile(system.UserName)
			*_chan.LoopData.IsWorking = false
		}
	}()

	defer func() {
		user.DeleteUserFromRunningFile(system.UserName)
		_chan.SaveVarsFn()
		_chan.UpdateChan("scripts__orbix_func")
		os.Exit(1)
	}()
}
