// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
//     "os"
// 	"net/http"
// )

// // IPResponse represents the structure of the IP address response
// type IPResponse struct {
// 	IP string `json:"ip"`
// }

// // getPublicIP fetches the public IP address from ipify.org
// func getPublicIP() (string, error) {
// 	resp, err := http.Get("https://api.ipify.org?format=json")
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get IP: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read response: %v", err)
// 	}

// 	var ipResp IPResponse
// 	err = json.Unmarshal(body, &ipResp)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to parse JSON: %v", err)
// 	}

// 	return ipResp.IP, nil
// }

// // ipHandler handles the IP lookup request
// func ipHandler(w http.ResponseWriter, r *http.Request) {
// 	ip, err := getPublicIP()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	response := IPResponse{IP: ip}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// func main() {
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080" // fallback to default if PORT not set
// 	}

// 	http.HandleFunc("/ip", ipHandler)
	
// 	fmt.Printf("Server starting on port %s...\n", port)
// 	log.Fatal(http.ListenAndServe(":" + port, nil))
// }

package main

import (
    "fmt"
    "log"
    "net"
    "net/http"
    "strings"
    "time"
)

// getClientIP extracts the real client IP from various headers
func getClientIP(r *http.Request) string {
    // Order of headers to check
    headers := []string{
        "HTTP_CLIENT_IP",
        "HTTP_X_FORWARDED_FOR",
        "HTTP_X_FORWARDED",
        "HTTP_X_CLUSTER_CLIENT_IP",
        "HTTP_FORWARDED_FOR",
        "HTTP_FORWARDED",
        "X-Forwarded-For",
        "X-Real-IP",
    }

    // Check each header
    for _, header := range headers {
        if ip := r.Header.Get(header); ip != "" {
            // Split in case of multiple IPs
            ips := strings.Split(ip, ",")
            for _, ip := range ips {
                // Clean the IP address
                ip = strings.TrimSpace(ip)
                
                // Parse the IP to validate it
                parsedIP := net.ParseIP(ip)
                if parsedIP == nil {
                    continue
                }

                // Check if it's a private or reserved IP
                if !parsedIP.IsPrivate() && !parsedIP.IsLoopback() && !parsedIP.IsLinkLocalUnicast() {
                    return ip
                }
            }
        }
    }

    // Fall back to RemoteAddr if no valid IP found in headers
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    return ip
}

// logIPHandler handles the IP logging endpoint
func logIPHandler(w http.ResponseWriter, r *http.Request) {
    ip := getClientIP(r)
    
    // Log the IP with timestamp
    log.Printf("Client IP: %s, Time: %s", ip, time.Now().Format(time.RFC3339))
    
    // Return the IP to the client
    fmt.Fprintf(w, "Your IP: %s", ip)
}

func main() {
    // Configure logging
    log.SetFlags(log.LstdFlags | log.Lmicroseconds)
    
    // Register the handler for the /log-ip endpoint
    http.HandleFunc("/log-ip", logIPHandler)
    
    // Start the server
    fmt.Println("Server starting on :8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}