package system

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"runtime"
	"strings"
)

const (
	Version         = "1.7.6"
	License         = "MIT"
	SystemName      = "Orbix"
	OperationSystem = runtime.GOOS
)

var (
	BetaVersion = SetBetaVersion(GetColorsMap())
	colors      = SetColorsMap()
)

func SetBetaVersion(colors map[string]func(...interface{}) string) bool {
	if len(os.Args) > 0 {
		if os.Args[len(os.Args)-1] == "beta" {
			return true
		}

		return false
	}

	for {
		var beta string
		reader := bufio.NewReader(os.Stdin)

		fmt.Print(colors["magenta"]("Use Beta Version[Y/N]:"))
		beta, _ = reader.ReadString('\n')

		if beta == "" {
			return false
		}

		if strings.ToLower(strings.TrimSpace(beta)) == "y" {
			return true
		}

		return false
	}
}

func SetColorsMap() map[string]func(...interface{}) string {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	colors := map[string]func(...interface{}) string{
		"green":   green,
		"red":     red,
		"blue":    blue,
		"yellow":  yellow,
		"cyan":    cyan,
		"magenta": magenta,
	}

	return colors
}

func GetColorsMap() map[string]func(...interface{}) string {
	return colors
}