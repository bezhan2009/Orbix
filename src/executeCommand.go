package src

import (
	"goCmd/cmd/commands/commandsWithSignaiture/printCommand"
	"goCmd/cmd/commands/commandsWithSignaiture/template"
	"goCmd/cmd/commands/commandsWithoutSignature/Clean"
	"goCmd/cmd/commands/commandsWithoutSignature/Ls"
	redis_server "goCmd/cmd/commands/commandsWithoutSignature/redis-server"
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
		"wifiutils":   wifiUtils.Start,
		"pingview":    func() { Network.Ping(executeCommand.CommandArgs) },
		"traceroute":  func() { Network.Traceroute(executeCommand.CommandArgs) },
		"extractzip":  func() { ExCommUtils.ExtractZipUtil(executeCommand.CommandArgs) },
		"scanport":    func() { ExCommUtils.ScanPortUtil(executeCommand.CommandArgs) },
		"whois":       func() { ExCommUtils.WhoisUtil(executeCommand.CommandArgs) },
		"dnslookup":   func() { ExCommUtils.DnsLookupUtil(executeCommand.CommandArgs) },
		"ipinfo":      func() { ExCommUtils.IPInfoUtil(executeCommand.CommandArgs) },
		"geoip":       func() { ExCommUtils.GeoIPUtil(executeCommand.CommandArgs) },
		"matrixmul":   func() { MatrixMultiplication.MatrixMulCommand(executeCommand.CommandArgs) },
		"primes":      func() { ExCommUtils.CalculatePrimesUtil(executeCommand.CommandArgs) },
		"picalc":      func() { ExCommUtils.CalculatePiUtil(executeCommand.CommandArgs) },
		"fileio":      func() { ExCommUtils.FileIOStressTestUtil(executeCommand.CommandArgs) },
		"newtemplate": func() { template.Make(executeCommand.CommandArgs) },
		"template":    func() { ExecuteShablonUtil(executeCommand.CommandArgs) },
		"copysource":  func() { ExCommUtils.CommandCopySourceUtil(executeCommand.CommandArgs) },
		"create":      func() { ExCommUtils.CreateFileUtil(executeCommand.CommandArgs, executeCommand.Dir) },
		"write":       func() { ExCommUtils.WriteFileUtil(executeCommand.CommandArgs) },
		"read":        func() { ExCommUtils.ReadFileUtil(executeCommand.CommandArgs) },
		"remove":      func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"del":         func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"rem":         func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"rename":      func() { ExCommUtils.RenameFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"cf":          func() { ExCommUtils.CFUtil(executeCommand.CommandArgs) },
		"df":          func() { ExCommUtils.DFUtil(executeCommand.CommandArgs) },
		"ren":         func() { ExCommUtils.RenameFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"panic":       func() { ExCommUtils.PanicOnError(executeCommand) },
		"cd":          func() { ExCommUtils.ChangeDirectoryUtil(executeCommand.CommandArgs) },
		"edit":        func() { ExCommUtils.EditFileUtil(executeCommand.CommandArgs) },
		"open_link":   func() { ExCommUtils.OpenLinkUtil(executeCommand.CommandArgs) },
		"print":       func() { printCommand.Print(executeCommand.CommandArgs) },

		"systemorbix":  utils.SystemInformation,
		"clean":        Clean.Screen,
		"cls":          Clean.Screen,
		"clear":        Clean.Screen,
		"help":         displayHelp,
		"ls":           Ls.PrintLS,
		"redis":        redis_server.StartRedisServer,
		"redis-server": redis_server.StartRedisServer,
		"redisserver":  redis_server.StartRedisServer,
		"ubuntu_redis": redis_server.StartRedisServer,
		"redis_server": redis_server.StartRedisServer,
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
