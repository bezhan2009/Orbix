package Orbix

import (
	"goCmd/cmd/commands"
	"goCmd/cmd/commands/resourceIntensive/MatrixMultiplication"
	"goCmd/cmd/commands/template"
	"goCmd/internal/Network"
	"goCmd/internal/Network/wifiUtils"
	"goCmd/src"
	"goCmd/src/environment"
	"goCmd/src/handlers"
	"goCmd/src/user"
	ExCommUtils "goCmd/src/utils"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
	"os"
	"strings"
)

func Command(executeCommand structs.ExecuteCommandFuncParams) {
	system.ExecutingCommand = true
	session := executeCommand.Session

	if src.CommandFile(strings.TrimSpace(executeCommand.CommandLower)) {
		src.FullFileName(&executeCommand.CommandArgs)
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
		"getvar":      func() { environment.GetVariableValueUtil(executeCommand) },
		"matrixmul":   func() { MatrixMultiplication.MatrixMulCommand(executeCommand.CommandArgs) },
		"primes":      func() { ExCommUtils.CalculatePrimesUtil(executeCommand.CommandArgs) },
		"picalc":      func() { ExCommUtils.CalculatePiUtil(executeCommand.CommandArgs) },
		"fileio":      func() { ExCommUtils.FileIOStressTestUtil(executeCommand.CommandArgs) },
		"newtemplate": func() { template.Make(executeCommand.CommandArgs) },
		"template":    func() { TemplateUtil(executeCommand.CommandArgs, executeCommand.SD) },
		"copysource":  func() { ExCommUtils.CommandCopySourceUtil(executeCommand.CommandArgs) },
		"create":      func() { ExCommUtils.CreateFileUtil(executeCommand.CommandArgs, system.UserDir) },
		"write":       func() { ExCommUtils.WriteFileUtil(executeCommand.CommandArgs) },
		"read":        func() { ExCommUtils.ReadFileUtil(executeCommand.CommandArgs) },
		"remove":      func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"del":         func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"delete":      func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"rem":         func() { ExCommUtils.RemoveFileUtil(executeCommand.CommandArgs, executeCommand.Command) },
		"rename": func() {
			ExCommUtils.RenameFileUtil(executeCommand.CommandArgs, executeCommand.Command, system.Yellow)
		},
		"cf": func() { ExCommUtils.CFUtil(executeCommand.CommandArgs) },
		"df": func() { ExCommUtils.DFUtil(executeCommand.CommandArgs) },
		"ren": func() {
			ExCommUtils.RenameFileUtil(executeCommand.CommandArgs, executeCommand.Command, system.Yellow)
		},
		"cd": func() {
			ExCommUtils.ChangeDirectoryUtil(executeCommand.CommandArgs, session)
			system.UserDir, _ = os.Getwd()
		},
		"edit":         func() { ExCommUtils.EditFileUtil(executeCommand.CommandArgs) },
		"open_link":    func() { ExCommUtils.OpenLinkUtil(executeCommand.CommandArgs) },
		"api_request":  func() { commands.ApiRequest() },
		"print":        func() { commands.Print(executeCommand.CommandArgs) },
		"kill":         func() { ExCommUtils.KillProcessUtil(executeCommand.CommandArgs) },
		"neofetch":     func() { ExCommUtils.NeofetchUtil(executeCommand, system.User, system.Commands) },
		"setvar":       func() { environment.SetVariableUtil(executeCommand.CommandArgs) },
		"stusenv":      func() { commands.SetUserFromENV(system.Path) },
		"set_user_env": func() { commands.SetUserFromENV(system.Path) },
		"new_prompt":   func() { session.IsAdmin = false },
		"old_prompt":   func() { session.IsAdmin = true },
		"delvar":       func() { environment.DeleteVariable(executeCommand.CommandArgs) },
		"new_window":   func() { src.OpenNewWindowForCommand(executeCommand) },
		"prompt":       func() { handlers.HandlePromptCommand(executeCommand.CommandArgs, executeCommand.Prompt) },

		"help":         handlers.DisplayHelp,
		"systemorbix":  environment.SystemInformation,
		"save":         environment.SaveVars,
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
			system.UserDir = dir

			Orbix("",
				true,
				structs.RebootedData{},
				executeCommand.SD)
			system.PreviousSessionPrefix = executeCommand.SessionPrefix
		},
		"newuser": func() { user.NewUser() },
		"signout": func() {
			SignOutUtil(executeCommand.Username, executeCommand.SD.Path, executeCommand.SD, executeCommand.SessionPrefix)
		},
		"exit": func() {
			*executeCommand.IsWorking = false
			*executeCommand.IsPermission = false
			user.DeleteUserFromRunningFile(executeCommand.Username)
			environment.SaveVars()
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
			isValid := utils.ValidCommand(executeCommand.CommandLower, system.AdditionalCommands)
			if !isValid {
				handlers.HandleUnknownCommandUtil(executeCommand.Command, system.Commands)
			}
		}
	}

	system.GlobalSession = *executeCommand.Session
	system.ExecutingCommand = false
}
