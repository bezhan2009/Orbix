package structs

import (
	"goCmd/system"
)

type OrbixLoopData struct {
	Username     string
	CommandInput string

	IsWorking        *bool
	IsPermission     *bool
	RestartAfterInit *bool

	SessionData *system.AppState
	Session     *system.Session
}

type ExecuteCommandFuncParams struct {
	Prompt        *string
	CommandLower  string
	Command       string
	CommandLine   string
	CommandInput  string
	Commands      []system.Command
	CommandArgs   []string
	SessionPrefix string
	LoopData      OrbixLoopData
}
