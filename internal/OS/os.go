package OS

import (
	"os"
	"runtime"
	"strconv"
	"strings"
)

func CheckOS() string {
	switch goos := runtime.GOOS; goos {
	case "windows":
		return "windows"
	case "darwin":
		return "darwin"
	case "linux":
		return "linux"
	default:
		return ""
	}
}

func KillProcess(pidStr string) error {
	pid, err := strconv.Atoi(strings.TrimSpace(pidStr))
	if err != nil {
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = process.Kill()
	if err != nil {
		return err
	}

	return nil
}
