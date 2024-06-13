package mycmd

import (
	"fmt"
)

func Interpreter(input string) {
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.errors) != 0 {
		for _, msg := range p.errors {
			fmt.Println(msg)
		}
		return
	}

	env := NewEnvironment()
	Eval(program, env)
}
