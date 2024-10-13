package structs

import "goCmd/system"

type ExecuteCommandFuncParams struct {
	CommandLower  string
	Command       string
	CommandLine   string
	Dir           string
	Commands      []Command
	CommandArgs   []string
	Username      string
	IsWorking     *bool
	IsPermission  *bool
	SD            *system.AppState
	SessionPrefix string
}
