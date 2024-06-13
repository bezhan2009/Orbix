package mycmd

import (
	"fmt"
	"io/ioutil"
)

func Interpreter(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: mycmd <script.mycmd>")
		return
	}

	fileName := commandArgs[0]

	script, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading script file: %v\n", err)
		return
	}

	input := string(script)

	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	env := NewEnvironment()
	Eval(program, env)
}
