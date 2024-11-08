package Network

import (
	"encoding/json"
	"fmt"
	"goCmd/utils"
	"net/http"
)

func GeoIP(ip string) {
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		utils.AnimatedPrint(fmt.Sprint("Error fetching Geo IP info:", err, "\n"), "red")
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	for key, value := range result {
		utils.AnimatedPrint(fmt.Sprintf("%s: %v\n", key, value), "red")
	}
}
