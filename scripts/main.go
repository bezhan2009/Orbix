package main

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/run"
	"goCmd/src"
	"time"
)

// OrbixLoop executes the basic Orbix logic in a loop with panic handling.
func OrbixLoop() {
	defer func() {
		if r := recover(); r != nil {
			red := color.New(color.FgRed).SprintFunc()
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Println(green("Recovery from panic:"), red(r))
			// Adding a delay before restarting
			time.Sleep(1 * time.Second)
		}
	}()

	run.Init()
	src.Orbix("", true)
}

func main() {
	for {
		OrbixLoop() // We perform the basic logic with the processing of panics
	}
}
