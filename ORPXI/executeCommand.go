package ORPXI

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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
)

func ExecuteCommand(commandLower, command, commandLine, dir string, commands []structs.Command, commandArgs []string, isWorking *bool, isPermission bool) {
	user := cmdPress.CmdUser(dir)
	switch commandLower {
	case "wifiutils":
		wifiUtils.Start()

	case "pingview":
		Network.Ping(commandArgs)

	case "traceroute":
		Network.Traceroute(commandArgs)

	case "extractzip":
		extractZip(commandArgs)

	case "scanport":
		scanPort(commandArgs)

	case "whois":
		whois(commandArgs)

	case "dnslookup":
		dnsLookup(commandArgs)

	case "ipinfo":
		ipInfo(commandArgs)

	case "geoip":
		geoIP(commandArgs)

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
			signOut(user, isWorking)
		}

	case "matrixmul":
		MatrixMultiplication.MatrixMulCommand()

	case "primes":
		calculatePrimes(commandArgs)

	case "picalc":
		calculatePi(commandArgs)

	case "fileio":
		fileIOStressTest(commandArgs)

	case "newshablon":
		shablon.Make()

	case "shablon":
		executeShablon(commandArgs)

	case "systemgocmd":
		utils.SystemInformation()

	case "exit":
		if isPermission {
			*isWorking = false
		}

	case "create":
		createFile(commandArgs, command, user, dir)

	case "write":
		Write.File(commandLower, commandArgs, user, dir)

	case "read":
		Read.File(commandLower, commandArgs, user, dir)

	case "remove":
		removeFile(commandArgs, command, user, dir)

	case "rename":
		renameFile(commandArgs, command, user, dir)

	case "clean":
		Clean.Screen()

	case "cd":
		changeDirectory(commandArgs)

	case "edit":
		editFile(commandArgs)

	case "ls":
		Ls.PrintLS()

	default:
		handleUnknownCommand(commandLower, commandLine, commands)
	}
}

func extractZip(commandArgs []string) {
	if len(commandArgs) < 2 {
		fmt.Println("Usage: extractzip <zipfile> <destination>")
		return
	}
	if err := ExtractZip.ExtractZip(commandArgs[0], commandArgs[1]); err != nil {
		fmt.Println("Error extracting ZIP file:", err)
	}
}

func scanPort(commandArgs []string) {
	if len(commandArgs) < 2 {
		fmt.Println("Usage: scanport <host> <ports>")
		return
	}
	var ports []int
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

func whois(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: whois <domain>")
		return
	}
	Network.Whois(commandArgs[0])
}

func dnsLookup(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: dnslookup <domain>")
		return
	}
	Network.DNSLookup(commandArgs[0])
}

func ipInfo(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: ipinfo <ip>")
		return
	}
	Network.IPInfo(commandArgs[0])
}

func geoIP(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: geoip <ip>")
		return
	}
	Network.GeoIP(commandArgs[0])
}

func signOut(user string, isWorking *bool) {
	if !CheckUser(user) {
		*isWorking = false
	}
}

func calculatePrimes(commandArgs []string) {
	limit := 100000
	if len(commandArgs) > 0 {
		if l, err := strconv.Atoi(commandArgs[0]); err == nil {
			limit = l
		}
	}
	PrimeNumbers.PrimeCommand(limit)
}

func calculatePi(commandArgs []string) {
	precision := 10000
	if len(commandArgs) > 0 {
		if p, err := strconv.Atoi(commandArgs[0]); err == nil {
			precision = p
		}
	}
	PiCalculation.PiCalcCommand(precision)
}

func fileIOStressTest(commandArgs []string) {
	filename := "largefile.dat"
	size := 100 * 1024 * 1024
	if len(commandArgs) > 0 {
		if s, err := strconv.Atoi(commandArgs[0]); err == nil {
			size = s
		}
	}
	FileIOStressTest.FileIOCommand(filename, size)
}

func executeShablon(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: shablon <template_name>")
		return
	}
	if err := Start(commandArgs[0]); err != nil {
		fmt.Println(err)
	}
}

func createFile(commandArgs []string, command, user, dir string) {
	name, err := Create.File(commandArgs)
	if err != nil {
		fmt.Println(err)
		debug.Commands(command, false, commandArgs, user, dir)
	} else if name != "" {
		fmt.Printf("File %s successfully created!\n", name)
		fmt.Printf("Directory of the new file: %s\n", filepath.Join(dir, name))
		debug.Commands(command, true, commandArgs, user, dir)
	}
}

func removeFile(commandArgs []string, command string, user, dir string) {
	name, err := Remove.File(commandArgs)
	if err != nil {
		debug.Commands(command, false, commandArgs, user, dir)
		fmt.Println(err)
	} else {
		debug.Commands(command, true, commandArgs, user, dir)
		fmt.Printf("File %s successfully deleted!\n", name)
	}
}

func renameFile(commandArgs []string, command string, user, dir string) {
	if err := Rename.Rename(commandArgs); err != nil {
		debug.Commands(command, false, commandArgs, user, dir)
		fmt.Println(err)
	} else {
		debug.Commands(command, true, commandArgs, user, dir)
	}
}

func changeDirectory(commandArgs []string) {
	if len(commandArgs) == 0 {
		dir, _ := os.Getwd()
		fmt.Println(dir)
		return
	}
	if err := CD.ChangeDirectory(commandArgs[0]); err != nil {
		fmt.Println(err)
	}
}

func editFile(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: edit <file>")
		return
	}
	if err := Edit.File(commandArgs[0]); err != nil {
		fmt.Println(err)
	}
}

func handleUnknownCommand(commandLower, commandLine string, commands []structs.Command) {
	if !utils.ValidCommand(commandLower, commands) {
		fmt.Printf("'%s' is not recognized as an internal or external command,\noperable program or batch file.\n", commandLine)
		if suggestedCommand := suggestCommand(commandLower); suggestedCommand != "" {
			fmt.Printf("Did you mean: %s?\n", suggestedCommand)
		}
	}
}
