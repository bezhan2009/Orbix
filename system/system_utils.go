package system

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

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

// CheckPackageExists checks if Package is installed on the system
func CheckPackageExists(packageName string) bool {
	// Попробуем выполнить команду "git --version"
	_, err := exec.LookPath(packageName)
	if err != nil {
		log.Println(fmt.Sprintf("%s is not installed.", packageName))
		return false
	}

	return true
}
