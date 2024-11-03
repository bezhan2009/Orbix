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

func SetCommands(session *system.Session) {
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
	SetGitBranch(session)
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
	{"newuser", "Adds a new user to Orbix"},
	{"template", "Executes a specific command template"},
	{"prompt", "Changes the command prompt in Orbix"},
	{"systemorbix", "Displays system information about Orbix"},
	{"rename", "Renames a file or directory"},
	{"rem", "Deletes a file"},
	{"del", "Deletes a file"},
	{"delete", "Deletes a file"},
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
	{"stusenv", "sets user from env file"},
	{"set_user_env", "sets user from env file"},
	{"setvar", "Sets a new value for variable in code"},
	{"getvar", "Gets value of variable in code"},
	{"scanport", "Scans open network ports"},
	{"dnslookup", "Performs DNS queries"},
	{"ipinfo", "Displays information about an IP address"},
	{"open_link", "Opens a URL in the default web browser"},
	{"geoip", "Displays geolocation information for an IP address"},
	{"api_request", "Sends an api request"},
	{"new_prompt", "Sets new prompt"},
	{"old_prompt", "Sets old prompt"},
	{"new_window", "Opens a new window"},
	{"kill", "kills process by PID"},
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
	{"delete", "Deletes a file"},
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
	{"stusenv", "sets user from env file"},
	{"set_user_env", "sets user from env file"},
	{"setvar", "Sets a new value for variable in code"},
	{"getvar", "Gets value of variable in code"},
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
	{"import", "Imports data from a file"},
	{"compress", "Compresses files into an archive"},
	{"decompress", "Decompresses files from an archive"},
	{"convert", "Converts files from one format to another"},
	{"monitor", "Monitors system resources"},
	{"network", "Displays network information and status"},
	{"api_request", "Sends an api request"},
	{"new_prompt", "Sets new prompt"},
	{"old_prompt", "Sets old prompt"},
	{"new_window", "Opens a new window"},
	{"kill", "kills process by PID"},
}

var CommandHistory []string

