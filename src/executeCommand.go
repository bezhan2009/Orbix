package src

import (
	"goCmd/cmd/commands/commandsWithSignaiture/template"
	"goCmd/cmd/commands/commandsWithoutSignature/Clean"
	"goCmd/cmd/commands/commandsWithoutSignature/Ls"
	"goCmd/cmd/commands/resourceIntensive/MatrixMultiplication"
	"goCmd/internal/Network"
	"goCmd/internal/Network/wifiUtils"
	ExCommUtils "goCmd/src/utils"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
)

func ExecuteCommand(executeCommand structs.ExecuteCommandFuncParams) {
	commandMap := map[string]func(){
		// Network related commands
		"wifiutils":  wifiUtils.Start,
		"pingview":   func() { Network.Ping(executeCommand.CommandArgs) },
		"traceroute": func() { Network.Traceroute(executeCommand.CommandArgs) },

		// File operations
		"extractzip": func() { ExCommUtils.ExtractZipUtil(executeCommand.CommandArgs) },
		"copysource": func() { ExCommUtils.CommandCopySourceUtil(executeCommand.CommandArgs) },
		"create":     func() { ExCommUtils.CreateFileUtil(executeCommand.CommandArgs, executeCommand.Dir) },
		"write":      func() { ExCommUtils.WriteFileUtil(executeCommand.CommandArgs) },
		"read":       func() { ExCommUtils.ReadFileUtil(executeCommand.CommandArgs) },
		"edit":       func() { ExCommUtils.EditFileUtil(executeCommand.CommandArgs) },
		"rename":     func() { ExCommUtils.RenameFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"ren":        func() { ExCommUtils.RenameFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"remove":     func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"del":        func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"rem":        func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"clean":      Clean.Screen,
		"cls":        Clean.Screen,
		"clear":      Clean.Screen,
		"cd":         func() { ExCommUtils.ChangeDirectoryUtil(executeCommand.CommandArgs) },

		// Utility commands
		"systemorbix": utils.SystemInformation,
		"open_link":   func() { ExCommUtils.OpenLinkUtil(executeCommand.CommandArgs) },

		// Calculation and resource intensive operations
		"matrixmul": func() { MatrixMultiplication.MatrixMulCommand(executeCommand.CommandArgs) },
		"primes":    func() { ExCommUtils.CalculatePrimesUtil(executeCommand.CommandArgs) },
		"picalc":    func() { ExCommUtils.CalculatePiUtil(executeCommand.CommandArgs) },
		"fileio":    func() { ExCommUtils.FileIOStressTestUtil(executeCommand.CommandArgs) },

		// Template and miscellaneous commands
		"newtemplate": func() { template.Make(executeCommand.CommandArgs) },
		"template":    func() { ExecuteShablonUtil(executeCommand.CommandArgs) },

		// Disk operations
		"cf": func() { ExCommUtils.CFUtil(executeCommand.CommandArgs) },
		"df": func() { ExCommUtils.DFUtil(executeCommand.CommandArgs) },

		// Listing and viewing
		"ls": Ls.PrintLS,
	}

	permissionRequiredCommands := map[string]func(){
		"orbix":   func() { Orbix("", true, structs.RebootedData{}) },
		"newuser": func() { NewUser(system.Path) },
		"signout": func() { SignOutUtil(executeCommand.Username, system.Path) },
		"exit": func() {
			*executeCommand.IsWorking = false
			removeUserFromRunningFile(executeCommand.Username)
		},
	}

	if handler, exists := commandMap[executeCommand.CommandLower]; exists {
		handler()
	} else if handler, exists := permissionRequiredCommands[executeCommand.CommandLower]; exists {
		if executeCommand.IsPermission {
			handler()
		}
	} else {
		HandleUnknownCommandUtil(executeCommand.CommandLower, executeCommand.CommandLine, executeCommand.Commands)
	}
}
