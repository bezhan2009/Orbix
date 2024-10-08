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
REMOVE             deletes a file.
READ               displays the contents of a file.
REDIS-SERVER       Starts redis server
PROMPT             changes src prompt.
PANIC              Panics inside the command line.
PINGVIEW           shows ping.
PRIMES             finds large prime numbers
PICALC             calculates the value of π.
PRINT              prints the text(args: font; example font=3d).
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
NEOFETCH           displays information about the system
EXIT               exit
`
	fmt.Println(helpText)
}

func displayHelpBeta() {
	fmt.Println(yellow("\tFor command information, type HELP\n"))
	const nameWidth = 20 // Задаем ширину для названий команд

	for _, command := range Commands {
		if command.Description != "" {
			// Используем форматирование с указанием ширины поля для команд
			fmt.Println(yellow(fmt.Sprintf("%-*s %s", nameWidth, command.Name, command.Description)))
		}
	}
}
