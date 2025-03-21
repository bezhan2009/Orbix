package utils

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/fatih/color"
	"goCmd/internal/Network"
)

func OpenLinkUtil(commandArgs []string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	if len(commandArgs) < 1 {
		fmt.Println(yellow("Usage: open_link <url>"))
		return
	}

	rawUrl := commandArgs[0]

	// Добавить "http://" если нет протокола
	if !strings.Contains(rawUrl, "://") {
		rawUrl = "http://" + rawUrl
	}

	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Println(red("Invalid URL format"))
		return
	}

	validProtocols := map[string]bool{
		"http":   true,
		"https":  true,
		"ftp":    true,
		"mailto": true,
		"file":   true,
	}

	protocol := parsedUrl.Scheme
	if !validProtocols[protocol] {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Println(red("Incorrect URL"))
		fmt.Println("Your URL address has an unsupported protocol!")
		fmt.Println("Valid protocols are:")
		for k := range validProtocols {
			fmt.Println("\t" + k)
		}
		return
	}

	if parsedUrl.Host == "" {
		parsedUrl.Host = parsedUrl.Path
		parsedUrl.Path = ""
	}
	if !strings.Contains(parsedUrl.Host, ".") {
		parsedUrl.Host += ".com"
	}

	err = Network.OpenBrowser(parsedUrl.String())
	if err != nil {
		fmt.Println(fmt.Sprintf("Error opening browser: %s", err))
	}
}
