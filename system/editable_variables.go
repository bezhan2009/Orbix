package system

var (
	Location = ""
	User     = ""
	Empty    = ""
	Prompt   = ""
)

// EditableVars Map, ограничивающий изменяемые переменные
var EditableVars = map[string]interface{}{
	"location": &Location,
	"prompt":   &Prompt,
	"user":     &User,
	"empty":    &Empty,
	"user_dir": &UserDir,
	"debug":    &Debug,
}

var CmdMap map[string]struct{} = BuildCommandMap(Commands)                     // for system.Command
var AdditionalCmdMap map[string]struct{} = BuildCommandMap(AdditionalCommands) // for system.Command

var AvailableEditableVars = []string{"location", "prompt", "user", "empty", "user_dir", "debug"}
var CustomEditableVars []string
