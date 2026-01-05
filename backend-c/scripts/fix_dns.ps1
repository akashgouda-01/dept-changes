# PowerShell script to fix DNS resolution issues
# Run this as Administrator if DNS resolution fails

Write-Host "🔍 Checking DNS resolution for Supabase hostname..." -ForegroundColor Cyan

$hostname = "db.pikfcjbdzairltrnxnwt.supabase.co"

# Try resolving with different DNS servers
Write-Host "`n1. Trying Google DNS (8.8.8.8)..." -ForegroundColor Yellow
try {
    $result = Resolve-DnsName $hostname -Server 8.8.8.8 -ErrorAction Stop
    Write-Host "✅ Resolved successfully!" -ForegroundColor Green
    $result | Select-Object Name, IPAddress, Type | Format-Table
} catch {
    Write-Host "❌ Failed: $_" -ForegroundColor Red
}

Write-Host "`n2. Trying Cloudflare DNS (1.1.1.1)..." -ForegroundColor Yellow
try {
    $result = Resolve-DnsName $hostname -Server 1.1.1.1 -ErrorAction Stop
    Write-Host "✅ Resolved successfully!" -ForegroundColor Green
    $result | Select-Object Name, IPAddress, Type | Format-Table
} catch {
    Write-Host "❌ Failed: $_" -ForegroundColor Red
}

Write-Host "`n3. Flushing DNS cache..." -ForegroundColor Yellow
ipconfig /flushdns | Out-Null
Write-Host "✅ DNS cache flushed" -ForegroundColor Green

Write-Host "`n4. Testing network connectivity..." -ForegroundColor Yellow
try {
    $test = Test-NetConnection -ComputerName $hostname -Port 6543 -WarningAction SilentlyContinue
    if ($test.TcpTestSucceeded) {
        Write-Host "✅ Port 6543 is reachable!" -ForegroundColor Green
    } else {
        Write-Host "❌ Port 6543 is not reachable" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ Network test failed: $_" -ForegroundColor Red
}

Write-Host "`n💡 If DNS resolution still fails:" -ForegroundColor Cyan
Write-Host "   1. Check if your network/firewall blocks DNS queries" -ForegroundColor White
Write-Host "   2. Try using a VPN" -ForegroundColor White
Write-Host "   3. Contact your network administrator" -ForegroundColor White
Write-Host "   4. Use Supabase IPv4 add-on (paid feature)" -ForegroundColor White

