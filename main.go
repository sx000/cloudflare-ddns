package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	configFile  = "/etc/cloudflare-ddns.conf"
	ipFile      = "/var/lib/cloudflare-ddns/current_ip"
	checkEvery  = 10 * time.Minute
	ipCheckURL  = "https://api.ipify.org?format=json"
	cfAPI       = "https://api.cloudflare.com/client/v4"
)

type Config struct {
	APIToken   string `json:"api_token"`
	ZoneName   string `json:"zone_name"`
	RecordName string `json:"record_name"`
}

type CloudflareResponse struct {
	Success bool `json:"success"`
	Result  []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Content string `json:"content"`
	} `json:"result"`
}

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting Cloudflare DDNS updater")

	cfg := loadConfig()
	client := &http.Client{Timeout: 10 * time.Second}

	for {
		currentIP := getPublicIP()
		if currentIP == "" {
			time.Sleep(checkEvery)
			continue
		}

		savedIP := readLastIP()
		if currentIP != savedIP {
			if updateDNS(client, cfg, currentIP) {
				writeLastIP(currentIP)
				log.Printf("Successfully updated %s to %s\n", cfg.RecordName, currentIP)
			}
		}

		time.Sleep(checkEvery)
	}
}

func loadConfig() Config {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	return cfg
}

func getPublicIP() string {
	resp, err := http.Get(ipCheckURL)
	if err != nil {
		log.Printf("IP check error: %v", err)
		return ""
	}
	defer resp.Body.Close()

	var result struct{ IP string `json:"ip"` }
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("IP parse error: %v", err)
		return ""
	}

	return result.IP
}

func readLastIP() string {
	data, err := ioutil.ReadFile(ipFile)
	if err != nil {
		return ""
	}
	return string(data)
}

func writeLastIP(ip string) {
	os.MkdirAll(filepath.Dir(ipFile), 0755)
	ioutil.WriteFile(ipFile, []byte(ip), 0644)
}

func updateDNS(client *http.Client, cfg Config, newIP string) bool {
	zoneID := getZoneID(client, cfg)
	if zoneID == "" {
		return false
	}

	recordID, currentDNSIP := getDNSRecord(client, zoneID, cfg)
	if currentDNSIP == newIP {
		return true
	}

	return updateDNSRecord(client, zoneID, recordID, cfg.RecordName, newIP, cfg.APIToken)
}

func getZoneID(client *http.Client, cfg Config) string {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/zones?name=%s", cfAPI, cfg.ZoneName), nil)
	req.Header.Set("Authorization", "Bearer "+cfg.APIToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Zone lookup error: %v", err)
		return ""
	}
	defer resp.Body.Close()

	var result CloudflareResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || !result.Success {
		log.Printf("Zone lookup failed: %v", err)
		return ""
	}

	if len(result.Result) == 0 {
		log.Printf("Zone %s not found", cfg.ZoneName)
		return ""
	}

	return result.Result[0].ID
}

func getDNSRecord(client *http.Client, zoneID string, cfg Config) (string, string) {
	url := fmt.Sprintf("%s/zones/%s/dns_records?type=A&name=%s", cfAPI, zoneID, cfg.RecordName)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+cfg.APIToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("DNS record lookup error: %v", err)
		return "", ""
	}
	defer resp.Body.Close()

	var result CloudflareResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || !result.Success {
		log.Printf("DNS record lookup failed: %v", err)
		return "", ""
	}

	if len(result.Result) == 0 {
		log.Printf("DNS record %s not found", cfg.RecordName)
		return "", ""
	}

	return result.Result[0].ID, result.Result[0].Content
}

func updateDNSRecord(client *http.Client, zoneID, recordID, name, newIP, token string) bool {
	data := map[string]interface{}{
		"type":    "A",
		"name":    name,
		"content": newIP,
		"ttl":     120,
		"proxied": false,
	}

	body, _ := json.Marshal(data)
	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", cfAPI, zoneID, recordID)

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Update error: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}