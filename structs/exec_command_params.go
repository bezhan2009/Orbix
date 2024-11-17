package structs

import "goCmd/system"

type ExecuteCommandFuncParams struct {
	Prompt        *string
	CommandLower  string
	Command       string
	CommandLine   string
	CommandInput  string
	Commands      []system.Command
	CommandArgs   []string
	Username      string
	IsWorking     *bool
	IsPermission  *bool
	SD            *system.AppState
	SessionPrefix string
	Session       *system.Session
}
