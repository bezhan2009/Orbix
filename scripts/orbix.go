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
	// Инициализация AppState
	appState := system.NewSystemData()

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", indexHandler)
	go func() {
		URL := fmt.Sprintf("localhost:%s", system.Port)
		err := http.ListenAndServe(URL, nil)
		if err != nil {
			log.Fatalf("Server failed to start: %v", err)
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

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	panicChan := make(chan any)

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
}
