$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $Root

$env:GOCACHE = Join-Path $Root ".gocache"
$env:GOMODCACHE = Join-Path $Root ".gomodcache"
$env:GOPATH = Join-Path $Root ".gopath"

if (-not (Get-Command wails -ErrorAction SilentlyContinue)) {
  Write-Host "Wails is not installed or is not on PATH."
  Write-Host "Install it with: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
  exit 1
}

if (-not (Test-Path "frontend\node_modules")) {
  Push-Location frontend
  npm install
  Pop-Location
}

if (-not (Test-Path "frontend\node_modules\go.mod")) {
  @"
module koda-ide/frontend/node_modules

go 1.22
"@ | Set-Content -Encoding utf8 "frontend\node_modules\go.mod"
}

Push-Location frontend
npm run build
Pop-Location

wails build -s -m -nopackage -tags native_webview2loader

$Exe = Join-Path $Root "build\bin\Koda Studio.exe"
if (-not (Test-Path $Exe)) {
$Exe = Join-Path $Root "build\bin\Koda Studio.exe"
if (-not (Test-Path $Exe)) {
  $Exe = Join-Path $Root "build\bin\koda-ide.exe"
}
}
if (-not (Test-Path $Exe)) {
  Write-Host "Build finished, but the app executable was not found at:"
  Write-Host $Exe
  exit 1
}

$RepoRoot = Split-Path -Parent $Root
if (Test-Path (Join-Path $RepoRoot "stdlib\math.koda")) {
  $env:KODA_HOME = $RepoRoot
  Write-Host "Using SDK at: $RepoRoot"
}

Write-Host "Starting Koda Studio..."
$ProjectPath = $args | Where-Object { $_ -and -not $_.StartsWith("-") } | Select-Object -First 1
$StartArgs = @{}
if ($ProjectPath) {
  $ProjectPath = (Resolve-Path $ProjectPath).Path
  Write-Host "Opening project: $ProjectPath"
  $StartArgs["ArgumentList"] = $ProjectPath
}
$Process = Start-Process -FilePath $Exe -WorkingDirectory $RepoRoot @StartArgs -PassThru
Start-Sleep -Seconds 3
if ($Process.HasExited) {
  Write-Host ""
  Write-Host "Koda Studio started, then closed immediately."
  Write-Host "On Windows this usually means the Microsoft Edge WebView2 runtime is busy or unhealthy."
  Write-Host "Close other WebView2-based apps, then run this script again."
  Write-Host "If it still closes, reinstall/repair Microsoft Edge WebView2 Runtime."
}
