@echo off
echo === SIMATIC-SMACKDOWN Test Script ===
echo.

echo Building the project...
go build -o simatic_smackdown.exe main.go
if %errorlevel% neq 0 (
    echo Build failed!
    exit /b 1
)

echo.
echo Test 1: Show usage with invalid IP
echo Command: simatic_smackdown.exe invalid-ip
simatic_smackdown.exe invalid-ip
echo.
echo ----------------------------------------
pause

echo.
echo Test 2: Target specific IPs (simulation - no real PLCs)
echo Command: simatic_smackdown.exe 192.168.1.50 192.168.1.51
simatic_smackdown.exe 192.168.1.50 192.168.1.51
echo.
echo ----------------------------------------
pause

echo.
echo Test 3: Network scan mode (no arguments)
echo Command: simatic_smackdown.exe
echo WARNING: This will scan your local network!
echo Press Ctrl+C to cancel, or any key to continue...
pause >nul
simatic_smackdown.exe
echo.
echo ----------------------------------------

echo.
echo All tests completed.
pause
