package src

import (
	"fmt"
	"goCmd/cmd/commands"
	"goCmd/cmd/commands/resourceIntensive/MatrixMultiplication"
	"goCmd/cmd/commands/template"
	"goCmd/internal/Network"
	"goCmd/internal/Network/wifiUtils"
	ExCommUtils "goCmd/src/utils"
	"goCmd/structs"
	"goCmd/system"
	"os"
	"strings"
)

func ExecuteCommand(executeCommand structs.ExecuteCommandFuncParams) {
	ExecutingCommand = true
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
		"getvar":      func() { GetVariableValueUtil(executeCommand) },
		"matrixmul":   func() { MatrixMultiplication.MatrixMulCommand(executeCommand.CommandArgs) },
		"primes":      func() { ExCommUtils.CalculatePrimesUtil(executeCommand.CommandArgs) },
		"picalc":      func() { ExCommUtils.CalculatePiUtil(executeCommand.CommandArgs) },
		"fileio":      func() { ExCommUtils.FileIOStressTestUtil(executeCommand.CommandArgs) },
		"newtemplate": func() { template.Make(executeCommand.CommandArgs) },
		"template":    func() { ExecuteTemplateUtil(executeCommand.CommandArgs, executeCommand.SD) },
		"copysource":  func() { ExCommUtils.CommandCopySourceUtil(executeCommand.CommandArgs) },
		"create":      func() { ExCommUtils.CreateFileUtil(executeCommand.CommandArgs, executeCommand.Dir) },
		"write":       func() { ExCommUtils.WriteFileUtil(executeCommand.CommandArgs) },
		"read":        func() { ExCommUtils.ReadFileUtil(executeCommand.CommandArgs) },
		"remove":      func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"del":         func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"delete":      func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"rem":         func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"rename":      func() { ExCommUtils.RenameFileUtil(executeCommand.CommandArgs, executeCommand.Command, yellow) },
		"cf":          func() { ExCommUtils.CFUtil(executeCommand.CommandArgs) },
		"df":          func() { ExCommUtils.DFUtil(executeCommand.CommandArgs) },
		"ren":         func() { ExCommUtils.RenameFileUtil(executeCommand.CommandArgs, executeCommand.Command, yellow) },
		"cd": func() {
			ExCommUtils.ChangeDirectoryUtil(executeCommand.CommandArgs, session)
			SetGitBranch(session)
		},
		"edit":         func() { ExCommUtils.EditFileUtil(executeCommand.CommandArgs) },
		"open_link":    func() { ExCommUtils.OpenLinkUtil(executeCommand.CommandArgs) },
		"print":        func() { commands.Print(executeCommand.CommandArgs) },
		"neofetch":     func() { ExCommUtils.NeofetchUtil(executeCommand, session, Commands) },
		"setvar":       func() { SetVariableUtil(executeCommand.CommandArgs) },
		"stusenv":      func() { commands.SetUserFromENV(system.Path) },
		"set_user_env": func() { commands.SetUserFromENV(system.Path) },

		"help":         displayHelp,
		"systemorbix":  SystemInformation,
		"clean":        commands.Screen,
		"cls":          commands.Screen,
		"clear":        commands.Screen,
		"ls":           commands.PrintLS,
		"redis":        commands.StartRedisServer,
		"redis-server": commands.StartRedisServer,
		"redisserver":  commands.StartRedisServer,
		"ubuntu_redis": commands.StartRedisServer,
		"redis_server": commands.StartRedisServer,
	}

	permissionRequiredCommands := map[string]func(){
		"panic": func() { ExCommUtils.PanicOnError(executeCommand) },
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
			*executeCommand.IsPermission = false
			RemoveUserFromRunningFile(executeCommand.Username)
		},
	}

	if *executeCommand.IsWorking {
		if strings.TrimSpace(executeCommand.CommandInput) != "" {
			*executeCommand.IsWorking = false
		}

		if handler, exists := commandMap[executeCommand.CommandLower]; exists {
			handler()
		} else if handler, exists = permissionRequiredCommands[executeCommand.CommandLower]; exists {
			if *executeCommand.IsPermission {
				handler()
			}
		} else {
			HandleUnknownCommandUtil(executeCommand.CommandLower, executeCommand.Commands)
		}
	}

	ExecutingCommand = false
}
