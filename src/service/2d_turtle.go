package service

import (
	"fmt"
	"goCmd/system"
	"goCmd/utils"
	"strconv"
	"strings"

	"github.com/gary23b/turtle"
	"github.com/gary23b/turtle/turtlemodel"
	"github.com/hajimehoshi/ebiten/v2"
)

func checkAndValidateIndexTurtle(index string) (int, error) {
	id, err := strconv.Atoi(index)
	if err != nil {
		return 0, fmt.Errorf("invalid command ID. It should be a number")
	}

	if id < 0 {
		return 0, fmt.Errorf("invalid command ID. It should be a positive number")
	}

	if len(system.TurtleDraw) == 0 {
		return 0, fmt.Errorf("there are no turtle commands")
	}

	if id >= len(system.TurtleDraw) {
		return 0, fmt.Errorf(
			"invalid command ID. It should be between 0 and %d",
			len(system.TurtleDraw)-1,
		)
	}

	return id, nil
}

func TranslateToTurtle(
	commandArgs []string,
	add bool,
	update bool,
	delete bool,
	replace bool,
	idCommand string,
	idCommandForReplace string,
) error {
	if replace {
		id, err := checkAndValidateIndexTurtle(idCommand)
		if err != nil {
			return err
		}

		idReplacement, err := checkAndValidateIndexTurtle(idCommandForReplace)
		if err != nil {
			return err
		}

		storedTurtleReplacer := system.TurtleDraw[idReplacement]
		storedTurtle := system.TurtleDraw[id]

		system.TurtleDraw[id] = storedTurtleReplacer
		system.TurtleDraw[idReplacement] = storedTurtle

		return nil
	}

	// DELETE does not need command/value parsing.
	if delete {
		id, err := checkAndValidateIndexTurtle(idCommand)
		if err != nil {
			return err
		}

		system.TurtleDraw = append(
			system.TurtleDraw[:id],
			system.TurtleDraw[id+1:]...,
		)

		return nil
	}

	// ADD and UPDATE require:
	// <command> <value>
	if len(commandArgs) != 2 {
		return fmt.Errorf(
			"every turtle command has to have a value after the command. Example: forward 10, left 90, right 45",
		)
	}

	command := strings.TrimSpace(strings.ToLower(commandArgs[0]))
	value := strings.TrimSpace(commandArgs[1])

	// Validate turtle command.
	if !utils.InStr(command, system.TurtleCommands) {
		return fmt.Errorf("invalid turtle command: %s", command)
	}

	// Try numeric conversion.
	// String values such as `red` for setcolor are still stored
	// inside ValueStr.
	digits, err := strconv.ParseFloat(value, 64)
	if err != nil {
		digits = 0
	}

	turtleCommand := system.Turtle{
		Command:    command,
		ValueFloat: digits,
		ValueStr:   value,
	}

	// UPDATE
	if update {
		id, err := checkAndValidateIndexTurtle(idCommand)
		if err != nil {
			return err
		}

		system.TurtleDraw[id] = turtleCommand
		return nil
	}

	// ADD
	if add {
		system.TurtleDraw = append(
			system.TurtleDraw,
			turtleCommand,
		)

		return nil
	}

	return fmt.Errorf("no turtle operation specified")
}

func ProcessTurtle() error {
	height, ok := utils.InTurtle("setheight", system.TurtleDraw)
	if !ok {
		return fmt.Errorf("setheight required")
	}

	width, ok := utils.InTurtle("setwidth", system.TurtleDraw)
	if !ok {
		return fmt.Errorf("setwidth required")
	}

	speed, ok := utils.InTurtle("speed", system.TurtleDraw)
	if !ok {
		return fmt.Errorf("speed required")
	}

	draw := make([]system.Turtle, len(system.TurtleDraw))
	copy(draw, system.TurtleDraw)

	request := system.TurtleWindowRequest{
		Width:  int(width),
		Height: int(height),
		Speed:  speed,
		Draw:   draw,
	}

	switch system.TurtleWindowState.Load() {

	case 0:
		// First ever turtle process.
		if system.TurtleWindowState.CompareAndSwap(0, 1) {
			system.TurtleWindowWidth = request.Width
			system.TurtleWindowHeight = request.Height

			system.TurtleStartChan <- request
			return nil
		}

		// Another call managed to initialize it first.
		system.TurtleDrawChan <- request
		return nil

	case 1:
		// Turtle window already exists.

		if request.Width != system.TurtleWindowWidth ||
			request.Height != system.TurtleWindowHeight {

			return fmt.Errorf(
				"cannot change turtle window size while it is running; current size is %dx%d",
				system.TurtleWindowWidth,
				system.TurtleWindowHeight,
			)
		}

		system.TurtleDrawChan <- request
		return nil

	case 2:
		return fmt.Errorf(
			"turtle window was closed; Ebitengine cannot be started twice in the same Orbix process, restart Orbix to open it again",
		)
	}

	return fmt.Errorf("invalid turtle window state")
}

func StartTurtleWindow(initial system.TurtleWindowRequest) {
	// Optional but useful for Orbix:
	// keep Turtle above the terminal so the animation is immediately visible.
	ebiten.SetWindowFloating(true)

	turtle.Start(
		turtle.Params{
			Width:  initial.Width,
			Height: initial.Height,
		},
		func(window turtle.Window) {
			runTurtleWindow(window, initial)
		},
	)
}

func runTurtleWindow(
	window turtle.Window,
	initial system.TurtleWindowRequest,
) {
	t := window.NewTurtle()

	t.ShowTurtle()
	t.PenDown()

	// Draw the first request.
	drawTurtleRequest(window, t, initial)

	// Now KEEP this goroutine alive.
	//
	// Every future `turtle process`
	// sends another drawing here.
	for request := range system.TurtleDrawChan {
		drawTurtleRequest(window, t, request)
	}
}
func drawTurtleRequest(
	window turtle.Window,
	t turtlemodel.Turtle,
	request system.TurtleWindowRequest,
) {
	// Clear previous drawing.
	window.GetCanvas().ClearScreen(turtle.Black)

	// Reset turtle to starting position.
	t.PenUp()

	t.Teleport(0, 0)
	t.Angle(0)

	t.Color(turtle.White)
	t.Speed(request.Speed)

	t.PenDown()
	t.ShowTurtle()

	// Draw animation.
	for _, command := range request.Draw {
		if fn, exists := CommandMap[command.Command]; exists {
			fn(t, command)
		}
	}
}

func drawTurtleCommands(
	window turtle.Window,
	commands []system.Turtle,
	speed float64,
) {
	window.GetCanvas().ClearScreen(turtle.Black)

	t := window.NewTurtle()
	t.ShowTurtle()
	t.PenDown()
	t.Speed(speed)

	for _, command := range commands {
		if fn, exists := CommandMap[command.Command]; exists {
			fn(t, command)
		}
	}
}
