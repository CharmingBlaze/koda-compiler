$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $Root

if (-not (Get-Command wails -ErrorAction SilentlyContinue)) {
  Write-Host "Install Wails: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
  exit 1
}

$RepoRoot = Split-Path -Parent $Root
if (Test-Path (Join-Path $RepoRoot "stdlib\math.koda")) {
  $env:KODA_HOME = $RepoRoot
}

Write-Host "Starting Koda Studio (wails dev — rebuilds Go on save)..."
wails dev -tags native_webview2loader
