package Network

import (
	"fmt"
	"github.com/go-ping/ping"
	"time"
)

func Ping(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: pingview <hostname>")
		return
	}
	hostname := args[0]

	pinger, err := ping.NewPinger(hostname)
	if err != nil {
		fmt.Println("Error creating pinger:", err)
		return
	}

	pinger.Count = 4
	pinger.Timeout = time.Second * 10

	err = pinger.Run()
	if err != nil {
		fmt.Println("Error running ping:", err)
		return
	}

	stats := pinger.Statistics()
	fmt.Printf("Ping statistics for %s:\n", hostname)
	fmt.Printf("Packets: Sent = %d, Received = %d, Lost = %d (%.2f%% loss)\n",
		stats.PacketsSent, stats.PacketsRecv, stats.PacketsSent-stats.PacketsRecv, stats.PacketLoss)
	fmt.Printf("Approximate round trip times in milli-seconds:\n")
	fmt.Printf("Minimum = %vms, Maximum = %vms, Average = %vms\n",
		stats.MinRtt.Milliseconds(), stats.MaxRtt.Milliseconds(), stats.AvgRtt.Milliseconds())
}
