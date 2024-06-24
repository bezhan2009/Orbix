package src

import (
	"goCmd/Network"
	"goCmd/Network/wifiUtils"
	"goCmd/cmdPress"
	"goCmd/commands/commandsWithSignaiture/AddOwnCommand"
	"goCmd/commands/commandsWithSignaiture/Read"
	"goCmd/commands/commandsWithSignaiture/Write"
	"goCmd/commands/commandsWithSignaiture/shablon"
	"goCmd/commands/commandsWithoutSignature/Clean"
	"goCmd/commands/commandsWithoutSignature/Ls"
	"goCmd/commands/resourceIntensive/MatrixMultiplication"
	ExCommUtils "goCmd/src/utils"
	"goCmd/structs"
	"goCmd/utils"
	"os"
	"os/exec"
)

func ExecuteCommand(commandLower, command, commandLine, dir string, commands []structs.Command, commandArgs []string, isWorking *bool, isPermission bool) {
	user := cmdPress.CmdUser(dir)

	commandMap := map[string]func(){
		"newcommand":  AddOwnCommand.Start,
		"wifiutils":   wifiUtils.Start,
		"pingview":    func() { Network.Ping(commandArgs) },
		"traceroute":  func() { Network.Traceroute(commandArgs) },
		"extractzip":  func() { ExCommUtils.ExtractZipUtil(commandArgs) },
		"scanport":    func() { ExCommUtils.ScanPortUtil(commandArgs) },
		"whois":       func() { ExCommUtils.WhoisUtil(commandArgs) },
		"dnslookup":   func() { ExCommUtils.DnsLookupUtil(commandArgs) },
		"ipinfo":      func() { ExCommUtils.IPInfoUtil(commandArgs) },
		"geoip":       func() { ExCommUtils.GeoIPUtil(commandArgs) },
		"matrixmul":   MatrixMultiplication.MatrixMulCommand,
		"primes":      func() { ExCommUtils.CalculatePrimesUtil(commandArgs) },
		"picalc":      func() { ExCommUtils.CalculatePiUtil(commandArgs) },
		"fileio":      func() { ExCommUtils.FileIOStressTestUtil(commandArgs) },
		"newshablon":  shablon.Make,
		"shablon":     func() { ExecuteShablonUtil(commandArgs) },
		"systemgocmd": utils.SystemInformation,
		"copysource":  func() { ExCommUtils.CommandCopySourceUtil(commandArgs) },
		"create":      func() { ExCommUtils.CreateFileUtil(commandArgs, command, user, dir) },
		"write":       func() { Write.File(commandLower, commandArgs, user, dir) },
		"read":        func() { Read.File(commandLower, commandArgs, user, dir) },
		"remove":      func() { ExCommUtils.RemoveFileUtil(commandArgs, command, user, dir) },
		"rename":      func() { ExCommUtils.RenameFileUtil(commandArgs, command, user, dir) },
		"clean":       Clean.Screen,
		"cd":          func() { ExCommUtils.ChangeDirectoryUtil(commandArgs) },
		"edit":        func() { ExCommUtils.EditFileUtil(commandArgs) },
		"ls":          Ls.PrintLS,
		"open_link":   func() { ExCommUtils.OpenLinkUtil(commandArgs) },
	}

	permissionRequiredCommands := map[string]func(){
		"orbix":   func() { Orbix("") },
		"newuser": NewUser,
		"signout": func() { SignOutUtil(user, isWorking) },
		"exit": func() {
			// Check if '-t' argument is present
			for _, arg := range commandArgs {
				if arg == "-t" {
					cmd := exec.Command("py", "exit.py")
					cmd.Start()
					cmd2 := exec.Command("exit")
					cmd2.Start()
					os.Exit(cmd2.ProcessState.ExitCode())
					return
				}
			}
			*isWorking = false
		},
	}

	if handler, exists := commandMap[commandLower]; exists {
		handler()
	} else if handler, exists := permissionRequiredCommands[commandLower]; exists {
		if isPermission {
			handler()
		}
	} else {
		HandleUnknownCommandUtil(commandLower, commandLine, commands)
	}
}
