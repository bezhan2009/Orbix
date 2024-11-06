package structs

import "goCmd/system"

type ProcessCommandParams struct {
	Command        string
	CommandInput   string
	CommandLower   string
	CommandLine    string
	CommandArgs    []string
	RunOnNewThread *bool
	EchoTime       *bool
	FirstCharIs    *bool
	LastCharIs     *bool
	IsWorking      *bool
	IsComHasFlag   *bool
	Session        *system.Session
	ExecCommand    ExecuteCommandFuncParams
}

type ExecuteCommandCatchErrs struct {
	CommandLower   string
	RunOnNewThread *bool
	EchoTime       *bool
}
