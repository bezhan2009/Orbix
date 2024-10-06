package main

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/internal/run"
	"goCmd/src"
	"goCmd/structs"
	"goCmd/system"
	"html/template"
	"log"
	"net/http"
	"os"
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
func indexHandler(w http.ResponseWriter, r *http.Request) {
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

// OrbixLoop runs the basic Orbix logic with panic handling.
func OrbixLoop(red func(a ...interface{}) string, panicChan chan any, appState *system.AppState) {
	defer func() {
		if r := recover(); r != nil {
			src.RemoveUserFromRunningFile(system.UserName)
			PanicText := fmt.Sprintf("Panic recovered: %v", r)
			fmt.Printf("\n%s\n", red(PanicText))
			log.Printf("Panic recovered: %v", r)
			panicChan <- r
		} else {
			panicChan <- nil
		}
	}()

	run.Init()
	src.Orbix("", true, structs.RebootedData{}, appState)
	panicChan <- nil
}

func main() {
	// Initialization AppState
	appState := system.NewSystemData()

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

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
				fmt.Println(red(fmt.Sprintf("Server failed to start: %v", err)))
				portInt, err = strconv.Atoi(system.Port)
				if err != nil {
					fmt.Println(red(fmt.Sprintf("PortError: %v", err)))
					break
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
				fmt.Print(magenta("enable secure[Y/N]: "))
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
		go OrbixLoop(red, panicChan, appState)

		// We are waiting for the result of OrbixLoop's work
		err := <-panicChan
		if err != nil {
			ErrorText := fmt.Sprintf("Error occurred: %v", err)
			fmt.Println(red(ErrorText))
			log.Printf("Error occurred: %v", err)
			isPanic = true
		}

		system.Attempts++
		if system.Attempts >= maxRetryAttempts {
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
		time.Sleep(time.Second * 1)
		if !system.OrbixWorking {
			src.RemoveUserFromRunningFile(system.UserName)
			os.Exit(1)
		}
	}()

	defer func() {
		src.RemoveUserFromRunningFile(system.UserName)
	}()
}
