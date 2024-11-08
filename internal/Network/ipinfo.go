package Network

import (
	"encoding/json"
	"fmt"
	"goCmd/utils"
	"net/http"
)

func IPInfo(ip string) {
	resp, err := http.Get("http://ipinfo.io/" + ip + "/json")
	if err != nil {
		utils.AnimatedPrint(fmt.Sprint("Error fetching IP info:", err), "red")
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	for key, value := range result {
		utils.AnimatedPrint(fmt.Sprintf("%s: %v\n", key, value), "red")
	}
}
