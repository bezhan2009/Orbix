package src

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/printCommand"
	"goCmd/cmd/commands/commandsWithSignaiture/template"
	"goCmd/cmd/commands/commandsWithoutSignature/Clean"
	"goCmd/cmd/commands/commandsWithoutSignature/Ls"
	redisserver "goCmd/cmd/commands/commandsWithoutSignature/redis-server"
	"goCmd/cmd/commands/resourceIntensive/MatrixMultiplication"
	"goCmd/internal/Network"
	"goCmd/internal/Network/wifiUtils"
	ExCommUtils "goCmd/src/utils"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
	"os"
)

func ExecuteCommand(executeCommand structs.ExecuteCommandFuncParams) {
	session, exists := executeCommand.SD.GetSession(executeCommand.SessionPrefix)
	if !exists {
		fmt.Println(red("Session Not Found!!!"))
		*executeCommand.IsWorking = false
	}

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
		"template":    func() { ExecuteShablonUtil(executeCommand.CommandArgs, executeCommand.SD) },
		"copysource":  func() { ExCommUtils.CommandCopySourceUtil(executeCommand.CommandArgs) },
		"create":      func() { ExCommUtils.CreateFileUtil(executeCommand.CommandArgs, executeCommand.Dir) },
		"write":       func() { ExCommUtils.WriteFileUtil(executeCommand.CommandArgs) },
		"read":        func() { ExCommUtils.ReadFileUtil(executeCommand.CommandArgs) },
		"remove":      func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"del":         func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"rem":         func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"rename":      func() { ExCommUtils.RenameFileUtil(executeCommand.CommandArgs, executeCommand.Command, yellow) },
		"cf":          func() { ExCommUtils.CFUtil(executeCommand.CommandArgs) },
		"df":          func() { ExCommUtils.DFUtil(executeCommand.CommandArgs) },
		"ren":         func() { ExCommUtils.RenameFileUtil(executeCommand.CommandArgs, executeCommand.Command, yellow) },
		"panic":       func() { ExCommUtils.PanicOnError(executeCommand) },
		"cd":          func() { ExCommUtils.ChangeDirectoryUtil(executeCommand.CommandArgs, session) },
		"edit":        func() { ExCommUtils.EditFileUtil(executeCommand.CommandArgs) },
		"open_link":   func() { ExCommUtils.OpenLinkUtil(executeCommand.CommandArgs) },
		"print":       func() { printCommand.Print(executeCommand.CommandArgs) },

		"systemorbix":  utils.SystemInformation,
		"neofetch":     func() { ExCommUtils.NeofetchUtil(executeCommand, session, Commands) },
		"clean":        Clean.Screen,
		"cls":          Clean.Screen,
		"clear":        Clean.Screen,
		"help":         displayHelp,
		"ls":           Ls.PrintLS,
		"redis":        redisserver.StartRedisServer,
		"redis-server": redisserver.StartRedisServer,
		"redisserver":  redisserver.StartRedisServer,
		"ubuntu_redis": redisserver.StartRedisServer,
		"redis_server": redisserver.StartRedisServer,
	}

	permissionRequiredCommands := map[string]func(){
		"orbix": func() {
			dir, _ := os.Getwd()

			session.Path = dir
			Orbix("", true, structs.RebootedData{}, executeCommand.SD)
			PreviousSessionPrefix = executeCommand.SessionPrefix
		},
		"newuser": func() { NewUser(system.Path) },
		"signout": func() {
			SignOutUtil(executeCommand.Username, executeCommand.SD.Path, executeCommand.SD, executeCommand.SessionPrefix)
		},
		"exit": func() {
			*executeCommand.IsWorking = false
			removeUserFromRunningFile(executeCommand.Username)
		},
	}

	if handler, exists := commandMap[executeCommand.CommandLower]; exists {
		handler()
	} else if handler, exists = permissionRequiredCommands[executeCommand.CommandLower]; exists {
		if executeCommand.IsPermission {
			handler()
		}
	} else {
		HandleUnknownCommandUtil(executeCommand.CommandLower, executeCommand.CommandLine, executeCommand.Commands)
	}
}
