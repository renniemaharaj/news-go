package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	apiToken   string
	zoneID     string
	recordID   string
	recordName string
	lastIP     string

	l = createLogger()
)

// Initialize is the main entry point
func Initialize() {
	if !validateEnv() {
		l.Error("Missing required environment variables. Please check your .env file.")
		os.Exit(1)
	}

	if err := TestCloudflareAccess(); err != nil {
		l.Error("Cloudflare access test failed: " + err.Error())
		os.Exit(1)
	}
	l.Success("Cloudflare access verified")

	startDNSLoop()
}

// Validate environment variables
func validateEnv() bool {
	if err := godotenv.Load(); err != nil {
		l.Error("No .env file found, or error loading it: " + err.Error())
	}

	missing := []string{}
	apiToken = os.Getenv("CF_API_TOKEN")
	if apiToken == "" {
		missing = append(missing, "CF_API_TOKEN")
	}
	zoneID = os.Getenv("CF_ZONE_ID")
	if zoneID == "" {
		missing = append(missing, "CF_ZONE_ID")
	}
	recordID = os.Getenv("CF_RECORD_ID")
	if recordID == "" {
		missing = append(missing, "CF_RECORD_ID")
	}
	recordName = os.Getenv("CF_RECORD_NAME")
	if recordName == "" {
		missing = append(missing, "CF_RECORD_NAME")
	}
	lastIP = os.Getenv("LOCAL_IP")
	if lastIP == "" {
		missing = append(missing, "LOCAL_IP")
	}
	if len(missing) > 0 {
		l.Error("Missing required environment variables: " + strings.Join(missing, ", "))
		return false
	}
	return true
}

// TestCloudflareAccess checks if the DNS record exists and API token works
func TestCloudflareAccess() error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if !result["success"].(bool) {
		return fmt.Errorf("cloudflare API error: %v", result["errors"])
	}

	return nil
}

// DNS update loop
func startDNSLoop() {
	for {
		currentIP, err := getPublicIP()
		if err != nil {
			l.Error("Failed to get public IP: " + err.Error())
			time.Sleep(5 * time.Minute)
			continue
		}

		if currentIP != lastIP {
			l.Warning(fmt.Sprintf("🔄 IP changed: %s → %s", lastIP, currentIP))

			if err := updateDNSRecord(currentIP); err != nil {
				l.Error("DNS update failed: " + err.Error())
			} else {
				l.Success("DNS record updated")
				os.Setenv("LOCAL_IP", currentIP)
				lastIP = currentIP
			}
		} else {
			l.Info("IP unchanged: " + currentIP)
		}

		time.Sleep(5 * time.Minute)
	}
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(ip)), err
}

func updateDNSRecord(ip string) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)

	payload := map[string]interface{}{
		"type":    "A",
		"name":    recordName,
		"content": ip,
		"ttl":     1,
		"proxied": true,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if !result["success"].(bool) {
		return fmt.Errorf("cloudflare API error: %v", result["errors"])
	}
	return nil
}
