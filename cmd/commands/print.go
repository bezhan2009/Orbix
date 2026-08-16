package commands

import (
	"fmt"
	"goCmd/system"
	"goCmd/utils"
	"regexp"
	"strings"

	"github.com/common-nighthawk/go-figure"
)

// Print prints text with support for:
//
//   - colors
//   - 2D fonts
//   - 3D fonts
//   - animated output
//
// Examples:
//
//	print Hello world
//	print red:Hello world
//	print red:Hello world; blue:Goodbye
//	print red:Hello font=2d
//	print green:Orbix font=3d
//	print red:Hello world animate
//	print animate green:Orbix font=3d
func Print(commandArgs []string) {
	var (
		font          string
		animatedPrint bool
		colorFuncs    = system.GetColorsMap()

		textArgs []string
	)

	if len(commandArgs) == 0 {
		fmt.Println(system.Yellow("Usage: print <text> [font=2d|3d] [animate]"))
		return
	}

	/*
		Extract special arguments first.

		IMPORTANT:
		Do not modify commandArgs while iterating over it.
		The previous implementation removed elements from the
		slice during range, which caused skipped and incorrectly
		indexed arguments.
	*/
	for i := 0; i < len(commandArgs); i++ {
		rawArg := commandArgs[i]
		arg := strings.TrimSpace(rawArg)

		if arg == "" {
			continue
		}

		// Detect:
		// font=2d
		// font=3d
		if key, value, found := strings.Cut(arg, "="); found {
			key = strings.TrimSpace(strings.ToLower(key))

			if key == "font" {
				font = strings.TrimSpace(strings.ToLower(value))

				switch font {
				case "":
					fmt.Println(
						system.Yellow(
							"Font value cannot be empty. Available fonts: 2d, 3d",
						),
					)
					return

				case "2d", "3d":
					// Valid font.

				default:
					fmt.Println(
						system.Yellow(
							fmt.Sprintf(
								"Invalid font %q. Available fonts: 2d, 3d",
								font,
							),
						),
					)
					return
				}

				continue
			}
		}

		if strings.EqualFold(arg, "font") {
			if i+1 >= len(commandArgs) {
				fmt.Println(
					system.Yellow(
						"Font value cannot be empty. Available fonts: 2d, 3d",
					),
				)
				return
			}

			value := strings.TrimSpace(
				strings.ToLower(commandArgs[i+1]),
			)

			switch value {
			case "2d", "3d":
				font = value
				i++
				continue

			default:
				fmt.Println(
					system.Yellow(
						fmt.Sprintf(
							"Invalid font %q. Available fonts: 2d, 3d",
							value,
						),
					),
				)
				return
			}
		}

		// Detect:
		// animate
		if strings.EqualFold(arg, "animate") {
			animatedPrint = true
			continue
		}

		// Everything else belongs to the actual text.
		textArgs = append(textArgs, rawArg)
	}

	if len(textArgs) == 0 {
		fmt.Println(system.Yellow("Nothing to print."))
		return
	}

	/*
		Orbix may split:

		    red:Hello my friend

		into:

		    []string{
		        "red:Hello",
		        "my",
		        "friend",
		    }

		So we join everything back together before parsing
		colors and semicolon-separated sections.
	*/
	fullText := strings.Join(textArgs, " ")

	/*
		A semicolon separates independently styled sections.

		Example:

		    red:Hello world; blue:Orbix

		becomes:

		    red:Hello world
		    blue:Orbix
	*/
	parts := strings.Split(fullText, ";")

	printedSomething := false

	for _, rawPart := range parts {
		part := strings.TrimSpace(rawPart)

		if part == "" {
			continue
		}

		text := part

		// Default means no terminal color.
		colorFunc := fmt.Sprint

		/*
			Try to interpret:

			    red:Hello

			as:

			    color = red
			    text  = Hello

			BUT only if "red" is actually a known Orbix color.

			This means ordinary strings such as:

			    Time: 12:30
			    https://example.com

			are not accidentally interpreted as color syntax.
		*/
		if possibleColor, possibleText, found := strings.Cut(part, ":"); found {
			colorName := strings.TrimSpace(
				strings.ToLower(possibleColor),
			)

			if foundColorFunc, exists := colorFuncs[colorName]; exists {
				colorFunc = foundColorFunc
				text = strings.TrimSpace(possibleText)
			}
		}

		if text == "" {
			continue
		}

		printWithFont(
			text,
			font,
			colorFunc,
			animatedPrint,
		)

		printedSomething = true
	}

	if printedSomething {
		fmt.Println()
	}
}

// printWithFont prints a single text section using the requested font,
// color function and animation mode.
func printWithFont(
	text string,
	font string,
	colorFunc func(a ...interface{}) string,
	animate bool,
) {
	font = strings.TrimSpace(strings.ToLower(font))

	/*
		ASCII-art fonts do not properly support arbitrary Unicode
		characters, so keep only characters that are safe for
		go-figure.
	*/
	if font == "2d" || font == "3d" {
		re := regexp.MustCompile(`[^a-zA-Z0-9 !@#+$%^&*()_]`)

		filteredText := re.ReplaceAllString(text, "")

		if strings.TrimSpace(filteredText) == "" {
			fmt.Println(
				system.Yellow(
					"Invalid characters for the selected font.",
				),
			)
			return
		}

		text = filteredText
	}

	switch font {

	// 3D ASCII-art output.
	case "3d":
		myFigure := figure.NewFigure(
			text,
			"larry3d",
			true,
		)

		output := colorFunc(myFigure.String())

		if animate {
			utils.PrintAnim(output)
		} else {
			fmt.Print(output)
		}

	// 2D ASCII-art output.
	case "2d":
		myFigure := figure.NewFigure(
			text,
			"",
			true,
		)

		output := colorFunc(myFigure.String())

		if animate {
			utils.PrintAnim(output)
		} else {
			fmt.Print(output)
		}

	// Normal text.
	default:
		output := colorFunc(text)

		if animate {
			utils.PrintAnim(output)
		} else {
			fmt.Print(output)
		}

		// Separate multiple styled sections.
		fmt.Print(" ")
	}
}
