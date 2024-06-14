package mycmd

import (
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">> "

func StartRepl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := NewLexer(line)
		p := NewParser(l)

		program := p.ParseProgram()
		if len(p.errors) != 0 {
			printParserErrors(out, p.errors)
			continue
		}

		evaluated := Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
