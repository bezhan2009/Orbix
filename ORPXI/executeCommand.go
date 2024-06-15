package ORPXI

import (
	"goCmd/Network"
	"goCmd/Network/wifiUtils"
	ExCommUtils "goCmd/ORPXI/utils"
	"goCmd/cmdPress"
	"goCmd/commands/commandsWithSignaiture/AddOwnCommand"
	"goCmd/commands/commandsWithSignaiture/Read"
	"goCmd/commands/commandsWithSignaiture/Write"
	"goCmd/commands/commandsWithSignaiture/shablon"
	"goCmd/commands/commandsWithoutSignature/Clean"
	"goCmd/commands/commandsWithoutSignature/Ls"
	"goCmd/commands/resourceIntensive/MatrixMultiplication"
	"goCmd/structs"
	"goCmd/utils"
)

func ExecuteCommand(commandLower, command, commandLine, dir string, commands []structs.Command, commandArgs []string, isWorking *bool, isPermission bool) {
	user := cmdPress.CmdUser(dir)
	switch commandLower {
	case "newcommand":
		AddOwnCommand.Start()

	case "wifiutils":
		wifiUtils.Start()

	case "pingview":
		Network.Ping(commandArgs)

	case "traceroute":
		Network.Traceroute(commandArgs)

	case "extractzip":
		ExCommUtils.ExtractZipUtil(commandArgs)

	case "scanport":
		ExCommUtils.ScanPortUtil(commandArgs)

	case "whois":
		ExCommUtils.WhoisUtil(commandArgs)

	case "dnslookup":
		ExCommUtils.DnsLookupUtil(commandArgs)

	case "ipinfo":
		ExCommUtils.IPInfoUtil(commandArgs)

	case "geoip":
		ExCommUtils.GeoIPUtil(commandArgs)

	case "orpxi":
		if isPermission {
			CMD("")
		}

	case "newuser":
		if isPermission {
			NewUser()
		}

	case "signout":
		if isPermission {
			SignOutUtil(user, isWorking)
		}

	case "matrixmul":
		MatrixMultiplication.MatrixMulCommand()

	case "primes":
		ExCommUtils.CalculatePrimesUtil(commandArgs)

	case "picalc":
		ExCommUtils.CalculatePiUtil(commandArgs)

	case "fileio":
		ExCommUtils.FileIOStressTestUtil(commandArgs)

	case "newshablon":
		shablon.Make()

	case "shablon":
		ExecuteShablonUtil(commandArgs)

	case "systemgocmd":
		utils.SystemInformation()

	case "exit":
		if isPermission {
			*isWorking = false
		}

	case "copysource":
		ExCommUtils.CommandCopySourceUtil(commandArgs)

	case "create":
		ExCommUtils.CreateFileUtil(commandArgs, command, user, dir)

	case "write":
		Write.File(commandLower, commandArgs, user, dir)

	case "read":
		Read.File(commandLower, commandArgs, user, dir)

	case "remove":
		ExCommUtils.RenameFileUtil(commandArgs, command, user, dir)

	case "rename":
		ExCommUtils.RenameFileUtil(commandArgs, command, user, dir)

	case "clean":
		Clean.Screen()

	case "cd":
		ExCommUtils.ChangeDirectoryUtil(commandArgs)

	case "edit":
		ExCommUtils.EditFileUtil(commandArgs)

	case "ls":
		Ls.PrintLS()

	default:
		HandleUnknownCommandUtil(commandLower, commandLine, commands)
	}
}
