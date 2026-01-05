# Database Connection Troubleshooting

## Current Issue
DNS resolution failing: `lookup db.pikfcjbdzairltrnxnwt.supabase.co: no such host`

## Root Cause
Your Supabase project shows "Not IPv4 compatible" for direct connections. The hostname resolves to IPv6, but Windows/Go DNS resolver can't connect.

## Solution Steps

### Step 1: Get Session Pooler Connection String from Supabase Dashboard

1. Go to your Supabase project dashboard
2. Click **"Connect"** button (top right)
3. Select **"Connection String"** tab
4. Change **"Method"** dropdown to **"Session Pooler"** (NOT "Direct connection")
5. Copy the connection string shown

### Step 2: Update .env File

Replace your current `DATABASE_URL` with the Session Pooler connection string from Step 1.

**Expected format:**
```
DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.pikfcjbdzairltrnxnwt.supabase.co:6543/postgres?sslmode=require
```

**Important:** 
- Port should be **6543** (Session Pooler), not 5432
- Must include `sslmode=require`
- Replace `[YOUR-PASSWORD]` with your actual database password

### Step 3: Test Connection

Run the diagnostic script:
```bash
go run ./scripts/test_db_connection.go
```

Or start the server:
```bash
go run ./cmd/server
```

## Alternative: Use IPv4 Add-on

If Session Pooler doesn't work, you can:
1. Purchase Supabase IPv4 add-on (from the dashboard)
2. Or use a VPN/proxy that supports IPv6

## Network Troubleshooting

If DNS still fails:

1. **Check DNS settings:**
   ```powershell
   Get-DnsClientServerAddress
   ```

2. **Try different DNS server:**
   ```powershell
   Set-DnsClientServerAddress -InterfaceAlias "Ethernet" -ServerAddresses "8.8.8.8","8.8.4.4"
   ```

3. **Flush DNS cache:**
   ```powershell
   ipconfig /flushdns
   ```

## Current .env Format

Your `.env` should look like:
```env
PORT=8080
APP_ENV=development
DATABASE_URL=postgresql://postgres:ySihqziq0WgqUwMD@db.pikfcjbdzairltrnxnwt.supabase.co:6543/postgres?sslmode=require
ALLOWED_EMAIL_DOMAIN=citchennai.net
```

**Note:** Make sure you're using the **Session Pooler** connection string from Supabase dashboard, not the direct connection.

