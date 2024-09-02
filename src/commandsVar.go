package src

import (
	"github.com/fatih/color"
	"goCmd/structs"
)

var (
	red     func(a ...interface{}) string
	green   func(a ...interface{}) string
	yellow  func(a ...interface{}) string
	blue    func(a ...interface{}) string
	magenta func(a ...interface{}) string
	cyan    func(a ...interface{}) string
)

// Commands available Orbix commands
var Commands = []structs.Command{
	{"help", "List of available commands"},
	{"whois", "Domain information"},
	{"pingview", "Displays ping"},
	{"traceroute", "Route tracing"},
	{"extractzip", "Extracts .zip archives"},
	{"copysource", "Copies source code from file"},
	{"signout", "Signs out the current user"},
	{"newtemplate", "Creates a new command template for execution"},
	{"template", "Executes a specific command template"},
	{"newuser", "Adds a new user to Orbix"},
	{"prompt", "Changes the command prompt in Orbix"},
	{"systemorbix", "Displays system information about Orbix"},
	{"rename", "Renames a file or directory"},
	{"rem", "Deletes a file"},
	{"del", "Deletes a file"},
	{"cf", "Creates a new folder"},
	{"df", "Deletes a folder"},
	{"read", "Displays the contents of a file"},
	{"write", "Writes data to a file"},
	{"create", "Creates a new file"},
	{"exit", "Exits the program"},
	{"orbix", "Starts another instance of Orbix"},
	{"wifiutils", "Launches utility for WiFi operations"},
	{"clean", "Clears the terminal screen"},
	{"clear", "Clears the terminal screen"},
	{"cls", "Clears the terminal screen"},
	{"matrixmul", "Performs matrix multiplication"},
	{"primes", "Searches for prime numbers"},
	{"picalc", "Calculates the value of π"},
	{"fileio", "Performs file input/output intensive tests"},
	{"cd", "Changes the current directory"},
	{"edit", "Opens a file for editing"},
	{"ls", "Lists the contents of a directory"},
	{"scanport", "Scans open network ports"},
	{"dnslookup", "Performs DNS queries"},
	{"ipinfo", "Displays information about an IP address"},
	{"open_link", "Opens a URL in the default web browser"},
	{"geoip", "Displays geolocation information for an IP address"},
}

// AdditionalCommands additional commands
var AdditionalCommands = []structs.Command{
	{"help", "List of available commands"},
	{"whois", "Domain information"},
	{"pingview", "Displays ping"},
	{"traceroute", "Route tracing"},
	{"extractzip", "Extracts .zip archives"},
	{"copysource", "Copies source code from file"},
	{"signout", "Signs out the current user"},
	{"newtemplate", "Creates a new command template for execution"},
	{"template", "Executes a specific command template"},
	{"newuser", "Adds a new user to Orbix"},
	{"prompt", "Changes the command prompt in Orbix"},
	{"systemorbix", "Displays system information about Orbix"},
	{"rename", "Renames a file or directory"},
	{"rem", "Deletes a file"},
	{"del", "Deletes a file"},
	{"cf", "Creates a new folder"},
	{"df", "Deletes a folder"},
	{"read", "Displays the contents of a file"},
	{"write", "Writes data to a file"},
	{"create", "Creates a new file"},
	{"exit", "Exits the program"},
	{"orbix", "Starts another instance of Orbix"},
	{"wifiutils", "Launches utility for WiFi operations"},
	{"clean", "Clears the terminal screen"},
	{"clear", "Clears the terminal screen"},
	{"cls", "Clears the terminal screen"},
	{"matrixmul", "Performs matrix multiplication"},
	{"primes", "Searches for prime numbers"},
	{"picalc", "Calculates the value of π"},
	{"fileio", "Performs file input/output intensive tests"},
	{"cd", "Changes the current directory"},
	{"edit", "Opens a file for editing"},
	{"ls", "Lists the contents of a directory"},
	{"scanport", "Scans open network ports"},
	{"dnslookup", "Performs DNS queries"},
	{"ipinfo", "Displays information about an IP address"},
	{"open_link", "Opens a URL in the default web browser"},
	{"geoip", "Displays geolocation information for an IP address"},
	{"git", "Runs git commands"},
	{"calc", "Launches a calculator"},
	{"cmd", "Launches the command prompt"},
	{"go", "Runs Go language commands"},
	{"pip", "Runs Python package installer"},
	{"py", "Runs Python interpreter"},
	{"deploy", "Deploys the application"},
	{"upgrade", "Upgrades installed packages"},
	{"export", "Exports data to a file"},
	{"import", "Imports data from a file"},
	{"compress", "Compresses files into an archive"},
	{"decompress", "Decompresses files from an archive"},
	{"convert", "Converts files from one format to another"},
	{"monitor", "Monitors system resources"},
	{"network", "Displays network information and status"},
}

