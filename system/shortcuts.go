package system

var ShortcutsJSONName = "shortcuts.json"

var Shortcuts = map[string]string{
	"dir":  "ls",
	"cd\\": "cd \\",
	"cd..": "cd ..",
}

var AvailableShortcuts = []string{"dir", "cd\\", "cd.."}
