package Network

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-ping/ping"
	"time"
)

func Ping(args []string) {
	red := color.New(color.FgRed).SprintFunc()
	if len(args) < 1 {
		fmt.Println(red("Usage: pingview <hostname>"))
		return
	}
	hostname := args[0]

	pinger, err := ping.NewPinger(hostname)
	if err != nil {
		fmt.Println(red("Error creating pinger:", err))
		return
	}

	pinger.Count = 4
	pinger.Timeout = time.Second * 10

	err = pinger.Run()
	if err != nil {
		fmt.Println(red("Error running ping:", err))
		return
	}

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Println("Successfully pinged " + green(hostname))

	stats := pinger.Statistics()
	fmt.Println(green("Ping statistics for", hostname, ":"))
	printStatistic := fmt.Sprintf("Packets: Sent = %d, Received = %d, Lost = %d (%.2f%% loss)\n",
		stats.PacketsSent, stats.PacketsRecv, stats.PacketsSent-stats.PacketsRecv, stats.PacketLoss)
	fmt.Println(yellow(printStatistic))
	fmt.Println(yellow("Approximate round trip times in milli-seconds:"))
	printStatistic = fmt.Sprintf("Minimum = %vms, Maximum = %vms, Average = %vms\n",
		stats.MinRtt.Milliseconds(), stats.MaxRtt.Milliseconds(), stats.AvgRtt.Milliseconds())
	fmt.Println(yellow(printStatistic))
}
