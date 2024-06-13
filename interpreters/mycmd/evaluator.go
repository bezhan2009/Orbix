package mycmd

import (
	"fmt"
)

func Eval(node Node, env *Environment) Object {
	switch node := node.(type) {
	case *Program:
		return evalProgram(node, env)
	case *ExpressionStatement:
		return Eval(node.Expression, env)
	case *LetStatement:
		val := Eval(node.Value, env)
		env.Set(node.Name.Value, val)
	case *Identifier:
		return evalIdentifier(node, env)
	case *PrintStatement:
		val := Eval(node.Value, env)
		fmt.Println(val)
		return val
	}
	return nil
}

func evalProgram(program *Program, env *Environment) Object {
	var result Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)
	}

	return result
}

func evalIdentifier(ident *Identifier, env *Environment) Object {
	val, ok := env.Get(ident.Value)
	if !ok {
		return fmt.Sprintf("identifier not found: %s", ident.Value)
	}
	return val
}