var CommandHistory []string

func Init() {
	// Initialize CommandHistory with package or tool names
	CommandHistory = append(CommandHistory, "help")
	CommandHistory = append(CommandHistory, "run")
	CommandHistory = append(CommandHistory, "push")
	CommandHistory = append(CommandHistory, "pull")
	CommandHistory = append(CommandHistory, "origin")
	CommandHistory = append(CommandHistory, "main")
	CommandHistory = append(CommandHistory, "master")
	CommandHistory = append(CommandHistory, "merge")
	CommandHistory = append(CommandHistory, ".")
	CommandHistory = append(CommandHistory, "remote")
	CommandHistory = append(CommandHistory, "add")
	CommandHistory = append(CommandHistory, "--version")
	CommandHistory = append(CommandHistory, "install")
	CommandHistory = append(CommandHistory, "django")
	CommandHistory = append(CommandHistory, "flask")
	CommandHistory = append(CommandHistory, "config")
	CommandHistory = append(CommandHistory, "--global")
	CommandHistory = append(CommandHistory, "-m")
	CommandHistory = append(CommandHistory, "-am")
	CommandHistory = append(CommandHistory, "--list")
	CommandHistory = append(CommandHistory, "config --global user.name \"Your name\"")
	CommandHistory = append(CommandHistory, "config --global user.email \"your_email@example.com\"")
	CommandHistory = append(CommandHistory, "branch")
	CommandHistory = append(CommandHistory, "checkout")
	CommandHistory = append(CommandHistory, "status")
	CommandHistory = append(CommandHistory, "commit")
	CommandHistory = append(CommandHistory, "clone")
	CommandHistory = append(CommandHistory, "log")
	CommandHistory = append(CommandHistory, "rebase")
	CommandHistory = append(CommandHistory, "cherry-pick")
	CommandHistory = append(CommandHistory, "stash")
	CommandHistory = append(CommandHistory, "reset")
	CommandHistory = append(CommandHistory, "diff")
	CommandHistory = append(CommandHistory, "grep")
	CommandHistory = append(CommandHistory, "fetch")
	CommandHistory = append(CommandHistory, "remote add")
	CommandHistory = append(CommandHistory, "remote remove")
	CommandHistory = append(CommandHistory, "tag")
	CommandHistory = append(CommandHistory, "show")
	CommandHistory = append(CommandHistory, "revert")
	CommandHistory = append(CommandHistory, "rm")
	CommandHistory = append(CommandHistory, "mv")
	CommandHistory = append(CommandHistory, "apply")
	CommandHistory = append(CommandHistory, "patch")
	CommandHistory = append(CommandHistory, "changelog")
	CommandHistory = append(CommandHistory, "upgrade")
	CommandHistory = append(CommandHistory, "export")
	CommandHistory = append(CommandHistory, "import")
	CommandHistory = append(CommandHistory, "tar compress")
	CommandHistory = append(CommandHistory, "tar decompress")
	CommandHistory = append(CommandHistory, "convert")
	CommandHistory = append(CommandHistory, "nmap monitor")

	// Initialize colors
	red = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan = color.New(color.FgCyan).SprintFunc()
	green = color.New(color.FgGreen).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	blue = color.New(color.FgBlue).SprintFunc()

	// Initialize git branch
	SetGitBranch()
}
