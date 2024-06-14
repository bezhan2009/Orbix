package mycmd

import (
	"fmt"
	"io/ioutil"
	"os"
)

func Start() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mycmd <script.mycmd>")
		return
	}

	fileName := os.Args[1]
	script, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading script file: %v\n", err)
		return
	}

	input := string(script)
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrorsP(p.Errors())
		return
	}

	env := NewEnvironment()
	result := Eval(program, env)

	if result != nil {
		fmt.Println(result.Inspect())
	}
}

func printParserErrorsP(errors []string) {
	for _, msg := range errors {
		fmt.Println("\t" + msg)
	}
}
