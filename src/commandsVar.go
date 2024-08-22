package src

import "goCmd/structs"

// Commands available Orbix commands
var Commands = []structs.Command{
	{"help", "List of available Commands"},
	{"whois", "Domain information"},
	{"pingview", "Displays ping"},
	{"traceroute", "Route tracing"},
	{"extractzip", "Extracts .zip archives"},
	{"copysource", "from file copy his source"},
	{"signout", "User signs out of Orbix"},
	{"newtemplate", "Creates a new command template for execution"},
	{"template", "Executes a specific command template"},
	{"newuser", "New user for Orbix"},
	{"prompt", "Changes Orbix"},
	{"systemorbix", "Displays information about Orbix"},
	{"rename", "Renames a file"},
	{"ren", "Renames a file"},
	{"remove", "Deletes a file"},
	{"rem", "Deletes a file"},
	{"del", "Deletes a file"},
	{"cf", "Creates a folder"},
	{"df", "Deletes a folder"},
	{"read", "Displays file contents"},
	{"write", "Writes data to a file"},
	{"create", "Creates a new file"},
	{"exit", "Exits the program"},
	{"orbix", "Starts another orbix"},
	{"wifiutils", "Launches utility for WiFi operations"},
	{"clean", "Clears the screen"},
	{"clear", "Clears the screen"},
	{"cls", "Clears the screen"},
	{"matrixmul", "Multiplication of large matrices"},
	{"primes", "Search for large prime numbers"},
	{"picalc", "Calculates the value of π"},
	{"fileio", "File I/O intensive test"},
	{"cd", "Changes the current directory"},
	{"edit", "Edits a file"},
	{"ls", "Displays directory contents"},
	{"scanport", "Port scanning"},
	{"dnslookup", "DNS queries"},
	{"ipinfo", "IP address information"},
	{"open_link", "opens the link in the browser"},
	{"geoip", "IP address geolocation"},
	{"()", "Close parenthesis"},
	{"[]", "Close square bracket"},
	{"{}", "Close curly brace"},
	{"\"\"", "Close double quote"},
	{"n", ""},
	{"и", ""},
	{"т", ""},
	{"b", ""},
	{"/", ""},
}

var CommandHistory []string

func Init() {
	CommandHistory = append(CommandHistory, "help")
	CommandHistory = append(CommandHistory, "n")
	CommandHistory = append(CommandHistory, "/")
	CommandHistory = append(CommandHistory, "b")
	CommandHistory = append(CommandHistory, "т")
	CommandHistory = append(CommandHistory, "и")
}
