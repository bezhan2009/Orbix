package service

import (
	"fmt"
	"goCmd/system"

	"github.com/gary23b/turtle/turtlemodel"
)

var (
	CommandMap = map[string]func(t turtlemodel.Turtle, value system.Turtle){
		"forward":  Forward,
		"backward": Backward,
		"left":     Left,
		"right":    Right,
		"setcolor": SetColor,
		"speed":    SetSpeed,
		"penup":    func(t turtlemodel.Turtle, value system.Turtle) { PenUp(t) },
		"pendown":  func(t turtlemodel.Turtle, value system.Turtle) { PenDown(t) },
	}
)

func Forward(t turtlemodel.Turtle, value system.Turtle) {
	t.Forward(value.ValueFloat)
}

func Backward(t turtlemodel.Turtle, value system.Turtle) {
	t.Backward(value.ValueFloat)
}

func Left(t turtlemodel.Turtle, value system.Turtle) {
	t.Left(value.ValueFloat)
}

func Right(t turtlemodel.Turtle, value system.Turtle) {
	t.Right(value.ValueFloat)
}

func SetColor(t turtlemodel.Turtle, value system.Turtle) {
	turtleColor, ok := system.TurtleColorsMap[value.ValueStr]
	if !ok {
		fmt.Println(system.Red("Error: Invalid color name " + value.ValueStr + ". Please use a valid color name."))
	} else {
		t.Color(turtleColor)
	}
}

func SetSpeed(t turtlemodel.Turtle, value system.Turtle) {
	t.Speed(value.ValueFloat)
}

func PenUp(t turtlemodel.Turtle) {
	t.PenUp()
}

func PenDown(t turtlemodel.Turtle) {
	t.PenDown()
}
