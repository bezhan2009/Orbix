package src

import (
	"goCmd/cmd/commands/commandsWithSignaiture/template"
	"goCmd/cmd/commands/commandsWithoutSignature/Clean"
	"goCmd/cmd/commands/commandsWithoutSignature/Ls"
	"goCmd/cmd/commands/resourceIntensive/MatrixMultiplication"
	"goCmd/cmd/dirInfo"
	"goCmd/internal/Network"
	"goCmd/internal/Network/wifiUtils"
	ExCommUtils "goCmd/src/utils"
	"goCmd/structs"
	"goCmd/system"
	"goCmd/utils"
)

func ExecuteCommand(commandLower, command, commandLine, dir string, commands []structs.Command, commandArgs []string, isWorking *bool, isPermission bool, username string) {
	user := dirInfo.CmdUser(dir)

	commandMap := map[string]func(){
		"wifiutils":   wifiUtils.Start,
		"pingview":    func() { Network.Ping(commandArgs) },
		"traceroute":  func() { Network.Traceroute(commandArgs) },
		"extractzip":  func() { ExCommUtils.ExtractZipUtil(commandArgs) },
		"scanport":    func() { ExCommUtils.ScanPortUtil(commandArgs) },
		"whois":       func() { ExCommUtils.WhoisUtil(commandArgs) },
		"dnslookup":   func() { ExCommUtils.DnsLookupUtil(commandArgs) },
		"ipinfo":      func() { ExCommUtils.IPInfoUtil(commandArgs) },
		"geoip":       func() { ExCommUtils.GeoIPUtil(commandArgs) },
		"matrixmul":   func() { MatrixMultiplication.MatrixMulCommand(commandArgs) },
		"primes":      func() { ExCommUtils.CalculatePrimesUtil(commandArgs) },
		"picalc":      func() { ExCommUtils.CalculatePiUtil(commandArgs) },
		"fileio":      func() { ExCommUtils.FileIOStressTestUtil(commandArgs) },
		"newtemplate": func() { template.Make(commandArgs) },
		"template":    func() { ExecuteShablonUtil(commandArgs) },
		"systemorbix": utils.SystemInformation,
		"copysource":  func() { ExCommUtils.CommandCopySourceUtil(commandArgs) },
		"create":      func() { ExCommUtils.CreateFileUtil(commandArgs, dir) },
		"write":       func() { ExCommUtils.WriteFileUtil(commandArgs) },
		"read":        func() { ExCommUtils.ReadFileUtil(commandLower, commandArgs, user, dir) },
		"remove":      func() { ExCommUtils.RemoveFileUtil(commandArgs, command) },
		"del":         func() { ExCommUtils.RemoveFileUtil(commandArgs, command) },
		"rem":         func() { ExCommUtils.RemoveFileUtil(commandArgs, command) },
		"rename":      func() { ExCommUtils.RenameFileUtil(commandArgs, command) },
		"cf":          func() { ExCommUtils.CFUtil(commandArgs) },
		"df":          func() { ExCommUtils.DFUtil(commandArgs) },
		"ren":         func() { ExCommUtils.RenameFileUtil(commandArgs, command) },
		"clean":       Clean.Screen,
		"cls":         Clean.Screen,
		"clear":       Clean.Screen,
		"cd":          func() { ExCommUtils.ChangeDirectoryUtil(commandArgs) },
		"edit":        func() { ExCommUtils.EditFileUtil(commandArgs) },
		"ls":          Ls.PrintLS,
		"open_link":   func() { ExCommUtils.OpenLinkUtil(commandArgs) },
	}

	permissionRequiredCommands := map[string]func(){
		"orbix":   func() { Orbix("", true) },
		"newuser": func() { NewUser(system.Path) },
		"signout": func() { SignOutUtil(username, system.Path) },
		"exit": func() {
			*isWorking = false
			removeUserFromRunningFile(username)
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
