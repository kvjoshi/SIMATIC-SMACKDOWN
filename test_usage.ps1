# SIMATIC-SMACKDOWN Test Script for PowerShell

Write-Host "=== SIMATIC-SMACKDOWN Test Script ===" -ForegroundColor Cyan
Write-Host ""

Write-Host "Building the project..." -ForegroundColor Yellow
go build -o simatic_smackdown.exe main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Test 1: Show usage with invalid IP" -ForegroundColor Green
Write-Host "Command: .\simatic_smackdown.exe invalid-ip" -ForegroundColor Gray
.\simatic_smackdown.exe invalid-ip
Write-Host ""
Write-Host "----------------------------------------" -ForegroundColor DarkGray
Write-Host "Press any key to continue..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

Write-Host ""
Write-Host "Test 2: Target specific IPs (simulation - no real PLCs)" -ForegroundColor Green
Write-Host "Command: .\simatic_smackdown.exe 192.168.1.50 192.168.1.51 10.0.0.100" -ForegroundColor Gray
.\simatic_smackdown.exe 192.168.1.50 192.168.1.51 10.0.0.100
Write-Host ""
Write-Host "----------------------------------------" -ForegroundColor DarkGray
Write-Host "Press any key to continue..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

Write-Host ""
Write-Host "Test 3: Mixed valid and invalid IPs" -ForegroundColor Green
Write-Host "Command: .\simatic_smackdown.exe 192.168.1.50 not-an-ip 10.0.0.1 256.256.256.256" -ForegroundColor Gray
.\simatic_smackdown.exe 192.168.1.50 not-an-ip 10.0.0.1 256.256.256.256
Write-Host ""
Write-Host "----------------------------------------" -ForegroundColor DarkGray

Write-Host ""
Write-Host "Test 4: Network scan mode (no arguments)" -ForegroundColor Green
Write-Host "Command: .\simatic_smackdown.exe" -ForegroundColor Gray
Write-Host "WARNING: This will scan your local network!" -ForegroundColor Red
Write-Host "Press Ctrl+C to cancel, or any key to continue..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
.\simatic_smackdown.exe
Write-Host ""
Write-Host "----------------------------------------" -ForegroundColor DarkGray

Write-Host ""
Write-Host "All tests completed." -ForegroundColor Cyan
Write-Host "Press any key to exit..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
