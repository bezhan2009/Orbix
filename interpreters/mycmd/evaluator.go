package mycmd

import (
	"fmt"
	"io/ioutil"
)

func Eval(node Node, env *Environment) Object {
	switch node := node.(type) {
	case *Program:
		return evalStatements(node.Statements, env)
	case *LetStatement:
		val := Eval(node.Value, env)
		env.store[node.Name.Value] = val
	case *Identifier:
		return env.store[node.Value]
	case *ExecuteStatement:
		return evalExecuteStatement(node, env)
	}

	return nil
}

func evalStatements(stmts []Statement, env *Environment) Object {
	var result Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)
	}

	return result
}

func evalExecuteStatement(stmt *ExecuteStatement, env *Environment) Object {
	scriptName := stmt.ScriptName

	script, err := ioutil.ReadFile(scriptName)
	if err != nil {
		fmt.Printf("Error reading script file: %v\n", err)
		return nil
	}

	input := string(script)

	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	return Eval(program, env)
}
