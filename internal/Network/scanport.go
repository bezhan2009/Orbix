package Network

import (
	"fmt"
	"goCmd/utils"
	"net"
	"strconv"
	"time"
)

func ScanPort(host string, ports []int) {
	for _, port := range ports {
		address := host + ":" + strconv.Itoa(port)
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err != nil {
			utils.AnimatedPrint(fmt.Sprintf("Port %d: Closed\n", port), "red")
			continue
		}
		conn.Close()
		utils.AnimatedPrint(fmt.Sprintf("Port %d: Open\n", port), "red")
	}
}
