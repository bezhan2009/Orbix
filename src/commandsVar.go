package src

import "goCmd/structs"

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

	// new commands
	{"branch", "Manages git branches"},
	{"checkout", "Switches between git branches"},
	{"status", "Shows the working tree status"},
	{"commit", "Records changes to the repository"},
	{"clone", "Clones a repository into a new directory"},
	{"log", "Shows commit logs"},
	{"rebase", "Reapply commits on top of another base tip"},
	{"cherry-pick", "Apply the changes introduced by some existing commits"},
	{"stash", "Stashes the changes in a dirty working directory away"},
	{"reset", "Resets the current HEAD to the specified state"},
	{"diff", "Shows changes between commits, commit and working tree, etc."},
	{"grep", "Prints lines matching a pattern"},
	{"fetch", "Downloads objects and refs from another repository"},
	{"remote", "Manages set of tracked repositories"},
	{"tag", "Lists, creates, deletes, or verifies tags object in repository"},
	{"show", "Displays various types of objects"},
	{"revert", "Reverts some existing commits"},
	{"rm", "Removes files from the working tree and from the index"},
	{"mv", "Moves or renames a file, a directory, or a symlink"},
	{"apply", "Applies a patch to files and/or to the index"},
	{"patch", "Creates patches from changes"},
	{"changelog", "Generates a changelog from git history"},
	{"build", "Builds the project from source"},
	{"test", "Runs tests for the project"},
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
	CommandHistory = append(CommandHistory, "--list")
	CommandHistory = append(CommandHistory, "user.name \"Your name\"")
	CommandHistory = append(CommandHistory, "user.email \"your_email@example.com\"")
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
	CommandHistory = append(CommandHistory, "build")
	CommandHistory = append(CommandHistory, "test")
	CommandHistory = append(CommandHistory, "deploy")
	CommandHistory = append(CommandHistory, "upgrade")
	CommandHistory = append(CommandHistory, "export")
	CommandHistory = append(CommandHistory, "import")
	CommandHistory = append(CommandHistory, "compress")
	CommandHistory = append(CommandHistory, "decompress")
	CommandHistory = append(CommandHistory, "convert")
	CommandHistory = append(CommandHistory, "monitor")
	CommandHistory = append(CommandHistory, "network")
}
