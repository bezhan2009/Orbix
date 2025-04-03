package Orbix

import (
	"fmt"
	_chan "goCmd/chan"
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
	"os"
	"strings"
)

func Command(executeCommand *structs.ExecuteCommandFuncParams) {
	system.ExecutingCommand = true
	session := executeCommand.LoopData.Session

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
		"template":    func() { TemplateUtil(executeCommand.CommandArgs, executeCommand.LoopData.SessionData) },
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
		"neofetch":     func() { ExCommUtils.NeofetchUtil(executeCommand, system.User, system.CmdMap) },
		"setvar":       func() { environment.SetVariableUtil(executeCommand.CommandArgs) },
		"shortcut":     func() { environment.SetShortCutUtil(executeCommand.CommandArgs) },
		"delvar":       func() { environment.DeleteVariable(executeCommand.CommandArgs) },
		"delshort":     func() { environment.DeleteShortcut(executeCommand.CommandArgs) },
		"getshort":     func() { environment.GetShortcutValueUtil(executeCommand) },
		"stusenv":      func() { commands.SetUserFromENV(system.Path) },
		"set_user_env": func() { commands.SetUserFromENV(system.Path) },
		"new_prompt":   func() { session.IsAdmin = false },
		"old_prompt":   func() { session.IsAdmin = true },
		"new_window":   func() { src.OpenNewWindowForCommand(executeCommand) },
		"prompt":       func() { handlers.HandlePromptCommand(executeCommand.CommandArgs, executeCommand.Prompt) },
		"getenv":       func() { ExCommUtils.GetEnvVarUtil(executeCommand.CommandArgs) },
		"setenv":       func() { commands.SetEnvVar(executeCommand.CommandArgs) },
		"chport":       func() { commands.IsPortOpen(executeCommand.CommandArgs) },
		"fileinfo":     func() { commands.GetFileInfo(executeCommand.CommandArgs) },

		"help":        handlers.DisplayHelp,
		"systemorbix": environment.SystemInformation,
		"save":        environment.SaveVars,
		"load": func() {
			err := environment.LoadUserConfigs()
			if err != nil {
				fmt.Println(system.Red("Error loading user configs:\n" + err.Error()))
			}
		},
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

			src.PrepareOrbix()

			Orbix("",
				true,
				structs.RebootedData{},
				executeCommand.LoopData.SessionData)
			system.PreviousSessionPrefix = executeCommand.SessionPrefix

			src.RestoreOrbix()
		},
		"newuser": func() { user.NewUser() },
		"signout": func() {
			SignOutUtil(
				executeCommand.LoopData.Username,
				executeCommand.LoopData.SessionData.Path,
				executeCommand.LoopData.SessionData,
				executeCommand.SessionPrefix,
			)
		},
		"exit": func() {
			*executeCommand.LoopData.IsWorking = false
			*executeCommand.LoopData.IsPermission = false

			user.DeleteUserFromRunningFile(executeCommand.LoopData.Username)

			if _chan.LoadConfigsFn() != nil &&
				session.IsAdmin &&
				!system.Unauthorized &&
				system.CntLaunchedOrbixes > 1 {
				_ = _chan.LoadConfigsFn()
			}

			if len(executeCommand.CommandArgs) > 0 && executeCommand.CommandArgs[0] == "\\k" {
				pid := os.Getpid()
				process, err := os.FindProcess(pid)
				if err != nil {
					fmt.Println(system.Red("Error finding Orbix process:", err))
					return
				}

				err = process.Kill()
				if err != nil {
					fmt.Println(system.Red("Error killing Orbix:", err))
					return
				}
			}
		},
	}
	defer func() {
		commandMap = nil
		permissionRequiredCommands = nil
	}()

	if *executeCommand.LoopData.IsWorking {
		if strings.TrimSpace(executeCommand.CommandInput) != "" {
			*executeCommand.LoopData.IsWorking = false
		}

		if handler, exists := commandMap[executeCommand.CommandLower]; exists {
			handler()
		} else if handler, exists = permissionRequiredCommands[executeCommand.CommandLower]; exists {
			if *executeCommand.LoopData.IsPermission {
				handler()
			}
		} else {
			handlers.HandleUnknownCommandUtil(executeCommand.Command, executeCommand.CommandLower, system.CmdMap)
		}
	}

	system.GlobalSession = *executeCommand.LoopData.Session
	system.ExecutingCommand = false
}
