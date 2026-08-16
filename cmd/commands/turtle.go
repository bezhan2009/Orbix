package commands

import (
	"fmt"
	"goCmd/src/service"
	"goCmd/system"
	"strings"
)

func TurtleCommand(commandArgs []string) {
	if len(commandArgs) == 0 {
		fmt.Println(system.Yellow("Usage: turtle add <command> <value>"))
		fmt.Println(system.Yellow("Usage: turtle print"))
		fmt.Println(system.Yellow("Usage: turtle clear"))
		fmt.Println(system.Yellow("Usage: turtle process"))
		fmt.Println(system.Yellow("Usage: turtle update <id> <command> <value>"))
		fmt.Println(system.Yellow("Usage: turtle replace <id> <id_for_replace>"))
		fmt.Println(system.Yellow("Usage: turtle delete <id>"))
		return
	}

	clearedCommand := strings.TrimSpace(
		strings.ToLower(commandArgs[0]),
	)

	switch clearedCommand {
	case "print":
		system.PrintTurtleDraw()
		return

	case "clear":
		TurtleClearCommands()
		return

	case "process":
		TurtleProcessCommands()
		return

	case "add":
		TurtleAddCommand(commandArgs)
		return

	case "update":
		TurtleUpdateCommand(commandArgs)
		return

	case "replace":
		TurtleReplaceCommand(commandArgs)
		return

	case "delete":
		TurtleDeleteCommand(commandArgs)
		return

	default:
		fmt.Println(
			system.Red(
				"Invalid turtle command. Available commands are: add, update, delete, print, clear, process",
			),
		)
	}
}

func TurtleAddCommand(commandArgs []string) {
	if len(commandArgs) != 3 {
		fmt.Println(
			system.Yellow(
				"Usage: turtle add <command> <value>",
			),
		)

		fmt.Println(
			system.Yellow(
				"Turtle add command requires exactly 2 arguments. Example: forward 10, left 90, right 45",
			),
		)

		fmt.Println(
			system.Yellow(
				"The list of available commands:",
			),
		)

		for _, cmd := range system.TurtleCommands {
			fmt.Println(system.Yellow(" - " + cmd))
		}

		return
	}

	err := service.TranslateToTurtle(
		commandArgs[1:],
		true,
		false,
		false,
		false,
		"",
		"",
	)

	if err != nil {
		fmt.Println(system.Red(err))
		return
	}

	fmt.Println(
		system.Green(
			"Turtle command added successfully.",
		),
	)
}

func TurtleUpdateCommand(commandArgs []string) {
	if len(commandArgs) != 4 {
		fmt.Println(
			system.Yellow(
				"Usage: turtle update <id> <command> <value>",
			),
		)

		fmt.Println(
			system.Yellow(
				"Turtle update command requires exactly 3 arguments. Example: update 0 forward 100",
			),
		)

		fmt.Println(
			system.Yellow(
				"The list of available commands:",
			),
		)

		for _, cmd := range system.TurtleCommands {
			fmt.Println(system.Yellow(" - " + cmd))
		}

		return
	}

	err := service.TranslateToTurtle(
		commandArgs[2:],
		false,
		true,
		false,
		false,
		commandArgs[1],
		"",
	)

	if err != nil {
		fmt.Println(system.Red(err))
		return
	}

	fmt.Println(
		system.Green(
			"Turtle command updated successfully.",
		),
	)
}

func TurtleReplaceCommand(commandArgs []string) {
	if len(commandArgs) != 3 {
		fmt.Println(system.Yellow("Usage: turtle replace <id> <id_for_replace>"))
		fmt.Println(system.Yellow("Turtle replace command requires exactly 2 arguments. Example: replace 0 1"))
		return
	}

	err := service.TranslateToTurtle(
		nil,
		false,
		false,
		false,
		true,
		commandArgs[1],
		commandArgs[2],
	)

	if err != nil {
		fmt.Println(system.Red(err))
		return
	}

	fmt.Println(system.Green("Turtle command replaced successfully."))
}

func TurtleDeleteCommand(commandArgs []string) {
	if len(commandArgs) != 2 {
		fmt.Println(
			system.Yellow(
				"Usage: turtle delete <id>",
			),
		)

		fmt.Println(
			system.Yellow(
				"Turtle delete command requires exactly 1 argument. Example: delete 0",
			),
		)

		return
	}

	err := service.TranslateToTurtle(
		nil,
		false,
		false,
		true,
		false,
		commandArgs[1],
		"",
	)

	if err != nil {
		fmt.Println(system.Red(err))
		return
	}

	fmt.Println(
		system.Green(
			"Turtle command deleted successfully.",
		),
	)
}

func TurtleClearCommands() {
	fmt.Print(system.Yellow("Are you sure you want to clear all turtle commands? This action cannot be undone. Type 'yes' to confirm: "))
	var confirmation string
	fmt.Scan(&confirmation)

	if confirmation != "yes" {
		fmt.Println(system.Yellow("Turtle command clearing canceled."))
		return
	}

	system.TurtleDraw = []system.Turtle{}
	fmt.Println(system.Green("All turtle commands cleared successfully."))
}

func TurtleProcessCommands() {
	err := service.ProcessTurtle()
	if err != nil {
		fmt.Println(system.Red(err))
		return
	}

	fmt.Println(system.Green("Turtle commands processed successfully."))
}
