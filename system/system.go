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
	Version         = "1.10.5"
	License         = "MIT"
	SystemName      = "Orbix"
	OperationSystem = runtime.GOOS
)

var (
	Beta          = false
	BetaVersion   = ""
	colors        = SetColorsMap()
	GlobalSession = Session{}
	SourcePath, _ = os.Getwd()
)

var (
	OrbixRunningUsersFileName = "running.env"
	OrbixFileNames            = map[string]uint{
		OrbixRunningUsersFileName: 1,
		"user.json":               1,
		".env":                    1,
		"commands.json":           1,
	}
)

func Init() *AppState {
	BetaVersion = string(strings.TrimSpace(strings.ToLower(os.Getenv("BETA"))))

	Beta = SetBetaVersion(colors)

	if UserDir == "" {
		UserDir, _ = os.Getwd()
	}

	// Initialization AppState
	return NewSystemData()
}

func SetBetaVersion(colors map[string]func(...interface{}) string) bool {
	if BetaVersion == "n" {
		return false
	}

	if len(os.Args) > 1 {
		if os.Args[len(os.Args)-1] == "beta" {
			return true
		}

		return false
	}

	for {
		var beta string
		reader := bufio.NewReader(os.Stdin)

		fmt.Print(colors["magenta"](fmt.Sprintf("Use Beta Version %s [Y/N]:", BetaVersion)))
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