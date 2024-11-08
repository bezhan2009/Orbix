package Network

import (
	"fmt"
	"goCmd/utils"
	"net"
)

func DNSLookup(domain string) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		utils.AnimatedPrint(fmt.Sprint("Error resolving DNS:", err), "red")
		return
	}
	for _, ip := range ips {
		utils.AnimatedPrint(fmt.Sprintf("%s IN A %s\n", domain, ip.String()), "red")
	}
}
