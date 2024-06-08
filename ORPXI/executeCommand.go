package ORPXI

import (
	"fmt"
	"goCmd/Network"
	"goCmd/Network/wifiUtils"
	"goCmd/cmdPress"
	"goCmd/commands/commandsWithSignaiture/Create"
	"goCmd/commands/commandsWithSignaiture/Edit"
	"goCmd/commands/commandsWithSignaiture/ExtractZip"
	"goCmd/commands/commandsWithSignaiture/Read"
	"goCmd/commands/commandsWithSignaiture/Remove"
	"goCmd/commands/commandsWithSignaiture/Rename"
	"goCmd/commands/commandsWithSignaiture/Write"
	"goCmd/commands/commandsWithSignaiture/shablon"
	"goCmd/commands/commandsWithoutSignature/CD"
	"goCmd/commands/commandsWithoutSignature/Clean"
	"goCmd/commands/commandsWithoutSignature/Ls"
	"goCmd/commands/resourceIntensive/FileIOStressTest"
	"goCmd/commands/resourceIntensive/MatrixMultiplication"
	"goCmd/commands/resourceIntensive/PiCalculation"
	"goCmd/commands/resourceIntensive/PrimeNumbers"
	"goCmd/debug"
	"goCmd/structs"
	"goCmd/utils"
	"os"
	"path/filepath"
	"strconv"
)

func ExecuteCommand(commandLower string, command string, commandLine string, dir string, commands []structs.Command, commandArgs []string, isWorking *bool, isPermission bool) {
	user := cmdPress.CmdUser(dir)
	switch commandLower {
	case "wifiutils":
		wifiUtils.Start()
	case "pingview":
		Network.Ping(commandArgs)
	case "traceroute":
		Network.Traceroute(commandArgs)
	case "extractzip":
		if len(commandArgs) < 2 {
			fmt.Println("Usage: extractzip <zipfile> <destination>")
		} else {
			err := ExtractZip.ExtractZip(commandArgs[0], commandArgs[1])
			if err != nil {
				fmt.Println("Error extracting ZIP file:", err)
			}
		}
	case "scanport":
		if len(commandArgs) < 2 {
			fmt.Println("Usage: scanport <host> <ports>")
		} else {
			ports := []int{}
			for _, p := range commandArgs[1:] {
				port, err := strconv.Atoi(p)
				if err != nil {
					fmt.Printf("Invalid port: %s\n", p)
					return
				}
				ports = append(ports, port)
			}
			Network.ScanPort(commandArgs[0], ports)
		}
	case "whois":
		if len(commandArgs) < 1 {
			fmt.Println("Usage: whois <domain>")
		} else {
			Network.Whois(commandArgs[0])
		}
	case "dnslookup":
		if len(commandArgs) < 1 {
			fmt.Println("Usage: dnslookup <domain>")
		} else {
			Network.DNSLookup(commandArgs[0])
		}
	case "ipinfo":
		if len(commandArgs) < 1 {
			fmt.Println("Usage: ipinfo <ip>")
		} else {
			Network.IPInfo(commandArgs[0])
		}
	case "geoip":
		if len(commandArgs) < 1 {
			fmt.Println("Usage: geoip <ip>")
		} else {
			Network.GeoIP(commandArgs[0])
		}
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
			if !CheckUser(user) {
				*isWorking = false
			}
		}
	case "matrixmul":
		MatrixMultiplication.MatrixMulCommand()
	case "primes":
		limit := 100000
		if len(commandArgs) > 0 {
			l, err := strconv.Atoi(commandArgs[0])
			if err == nil {
				limit = l
			}
		}
		PrimeNumbers.PrimeCommand(limit)
	case "picalc":
		precision := 10000
		if len(commandArgs) > 0 {
			p, err := strconv.Atoi(commandArgs[0])
			if err == nil {
				precision = p
			}
		}
		PiCalculation.PiCalcCommand(precision)
	case "fileio":
		filename := "largefile.dat"
		size := 100 * 1024 * 1024
		if len(commandArgs) > 0 {
			s, err := strconv.Atoi(commandArgs[0])
			if err == nil {
				size = s
			}
		}
		FileIOStressTest.FileIOCommand(filename, size)
	case "newshablon":
		shablon.Make()
	case "shablon":
		if len(commandArgs) < 1 {
			fmt.Println("Использование: shablon <название_шаблона>")
			return
		}

		nameShablon := commandArgs[0]
		err := Start(nameShablon)
		if err != nil {
			fmt.Println(err)
		}
	case "systemgocmd":
		utils.SystemInformation()
	case "exit":
		if isPermission {
			*isWorking = false
		}
	case "create":
		name, err := Create.File(commandArgs)
		if err != nil {
			fmt.Println(err)
			debug.Commands(command, false)
		} else if name != "" {
			fmt.Printf("Файл %s успешно создан!\n", name)
			fmt.Printf("Директория нового файла: %s\n", filepath.Join(dir, name))
			debug.Commands(command, true)
		}
	case "write":
		Write.File(commandLower, commandArgs)
	case "read":
		Read.File(commandLower, commandArgs)
	case "remove":
		name, err := Remove.File(commandArgs)
		if err != nil {
			debug.Commands(command, false)
			fmt.Println(err)
		} else {
			debug.Commands(command, true)
			fmt.Printf("Файл %s успешно удален!\n", name)
		}
	case "rename":
		errRename := Rename.Rename(commandArgs)
		if errRename != nil {
			debug.Commands(command, false)
			fmt.Println(errRename)
		} else {
			debug.Commands(command, true)
		}
	case "clean":
		Clean.Screen()
	case "cd":
		if len(commandArgs) == 0 {
			dir, _ := os.Getwd()
			fmt.Println(dir)
		} else {
			err := CD.ChangeDirectory(commandArgs[0])
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	case "edit":
		if len(commandArgs) < 1 {
			fmt.Println("Использование: edit <файл>")
			return
		}
		filename := commandArgs[0]
		err := Edit.File(filename)
		if err != nil {
			fmt.Println(err)
		}
	case "ls":
		Ls.PrintLS()
	default:
		validCommand := utils.ValidCommand(commandLower, commands)
		if !validCommand {
			fmt.Printf("'%s' не является внутренней или внешней командой,\nисполняемой программой или пакетным файлом.\n", commandLine)
			suggestedCommand := suggestCommand(commandLower)
			if suggestedCommand != "" {
				fmt.Printf("Возможно, вы имели в виду: %s?\n", suggestedCommand)
			}
		}
	}
}
