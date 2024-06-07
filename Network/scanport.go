package Network

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func ScanPort(host string, ports []int) {
	for _, port := range ports {
		address := host + ":" + strconv.Itoa(port)
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err != nil {
			fmt.Printf("Port %d: Closed\n", port)
			continue
		}
		conn.Close()
		fmt.Printf("Port %d: Open\n", port)
	}
}
