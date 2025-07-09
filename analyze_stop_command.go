package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	// The STOP command from your code
	stop := "\x03\x00\x00\x21\x02\xf0\x80\x32\x01\x00\x00\x06\x00\x00\x10\x00\x00\x29\x00\x00\x00\x00\x00\x09\x50\x5f\x50\x52\x4f\x47\x52\x41\x4d"

	// Convert to byte array
	stopBytes := []byte(stop)

	fmt.Println("=== S7 PLC STOP Command Analysis ===")
	fmt.Printf("Total Length: %d bytes\n", len(stopBytes))
	fmt.Printf("Hex Dump:\n%s\n", hex.Dump(stopBytes))

	fmt.Println("\n=== Protocol Breakdown ===")

	// TPKT Header (RFC 1006)
	fmt.Println("\nTPKT Header (ISO on TCP - RFC 1006):")
	fmt.Printf("  [0] Version: 0x%02X (3 = TPKT version 3)\n", stopBytes[0])
	fmt.Printf("  [1] Reserved: 0x%02X\n", stopBytes[1])
	tpktLen := (int(stopBytes[2]) << 8) | int(stopBytes[3])
	fmt.Printf("  [2-3] Length: 0x%02X%02X = %d bytes (total packet)\n", stopBytes[2], stopBytes[3], tpktLen)

	// COTP Header (ISO 8073)
	fmt.Println("\nCOTP Header (Connection-Oriented Transport Protocol):")
	fmt.Printf("  [4] Length: 0x%02X = %d bytes (remaining COTP header)\n", stopBytes[4], stopBytes[4])
	fmt.Printf("  [5] PDU Type: 0x%02X (0xF0 = DT Data)\n", stopBytes[5])
	fmt.Printf("  [6] TPDU Number: 0x%02X (0x80 = Last data unit)\n", stopBytes[6])

	// S7 Header
	fmt.Println("\nS7 Protocol Header:")
	fmt.Printf("  [7] Protocol ID: 0x%02X (0x32 = S7 protocol identifier)\n", stopBytes[7])
	fmt.Printf("  [8] Message Type: 0x%02X (1 = Job Request)\n", stopBytes[8])
	fmt.Printf("  [9-10] Reserved: 0x%02X%02X\n", stopBytes[9], stopBytes[10])
	pduRef := (int(stopBytes[11]) << 8) | int(stopBytes[12])
	fmt.Printf("  [11-12] PDU Reference: 0x%02X%02X = %d\n", stopBytes[11], stopBytes[12], pduRef)
	paramLen := (int(stopBytes[13]) << 8) | int(stopBytes[14])
	fmt.Printf("  [13-14] Parameter Length: 0x%02X%02X = %d bytes\n", stopBytes[13], stopBytes[14], paramLen)
	dataLen := (int(stopBytes[15]) << 8) | int(stopBytes[16])
	fmt.Printf("  [15-16] Data Length: 0x%02X%02X = %d bytes\n", stopBytes[15], stopBytes[16], dataLen)

	// S7 Parameters
	fmt.Println("\nS7 Job Request Parameters:")
	fmt.Printf("  [17] Function Code: 0x%02X (0x29 = PLC Control Services)\n", stopBytes[17])

	// Remaining bytes analysis
	fmt.Println("\nPLC Control Parameters:")
	fmt.Printf("  [18-23] Unknown/Reserved: ")
	for i := 18; i <= 23; i++ {
		fmt.Printf("0x%02X ", stopBytes[i])
	}
	fmt.Println()

	// ASCII part (P_PROGRAM)
	fmt.Printf("  [24] String Length: 0x%02X = %d\n", stopBytes[24], stopBytes[24])
	asciiPart := string(stopBytes[25:])
	fmt.Printf("  [25-33] ASCII String: '%s'\n", asciiPart)
	fmt.Printf("           Hex values: ")
	for i := 25; i < len(stopBytes); i++ {
		fmt.Printf("0x%02X ", stopBytes[i])
	}
	fmt.Println()

	fmt.Println("\n=== Command Analysis ===")
	fmt.Println("This command structure indicates:")
	fmt.Println("1. Valid S7 protocol packet with proper TPKT/COTP headers")
	fmt.Println("2. Function 0x29 = PLC Control Services (START/STOP operations)")
	fmt.Println("3. Parameter 'P_PROGRAM' = Stop program execution")
	fmt.Println("4. This is a legitimate S7 STOP CPU command")
	fmt.Println("\nNOTE: This command will only work if:")
	fmt.Println("- The PLC has no password protection")
	fmt.Println("- The PLC allows remote STOP commands")
	fmt.Println("- The communication is not encrypted/authenticated")

	// Let's also show what a START command might look like
	fmt.Println("\n=== Related Commands ===")
	fmt.Println("For reference, related S7 control commands use similar structure:")
	fmt.Println("- STOP: Function 0x29 with 'P_PROGRAM'")
	fmt.Println("- COLD RESTART: Function 0x28")
	fmt.Println("- WARM RESTART: Function 0x28 with different parameters")
}
