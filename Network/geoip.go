package Network

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GeoIP(ip string) {
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		fmt.Println("Error fetching Geo IP info:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	for key, value := range result {
		fmt.Printf("%s: %v\n", key, value)
	}
}
