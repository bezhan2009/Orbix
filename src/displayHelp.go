package src

import (
	"fmt"
)

func displayHelp() {
	helpText := `
For command information, type HELP
CREATE             creates a new file	
CLEAN              clears the screen	
CD                 changes the current directory
COPUSOURCE         copies the source code of the file
CF                 Creates a folder
DNSLOOKUP          DNS queries
DF                 Deletes a folder
LS                 lists directory contents
NEWSHABLON         creates a new command template for execution
REMOVE             deletes a file
READ               displays the contents of a file
PROMPT             changes src prompt.
PINGVIEW           shows ping.
PRIMES             finds large prime numbers
PICALC             calculates the value of π.
NEWUSER            creates a new user for src.
NEWCOMMAND         created a new command
ORBIX              starts another ORBIX session
SHABLON            executes a specific command template
SYSTEMGOCMD        displays information about src
SYSTEMINFO         displays system information
SIGNOUT            user signs out of src
TREE               graphically displays directory structure
WRITE              writes data to a file
EDIT               edits a file
WIFIUTILS          launches a utility for working with WiFi
EXTRACTZIP         extracts .zip archives
SCANPORT           scans ports
OPEN_LINK          opens the link in the browser
WHOIS              domain information
FILEIO             intensive file operation test
IPINFO             IP address information
GEOIP              IP address geolocation
MATRIXMUL          multiplies large matrices
EXIT               exit
`
	fmt.Println(helpText)
}
