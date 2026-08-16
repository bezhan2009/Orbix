package system

import (
	"fmt"
	"github.com/gary23b/turtle"
	"image/color"
	"sync/atomic"
)

type Turtle struct {
	Command    string
	ValueFloat float64
	ValueStr   string
}

type TurtleWindowRequest struct {
	Width  int
	Height int
	Speed  float64
	Draw   []Turtle
}

var TurtleWindowChan = make(chan TurtleWindowRequest)

var (
	// TurtleStartChan Only used to START Ebiten once.
	TurtleStartChan = make(chan TurtleWindowRequest, 1)

	// TurtleDrawChan Used for every subsequent `turtle process`.
	TurtleDrawChan = make(chan TurtleWindowRequest, 32)

	// TurtleWindowState 0 = not started
	// 1 = running
	// 2 = closed permanently for this Orbix process
	TurtleWindowState atomic.Int32

	TurtleWindowWidth  int
	TurtleWindowHeight int
)

var (
	// TurtleCommands Orbix turtle commands
	TurtleCommands = []string{"forward", "left", "right", "backward", "speed", "penup", "pendown", "setcolor", "setwidth", "setheight"}

	// TurtleDraw Orbix's processed turtle draw commands
	TurtleDraw = []Turtle{}

	TurtleColorsMap = map[string]color.Color{
		"black":     turtle.Black,
		"red":       turtle.Red,
		"yellow":    turtle.Yellow,
		"green":     turtle.Green,
		"purple":    turtle.Purple,
		"blue":      turtle.Blue,
		"orange":    turtle.Orange,
		"white":     turtle.White,
		"cyan":      turtle.Cyan,
		"magenta":   turtle.Magenta,
		"pink":      turtle.Pink,
		"lime":      turtle.Lime,
		"teal":      turtle.Teal,
		"indigo":    turtle.Indigo,
		"violet":    turtle.Violet,
		"maroon":    turtle.Maroon,
		"olive":     turtle.Olive,
		"gold":      turtle.Gold,
		"salmon":    turtle.Salmon,
		"turquoise": turtle.Turquoise,
		"crimson":   turtle.Crimson,
		"tan":       turtle.Tan,
	}

	Speed = 100.00
)

func PrintTurtleDraw() {
	for i, cmd := range TurtleDraw {
		if cmd.ValueStr != "" {
			fmt.Printf("Command %d: %s %s\n", i, cmd.Command, cmd.ValueStr)
		} else {
			fmt.Printf("Command %d: %s %f\n", i, cmd.Command, cmd.ValueFloat)
		}
	}
}
