package system

import "github.com/fatih/color"

var (
	Red         func(a ...interface{}) string
	Green       func(a ...interface{}) string
	Yellow      func(a ...interface{}) string
	Blue        func(a ...interface{}) string
	Magenta     func(a ...interface{}) string
	Cyan        func(a ...interface{}) string
	RedBold     func(a ...interface{}) string
	GreenBold   func(a ...interface{}) string
	YellowBold  func(a ...interface{}) string
	BlueBold    func(a ...interface{}) string
	MagentaBold func(a ...interface{}) string
	CyanBold    func(a ...interface{}) string
)

func SetColorsMap() map[string]func(...interface{}) string {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	redBold := color.New(color.FgRed, color.Bold).SprintFunc()
	yellowBold := color.New(color.FgYellow, color.Bold).SprintFunc()
	cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	greenBold := color.New(color.FgGreen, color.Bold).SprintFunc()
	magentaBold := color.New(color.FgMagenta, color.Bold).SprintFunc()
	blueBold := color.New(color.FgBlue, color.Bold).SprintFunc()

	colors := map[string]func(...interface{}) string{
		"green":       green,
		"red":         red,
		"blue":        blue,
		"yellow":      yellow,
		"cyan":        cyan,
		"magenta":     magenta,
		"redBold":     redBold,
		"yellowBold":  yellowBold,
		"cyanBold":    cyanBold,
		"magentaBold": magentaBold,
		"blueBold":    blueBold,
		"greenBold":   greenBold,
	}

	return colors
}

func GetColorsMap() map[string]func(...interface{}) string {
	return colors
}

func InitColors() {
	colors = GetColorsMap()

	Red = colors["red"]
	Yellow = colors["yellow"]
	Cyan = colors["cyan"]
	Green = colors["green"]
	Magenta = colors["magenta"]
	Blue = colors["blue"]
	RedBold = colors["redBold"]
	YellowBold = colors["yellowBold"]
	CyanBold = colors["cyanBold"]
	MagentaBold = colors["magentaBold"]
	BlueBold = colors["blueBold"]
}
