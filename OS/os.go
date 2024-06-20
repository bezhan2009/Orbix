package OS

import "runtime"

func CheckOS() string {
	switch os := runtime.GOOS; os {
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
