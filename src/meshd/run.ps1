# run.ps1 — PiN meshd startup script
# Clears port conflicts and starts the meshd daemon cleanly.

$APIPort = 4002
$ListenPort = 4001

Write-Host "PiN — Pi Integrated Network" -ForegroundColor Cyan
Write-Host "Starting meshd..." -ForegroundColor Cyan

# Kill any process holding the API port
$netstatOutput = netstat -ano 2>$null
$netstatOutput | Select-String ":$APIPort\s" | ForEach-Object {
    $parts = $_.ToString().Trim() -split '\s+'
    $p = $parts[-1]
    if ($p -match '^\d+$' -and $p -ne '0') {
        Write-Host "Clearing port $APIPort (PID $p)..." -ForegroundColor Yellow
        taskkill /PID $p /F 2>$null | Out-Null
    }
}

# Kill any process holding the listen port
$netstatOutput | Select-String ":$ListenPort\s" | ForEach-Object {
    $parts = $_.ToString().Trim() -split '\s+'
    $p = $parts[-1]
    if ($p -match '^\d+$' -and $p -ne '0') {
        Write-Host "Clearing port $ListenPort (PID $p)..." -ForegroundColor Yellow
        taskkill /PID $p /F 2>$null | Out-Null
    }
}

# Brief pause to let ports release
Start-Sleep -Milliseconds 500

# Set CGO for SQLite
$env:CGO_ENABLED = 1

# Parse arguments — pass through to meshd
$args_str = $args -join " "
if ($args_str -eq "") {
    $args_str = "--dev"
}

Write-Host "Launching meshd $args_str" -ForegroundColor Green
Write-Host ""

# Start meshd
Invoke-Expression "go run . $args_str"