func Init(session *system.Session) {
	// Initialize CommandHistory with package or tool names
	session.CommandHistory = append(session.CommandHistory, "help")
	session.CommandHistory = append(session.CommandHistory, "run")
	session.CommandHistory = append(session.CommandHistory, "push")
	session.CommandHistory = append(session.CommandHistory, "pull")
	session.CommandHistory = append(session.CommandHistory, "origin")
	session.CommandHistory = append(session.CommandHistory, "main")
	session.CommandHistory = append(session.CommandHistory, "master")
	session.CommandHistory = append(session.CommandHistory, "merge")
	session.CommandHistory = append(session.CommandHistory, "run")
	session.CommandHistory = append(session.CommandHistory, "start")
	session.CommandHistory = append(session.CommandHistory, ".")
	session.CommandHistory = append(session.CommandHistory, "remote")
	session.CommandHistory = append(session.CommandHistory, "newthread")
	session.CommandHistory = append(session.CommandHistory, "neofetch")
	session.CommandHistory = append(session.CommandHistory, "location")
	session.CommandHistory = append(session.CommandHistory, "diruser")
	session.CommandHistory = append(session.CommandHistory, "prompt")
	session.CommandHistory = append(session.CommandHistory, "remote -v")
	session.CommandHistory = append(session.CommandHistory, "add")
	session.CommandHistory = append(session.CommandHistory, "add .")
	session.CommandHistory = append(session.CommandHistory, "add README.md")
	session.CommandHistory = append(session.CommandHistory, "--version")
	session.CommandHistory = append(session.CommandHistory, "install")
	session.CommandHistory = append(session.CommandHistory, "django")
	session.CommandHistory = append(session.CommandHistory, "flask")
	session.CommandHistory = append(session.CommandHistory, "config")
	session.CommandHistory = append(session.CommandHistory, "--global")
	session.CommandHistory = append(session.CommandHistory, "--timing")
	session.CommandHistory = append(session.CommandHistory, "--run-in-new-thread")
	session.CommandHistory = append(session.CommandHistory, "-t")
	session.CommandHistory = append(session.CommandHistory, "-m")
	session.CommandHistory = append(session.CommandHistory, "-am")
	session.CommandHistory = append(session.CommandHistory, "--list")
	session.CommandHistory = append(session.CommandHistory, "getvar *")
	session.CommandHistory = append(session.CommandHistory, "\"Your name\"")
	session.CommandHistory = append(session.CommandHistory, "\"your_email@example.com\"")
	session.CommandHistory = append(session.CommandHistory, "config")
	session.CommandHistory = append(session.CommandHistory, "--global user.name")
	session.CommandHistory = append(session.CommandHistory, "--global user.email")
	session.CommandHistory = append(session.CommandHistory, "branch")
	session.CommandHistory = append(session.CommandHistory, "checkout")
	session.CommandHistory = append(session.CommandHistory, "status")
	session.CommandHistory = append(session.CommandHistory, "commit")
	session.CommandHistory = append(session.CommandHistory, "clone")
	session.CommandHistory = append(session.CommandHistory, "log")
	session.CommandHistory = append(session.CommandHistory, "rebase")
	session.CommandHistory = append(session.CommandHistory, "cherry-pick")
	session.CommandHistory = append(session.CommandHistory, "stash")
	session.CommandHistory = append(session.CommandHistory, "reset")
	session.CommandHistory = append(session.CommandHistory, "diff")
	session.CommandHistory = append(session.CommandHistory, "grep")
	session.CommandHistory = append(session.CommandHistory, "fetch")
	session.CommandHistory = append(session.CommandHistory, "remote add")
	session.CommandHistory = append(session.CommandHistory, "remote remove")
	session.CommandHistory = append(session.CommandHistory, "tag")
	session.CommandHistory = append(session.CommandHistory, "show")
	session.CommandHistory = append(session.CommandHistory, "revert")
	session.CommandHistory = append(session.CommandHistory, "rm")
	session.CommandHistory = append(session.CommandHistory, "mv")
	session.CommandHistory = append(session.CommandHistory, "apply")
	session.CommandHistory = append(session.CommandHistory, "3d")
	session.CommandHistory = append(session.CommandHistory, "2d")
	session.CommandHistory = append(session.CommandHistory, "font")
	session.CommandHistory = append(session.CommandHistory, "hello")
	session.CommandHistory = append(session.CommandHistory, "patch")
	session.CommandHistory = append(session.CommandHistory, "delete")
	session.CommandHistory = append(session.CommandHistory, "echo")
	session.CommandHistory = append(session.CommandHistory, "echo=on")
	session.CommandHistory = append(session.CommandHistory, "echo=off")
	session.CommandHistory = append(session.CommandHistory, "changelog")
	session.CommandHistory = append(session.CommandHistory, "beta")
	session.CommandHistory = append(session.CommandHistory, system.Localhost)
	session.CommandHistory = append(session.CommandHistory, system.GitHubURL)
	session.CommandHistory = append(session.CommandHistory, "upgrade")
	session.CommandHistory = append(session.CommandHistory, "export")
	session.CommandHistory = append(session.CommandHistory, "import")
	session.CommandHistory = append(session.CommandHistory, "tar compress")
	session.CommandHistory = append(session.CommandHistory, "tar decompress")
	session.CommandHistory = append(session.CommandHistory, "convert")
	session.CommandHistory = append(session.CommandHistory, "nmap monitor")

	// Initialize session data
	SetGitBranch(session)
	SetPath(session)

	system.Attempts = 0
}

func InitColors() {
	colors := make(map[string]func(...interface{}) string)
	colors = system.GetColorsMap()

	red = colors["red"]
	yellow = colors["yellow"]
	cyan = colors["cyan"]
	green = colors["green"]
	magenta = colors["magenta"]
	blue = colors["blue"]
}

// OrbixFlags — список флагов, которые нужно удалить
var OrbixFlags = []string{"--timing", "-t", "--run-in-new-thread"}
