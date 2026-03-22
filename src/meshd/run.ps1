$env:CGO_ENABLED = 1
netstat -ano | Select-String ":4002 " | ForEach-Object {
    $p = ($_.ToString().Trim() -split "\s+")[-1]
    if ($p -match "^\d+$" -and $p -ne "0") {
        taskkill /PID $p /F 2>$null | Out-Null
    }
}
netstat -ano | Select-String ":4001 " | ForEach-Object {
    $p = ($_.ToString().Trim() -split "\s+")[-1]
    if ($p -match "^\d+$" -and $p -ne "0") {
        taskkill /PID $p /F 2>$null | Out-Null
    }
}
Start-Sleep -Milliseconds 500
go run . --dev