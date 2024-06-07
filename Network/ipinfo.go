package Network

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func IPInfo(ip string) {
	resp, err := http.Get("http://ipinfo.io/" + ip + "/json")
	if err != nil {
		fmt.Println("Error fetching IP info:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	for key, value := range result {
		fmt.Printf("%s: %v\n", key, value)
	}
}
