package system

import (
	"path/filepath"
	"time"
)

// Commands available Orbix commands
var Commands = []Command{
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
	{"save", "Saves your custom vars to .env file"},
	{"delvar", "Deletes your custom variable"},
	{"gocode", "Executes go code"},
	{"sc", ""},
}

// AdditionalCommands additional commands
var AdditionalCommands = []Command{
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
	{"save", "Saves your custom vars to .env file"},
	{"delvar", "Deletes your custom variable"},
	{"gocode", "Executes go code"},
	{"sc", ""},
}

var (
	Absdir, _             = filepath.Abs("")
	RunningPath           = filepath.Join(Absdir, OrbixRunningUsersFileName)
	PreviousSessionPath   = ""
	PreviousSessionPrefix = ""
	Prefix                = ""
	ExecutingCommand      = false
	Unauthorized          = true
	RebootAttempts        = uint(0)
	SessionsStarted       = uint(0)
)

const (
	// MaxRetryAttempts Maximum number of restart attempts
	MaxRetryAttempts = 5
	// RetryDelay Delay before restart
	RetryDelay = 1 * time.Second
)

var (
	Port                = "6060"
	ErrorStartingServer = false
	UserName            = ""
	OrbixWorking        = false
	Localhost           = ""
	UserDir             = ""
	GitHubURL           = "https://github.com/bezhan2009/Orbix"
)