package cache

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var (
	// CommandCache хранит для каждой команды её полный путь.
	CommandCache = make(map[string]string)
	cacheMutex   sync.RWMutex
	CacheOnce    sync.Once
)

// InitCommandCache сканирует текущую директорию и каталоги из PATH и заполняет кеш.
func InitCommandCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Очищаем кеш.
	CommandCache = make(map[string]string)

	// Кэширование файлов из текущей директории.
	curDir, err := os.Getwd()
	if err == nil {
		_ = filepath.Walk(curDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				CommandCache[info.Name()] = path
			}
			return nil
		})
	}

	// Кэширование файлов из директорий PATH.
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))
	for _, dir := range paths {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			fullPath := filepath.Join(dir, entry.Name())
			CommandCache[entry.Name()] = fullPath
		}
	}
}

// GetCommandFromCache возвращает полный путь для заданной команды (ключ) из кеша.
func GetCommandFromCache(name string) string {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	if path, ok := CommandCache[name]; ok {
		return path
	}
	return ""
}

// watchDirectories запускает наблюдатель за изменениями в указанных директориях.
func watchDirectories(dirs []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Error creating watcher:", err)
		return
	}
	defer watcher.Close()

	// Добавляем все нужные директории в наблюдение.
	for _, dir := range dirs {
		err = watcher.Add(dir)
		if err != nil {
			log.Println("Error adding directory to watcher:", err)
		}
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// Если произошли изменения (создание, удаление, изменение, переименование) – обновляем кеш.
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename {
				log.Println("Directory change detected in", event.Name, "-> updating command cache")
				// Можно обновить кеш в отдельной горутине.
				go InitCommandCache()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}

// StartWatchers запускает наблюдение за текущей директорией и за всеми директориями из PATH.
func StartWatchers() {
	var dirs []string

	// Текущая директория.
	curDir, err := os.Getwd()
	if err == nil {
		dirs = append(dirs, curDir)
	}
	// Директории из PATH.
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))
	dirs = append(dirs, paths...)

	// Запускаем наблюдение в отдельной горутине.
	go watchDirectories(dirs)
}
