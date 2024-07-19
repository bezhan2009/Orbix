package NetworkScan

import (
	"fmt"
	"github.com/google/gopacket/pcap"
)

func WiFi() {
	ifaces, err := pcap.FindAllDevs()
	if err != nil {
		fmt.Println("Ошибка при получении списка интерфейсов:", err)
		return
	}
	for _, iface := range ifaces {
		fmt.Printf("Интерфейс: %s\n", iface.Name)
		for _, addr := range iface.Addresses {
			fmt.Printf("  Адрес: %s\n", addr.IP)
		}
	}
}
