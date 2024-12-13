package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
    "os"
	"net/http"
)

// IPResponse represents the structure of the IP address response
type IPResponse struct {
	IP string `json:"ip"`
}

// getPublicIP fetches the public IP address from ipify.org
func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", fmt.Errorf("failed to get IP: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var ipResp IPResponse
	err = json.Unmarshal(body, &ipResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}

	return ipResp.IP, nil
}

// ipHandler handles the IP lookup request
func ipHandler(w http.ResponseWriter, r *http.Request) {
	ip, err := getPublicIP()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := IPResponse{IP: ip}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback to default if PORT not set
	}

	http.HandleFunc("/ip", ipHandler)
	
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}