package structs

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
	IsComHasFlag   *bool
	ExecCommand    ExecuteCommandFuncParams
	LoopData       OrbixLoopData
}

type ExecuteCommandCatchErrs struct {
	CommandLower   string
	RunOnNewThread *bool
	EchoTime       *bool
}
