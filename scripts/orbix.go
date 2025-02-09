package main

import (
	"fmt"
	_chan "goCmd/chan"
	"goCmd/internal/run"
	"goCmd/src/Orbix"
	"goCmd/src/user"
	"goCmd/structs"
	"goCmd/system"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

const maxRetryAttempts = system.MaxRetryAttempts // Maximum number of restart attempts
const retryDelay = system.RetryDelay             // Delay before restart

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Handler for rendering index.html
func indexHandler(w http.ResponseWriter,
	r *http.Request) {
	setHeaders(w)
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return
	}
}

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

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", indexHandler)
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
				fmt.Print(green(" >"))
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

		// Launching OrbixLoop in a separate goroutine
		go OrbixLoop(panicChan, appState)

		// We are waiting for the result of OrbixLoop's work
		err := <-panicChan
		if err != nil {
			ErrorText := fmt.Sprintf("Error occurred: %v", err)
			fmt.Println(red(ErrorText))
			log.Printf("Error occurred: %v", err)
			isPanic = true
		}

		system.Attempts++
		if system.Attempts > maxRetryAttempts {
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
