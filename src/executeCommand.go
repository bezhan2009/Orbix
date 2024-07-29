package src

import (
	"fmt"
	"goCmd/cmd/cmdPress"
	"goCmd/cmd/commands/commandsWithSignaiture/Read"
	"goCmd/cmd/commands/commandsWithSignaiture/Write"
	"goCmd/cmd/commands/commandsWithSignaiture/shablon"
	"goCmd/cmd/commands/commandsWithoutSignature/Clean"
	"goCmd/cmd/commands/commandsWithoutSignature/Ls"
	"goCmd/cmd/commands/resourceIntensive/MatrixMultiplication"
	"goCmd/internal/Network"
	"goCmd/internal/Network/wifiUtils"
	ExCommUtils "goCmd/src/utils"
	"goCmd/structs"
	"goCmd/utils"
	"os"
	"os/exec"
)

func ExecuteCommand(commandLower, command, commandLine, dir string, commands []structs.Command, commandArgs []string, isWorking *bool, isPermission bool, username string) {
	user := cmdPress.CmdUser(dir)

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
		"newtemplate": shablon.Make,
		"template":    func() { ExecuteShablonUtil(commandArgs) },
		"systemorbix": utils.SystemInformation,
		"copysource":  func() { ExCommUtils.CommandCopySourceUtil(commandArgs) },
		"create":      func() { ExCommUtils.CreateFileUtil(commandArgs, command, user, dir) },
		"write":       func() { Write.File(commandLower, commandArgs, user, dir) },
		"read":        func() { Read.File(commandLower, commandArgs, user, dir) },
		"remove":      func() { ExCommUtils.RemoveFileUtil(commandArgs, command, user, dir) },
		"rename":      func() { ExCommUtils.RenameFileUtil(commandArgs, command, user, dir) },
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
		"newuser": NewUser,
		"signout": func() { SignOutUtil(username) },
		"exit": func() {
			removeFromRunningFile := true
			for _, arg := range commandArgs {
				if arg == "-t" {
					removeFromRunningFile = false
					cmd := exec.Command("py", "exit.py")
					if err := cmd.Start(); err != nil {
						fmt.Printf("Ошибка запуска exit.py: %v\n", err)
						return
					}
					cmd2 := exec.Command("exit")
					if err := cmd2.Start(); err != nil {
						fmt.Printf("Ошибка запуска команды exit: %v\n", err)
						return
					}
					os.Exit(cmd2.ProcessState.ExitCode())
					return
				}
			}
			if removeFromRunningFile {
				removeUserFromRunningFile(username)
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
