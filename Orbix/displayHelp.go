package Orbix

import "goCmd/debug"

func displayHelp(commandArgs []string, user, dir string) {
	helpText := `
For command information, type HELP
CREATE             creates a new file
CLEAN              clears the screen
CD                 changes the current directory
COPUSOURCE         copies the source code of the file
LS                 lists directory contents
NEWSHABLON         creates a new command template for execution
REMOVE             deletes a file
READ               displays the contents of a file
PROMPT             changes Orbix prompt.
PINGVIEW           shows ping.
PRIMES             finds large prime numbers
PICALC             calculates the value of π.
NEWUSER            creates a new user for Orbix.
NEWCOMMAND         created a new command
Orbix              starts another Orbix session
SHABLON            executes a specific command template
SYSTEMGOCMD        displays information about Orbix
SYSTEMINFO         displays system information
SIGNOUT            user signs out of Orbix
TREE               graphically displays directory structure
WRITE              writes data to a file
EDIT               edits a file
WIFIUTILS          launches a utility for working with WiFi
EXTRACTZIP         extracts .zip archives
SCANPORT           scans ports
WHOIS              domain information
DNSLOOKUP          DNS queries
FILEIO             intensive file operation test
IPINFO             IP address information
GEOIP              IP address geolocation
MATRIXMUL          multiplies large matrices
MYCMD              native command interpreter
EXIT               exit
`
	animatedPrint(helpText)
	errDebug := debug.Commands("help", true, commandArgs, user, dir)
	if errDebug != nil {
		animatedPrint(errDebug.Error() + "\n")
	}
}
