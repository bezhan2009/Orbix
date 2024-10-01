package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"goCmd/system"
	"io/ioutil"
	"log"
	"os"

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

// Config represents the structure of the JSON configuration file
type Config struct {
	Colors             map[string]string `json:"colors"`
	Commands           []structs.Command `json:"commands"`
	AdditionalCommands []structs.Command `json:"additionalCommands"`
	CommandHistory     []string          `json:"commandHistory"`
}

//var Commands []structs.Command
//var AdditionalCommands []structs.Command
//var CommandHistory []string

func SetCommands() {
	// Load config from JSON
	config, err := ReadConfigs("commands.json")
	if err != nil {
		panic(err)
		return
	}

	// Initialize colors
	red = color.New(getColor(config.Colors["red"])).SprintFunc()
	yellow = color.New(getColor(config.Colors["yellow"])).SprintFunc()
	cyan = color.New(getColor(config.Colors["cyan"])).SprintFunc()
	green = color.New(getColor(config.Colors["green"])).SprintFunc()
	magenta = color.New(getColor(config.Colors["magenta"])).SprintFunc()
	blue = color.New(getColor(config.Colors["blue"])).SprintFunc()

	// Initialize commands and history
	Commands = config.Commands
	AdditionalCommands = config.AdditionalCommands
	CommandHistory = config.CommandHistory

	// Initialize git branch
	SetGitBranch()
}

func loadConfig(filename string) Config {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	return config
}

func ReadConfigs(filename string) (Config, error) {
	var AppSettings Config
	configFile, err := os.Open(filename)
	if err != nil {
		return Config{}, errors.New(fmt.Sprintf("Couldn't open config file. Error is: %s", err.Error()))
	}

	defer func(configFile *os.File) {
		err = configFile.Close()
		if err != nil {
			log.Fatal("Couldn't close config file. Error is: ", err.Error())
		}
	}(configFile)

	if err = json.NewDecoder(configFile).Decode(&AppSettings); err != nil {
		return Config{}, errors.New(fmt.Sprintf("Couldn't decode settings json file. Error is: %s", err.Error()))
	}

	return AppSettings, nil
}

func getColor(colorName string) color.Attribute {
	switch colorName {
	case "FgRed":
		return color.FgRed
	case "FgGreen":
		return color.FgGreen
	case "FgYellow":
		return color.FgYellow
	case "FgBlue":
		return color.FgBlue
	case "FgMagenta":
		return color.FgMagenta
	case "FgCyan":
		return color.FgCyan
	default:
		return color.FgWhite
	}
}

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
	{"redis", "Starts redis server"},
	{"redisserver", "Starts redis server"},
	{"redis-server", "Starts redis server"},
	{"ubuntu_redis", "Starts redis server"},
	{"redis_server", "Starts redis server"},
	{"panic", "Panics inside the command line"},
	{"print", "prints the text(args: font; example font=3d)"},
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
	{"redis", "Starts redis server"},
	{"redisserver", "Starts redis server"},
	{"redis-server", "Starts redis server"},
	{"ubuntu_redis", "Starts redis server"},
	{"redis_server", "Starts redis server"},
	{"panic", "Panics inside the command line"},
	{"print", "prints the text(args: font; example font=3d)"},
	{"pip", "Runs Python package installer"},
	{"py", "Runs Python interpreter"},
	{"npm", "Runs NPM package"},
	{"gcc", "Runs C compiled gcc"},
	{"deploy", "Deploys the application"},
	{"upgrade", "Upgrades installed packages"},
	{"export", "Exports data to a file"},
	{"import", "Imports data from a file"},
	{"compress", "Compresses files into an archive"},
	{"decompress", "Decompresses files from an archive"},
	{"convert", "Converts files from one format to another"},
	{"monitor", "Monitors system resources"},
	{"network", "Displays network information and status"},
	{"ubuntu", ""},
}

var CommandHistory []string

func Init() {
	// Initialize system variables
	system.Attempts = 0

	// Initialize CommandHistory with package or tool names
	CommandHistory = append(CommandHistory, "help")
	CommandHistory = append(CommandHistory, "run")
	CommandHistory = append(CommandHistory, "push")
	CommandHistory = append(CommandHistory, "pull")
	CommandHistory = append(CommandHistory, "origin")
	CommandHistory = append(CommandHistory, "main")
	CommandHistory = append(CommandHistory, "master")
	CommandHistory = append(CommandHistory, "merge")
	CommandHistory = append(CommandHistory, "run")
	CommandHistory = append(CommandHistory, "start")
	CommandHistory = append(CommandHistory, ".")
	CommandHistory = append(CommandHistory, "remote")
	CommandHistory = append(CommandHistory, "remote -v")
	CommandHistory = append(CommandHistory, "add")
	CommandHistory = append(CommandHistory, "add .")
	CommandHistory = append(CommandHistory, "add README.md")
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
	CommandHistory = append(CommandHistory, "3d")
	CommandHistory = append(CommandHistory, "2d")
	CommandHistory = append(CommandHistory, "font")
	CommandHistory = append(CommandHistory, "hello")
	CommandHistory = append(CommandHistory, "patch")
	CommandHistory = append(CommandHistory, "delete")
	CommandHistory = append(CommandHistory, "echo")
	CommandHistory = append(CommandHistory, "echo=on")
	CommandHistory = append(CommandHistory, "echo=off")
	CommandHistory = append(CommandHistory, "changelog")
	CommandHistory = append(CommandHistory, "beta")
	CommandHistory = append(CommandHistory, "http://localhost:6060")
	CommandHistory = append(CommandHistory, "https://github.com/bezhan2009/Orbix")
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
