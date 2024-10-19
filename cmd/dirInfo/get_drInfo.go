package dirInfo

import (
	"goCmd/system"
	"os"
	"path/filepath"
	"strings"
	_ "strings"
	_ "unicode/utf8"
)

func CmdDir(dir string) string {
	if system.OperationSystem != "windows" {
		return dir
	}

	// Считаем количество слэшей и обрезаем путь
	parts := strings.SplitN(dir, "\\", 4)
	if len(parts) >= 4 {
		return parts[3]
	}

	return dir
}

// Пример функции декодирования пути для Windows
func decodeWindowsPath(dir string) string {
	// Преобразуем путь в нормализованный
	return filepath.ToSlash(dir)
}

func CmdUser(dir string) string {
	if system.OperationSystem == "linux" {
		return os.Getenv("USER")
	} else {
		// Считаем количество слэшей и получаем имя пользователя между 1 и 2 слэшом
		parts := strings.SplitN(dir, "\\", 3)
		if len(parts) >= 2 {
			return parts[1]
		}

		return "Unknown"
	}
}
