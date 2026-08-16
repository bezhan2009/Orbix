package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goCmd/system"
	"goCmd/utils"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/gary23b/turtle"
)

var (
	turtleProcessMu sync.Mutex
	turtleProcess   *exec.Cmd
)

const (
	defaultTurtleWindowWidth  = 800
	defaultTurtleWindowHeight = 600
	defaultTurtleSpeed        = 100.0
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
		height = defaultTurtleWindowHeight
	}

	width, ok := utils.InTurtle("setwidth", system.TurtleDraw)
	if !ok {
		width = defaultTurtleWindowWidth
	}

	speed, ok := utils.InTurtle("speed", system.TurtleDraw)
	if !ok {
		speed = defaultTurtleSpeed
	}

	if width <= 0 {
		return fmt.Errorf("invalid turtle window width")
	}

	if height <= 0 {
		return fmt.Errorf("invalid turtle window height")
	}

	if speed <= 0 {
		return fmt.Errorf("invalid turtle speed")
	}

	draw := make([]system.Turtle, len(system.TurtleDraw))
	copy(draw, system.TurtleDraw)

	request := system.TurtleWindowRequest{
		Width:  int(width),
		Height: int(height),
		Speed:  speed,
		Draw:   draw,
	}

	return startTurtleRendererProcess(request)
}

func startTurtleRendererProcess(
	request system.TurtleWindowRequest,
) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf(
			"failed to encode turtle commands: %w",
			err,
		)
	}

	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf(
			"failed to locate Orbix executable: %w",
			err,
		)
	}

	cmd := exec.Command(
		executable,
		"--turtle-renderer",
	)

	cmd.Stdin = bytes.NewReader(payload)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	turtleProcessMu.Lock()

	if turtleProcess != nil &&
		turtleProcess.Process != nil {

		_ = turtleProcess.Process.Kill()
	}

	if err := cmd.Start(); err != nil {
		turtleProcessMu.Unlock()

		return fmt.Errorf(
			"failed to start turtle renderer: %w",
			err,
		)
	}

	turtleProcess = cmd

	turtleProcessMu.Unlock()

	go func(current *exec.Cmd) {
		_ = current.Wait()

		turtleProcessMu.Lock()
		defer turtleProcessMu.Unlock()

		if turtleProcess == current {
			turtleProcess = nil
		}
	}(cmd)

	return nil
}

func RunTurtleRendererFromStdin() error {
	var request system.TurtleWindowRequest

	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		return fmt.Errorf(
			"failed to decode turtle render request: %w",
			err,
		)
	}

	if request.Width <= 0 {
		return fmt.Errorf("invalid turtle window width")
	}

	if request.Height <= 0 {
		return fmt.Errorf("invalid turtle window height")
	}

	if request.Speed <= 0 {
		return fmt.Errorf("invalid turtle speed")
	}

	system.Speed = request.Speed

	turtle.Start(
		turtle.Params{
			Width:  request.Width,
			Height: request.Height,
		},
		func(window turtle.Window) {
			drawTurtleRequest(
				window,
				request,
			)
		},
	)

	return nil
}

func drawTurtleRequest(
	window turtle.Window,
	request system.TurtleWindowRequest,
) {
	window.GetCanvas().ClearScreen(turtle.Black)

	t := window.NewTurtle()

	t.ShowTurtle()
	t.PenDown()
	t.Speed(request.Speed)

	for _, command := range request.Draw {
		if fn, exists := CommandMap[command.Command]; exists {
			fn(t, command)
		}
	}
}
