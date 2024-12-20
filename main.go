// package main

// import (
//     "fmt"
//     "log"
//     "net"
//     "net/http"
//     "strings"
//     "time"
// )

// // getClientIP extracts the real client IP from various headers
// func getClientIP(r *http.Request) string {
//     // Order of headers to check
//     headers := []string{
//         "HTTP_CLIENT_IP",
//         "HTTP_X_FORWARDED_FOR",
//         "HTTP_X_FORWARDED",
//         "HTTP_X_CLUSTER_CLIENT_IP",
//         "HTTP_FORWARDED_FOR",
//         "HTTP_FORWARDED",
//         "X-Forwarded-For",
//         "X-Real-IP",
//     }

//     // Check each header
//     for _, header := range headers {
//         if ip := r.Header.Get(header); ip != "" {
//             // Split in case of multiple IPs
//             ips := strings.Split(ip, ",")
//             for _, ip := range ips {
//                 // Clean the IP address
//                 ip = strings.TrimSpace(ip)
                
//                 // Parse the IP to validate it
//                 parsedIP := net.ParseIP(ip)
//                 if parsedIP == nil {
//                     continue
//                 }

//                 // Check if it's a private or reserved IP
//                 if !parsedIP.IsPrivate() && !parsedIP.IsLoopback() && !parsedIP.IsLinkLocalUnicast() {
//                     return ip
//                 }
//             }
//         }
//     }

//     // Fall back to RemoteAddr if no valid IP found in headers
//     ip, _, _ := net.SplitHostPort(r.RemoteAddr)
//     return ip
// }

// // logIPHandler handles the IP logging endpoint
// func logIPHandler(w http.ResponseWriter, r *http.Request) {
//     ip := getClientIP(r)
    
//     // Log the IP with timestamp
//     log.Printf("Client IP: %s, Time: %s", ip, time.Now().Format(time.RFC3339))
    
//     // Return the IP to the client
//     fmt.Fprintf(w, "Your IP: %s", ip)
// }

// func main() {
//     // Configure logging
//     log.SetFlags(log.LstdFlags | log.Lmicroseconds)
    
//     // Register the handler for the /log-ip endpoint
//     http.HandleFunc("/log-ip", logIPHandler)
    
//     // Start the server
//     fmt.Println("Server starting on :8080...")
//     if err := http.ListenAndServe(":8080", nil); err != nil {
//         log.Fatal(err)
//     }
// }

package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net"
    "net/http"
    "strings"
    "time"
)

// IPDetails stores the response from ip-api.com
type IPDetails struct {
    Status      string  `json:"status"`
    Country     string  `json:"country"`
    CountryCode string  `json:"countryCode"`
    Region      string  `json:"region"`
    RegionName  string  `json:"regionName"`
    City        string  `json:"city"`
    Zip         string  `json:"zip"`
    Lat         float64 `json:"lat"`
    Lon         float64 `json:"lon"`
    Timezone    string  `json:"timezone"`
    ISP         string  `json:"isp"`
    Org         string  `json:"org"`
    AS          string  `json:"as"`
    Query       string  `json:"query"`
}

// getIPDetails fetches IP information from ip-api.com
// func getIPDetails(ip string) (*IPDetails, error) {
//     url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
    
//     resp, err := http.Get(url)
//     if err != nil {
//         return nil, fmt.Errorf("failed to fetch IP details: %v", err)
//     }
//     defer resp.Body.Close()

//     body, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return nil, fmt.Errorf("failed to read response: %v", err)
//     }

//     var details IPDetails
//     if err := json.Unmarshal(body, &details); err != nil {
//         return nil, fmt.Errorf("failed to parse response: %v", err)
//     }

//     return &details, nil
// }
func getIPDetails(ip string) (*IPDetails, error) {
    // Create custom client
    client := &http.Client{}
    
    // Create request with the HTTPS endpoint
    url := fmt.Sprintf("https://ipapi.co/%s/json/", ip)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }
    
    // Set required headers
    req.Header.Set("User-Agent", "ipapi.co/#go-v1.5")
    
    // Make the request
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch IP details: %v", err)
    }
    defer resp.Body.Close()

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %v", err)
    }

    // Check if response status is not OK
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
    }

    // Parse the response
    var details IPDetails
    if err := json.Unmarshal(body, &details); err != nil {
        return nil, fmt.Errorf("failed to parse response: %v", err)
    }

    return &details, nil
}
// getClientIP extracts the real client IP from various headers
func getClientIP(r *http.Request) string {
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

    for _, header := range headers {
        if ip := r.Header.Get(header); ip != "" {
            ips := strings.Split(ip, ",")
            for _, ip := range ips {
                ip = strings.TrimSpace(ip)
                parsedIP := net.ParseIP(ip)
                if parsedIP == nil {
                    continue
                }

                if !parsedIP.IsPrivate() && !parsedIP.IsLoopback() && !parsedIP.IsLinkLocalUnicast() {
                    return ip
                }
            }
        }
    }

    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    return ip
}

// logIPHandler handles the IP logging endpoint
func logIPHandler(w http.ResponseWriter, r *http.Request) {
    ip := getClientIP(r)
    
    // Get IP details
    details, err := getIPDetails(ip)
    if err != nil {
        log.Printf("Error getting IP details: %v", err)
        http.Error(w, "Failed to get IP details", http.StatusInternalServerError)
        return
    }
    
    // Log the IP with details and timestamp
    log.Printf("Client IP: %s, Location: %s, %s, ISP: %s, Time: %s",
        ip,
        details.City,
        details.Country,
        details.ISP,
        time.Now().Format(time.RFC3339))
    
    // Return the details to the client as JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(details)
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