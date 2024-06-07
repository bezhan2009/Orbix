package Network

import (
	"fmt"
	"net"
)

func DNSLookup(domain string) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Println("Error resolving DNS:", err)
		return
	}
	for _, ip := range ips {
		fmt.Printf("%s IN A %s\n", domain, ip.String())
	}
}
