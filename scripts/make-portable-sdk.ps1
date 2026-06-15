# One-shot: build release koda + Koda Studio, assemble portable SDK folder (and optional zip).
#
# Usage (from repo root):
#   powershell -File scripts\make-portable-sdk.ps1
#   powershell -File scripts\make-portable-sdk.ps1 -Zip
#
# Output: dist\koda-<version>-sdk-windows-amd64\  (and .zip with -Zip)

param(
    [switch]$Zip,
    [switch]$SkipRelease,
    [switch]$SkipStudio
)

$ErrorActionPreference = "Stop"
$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $RepoRoot

$kodaExe = Join-Path $RepoRoot "koda-release.exe"
$wrapExe = Join-Path $RepoRoot "kodawrap.exe"

if (-not $SkipRelease) {
    Write-Host "==> Building release koda + kodawrap ..."
    & (Join-Path $PSScriptRoot "build-release.ps1")
}

if (-not (Test-Path $kodaExe) -or -not (Test-Path $wrapExe)) {
    throw "Missing koda-release.exe or kodawrap.exe — run without -SkipRelease or build manually."
}

if (-not $SkipStudio) {
    Write-Host "==> Building Koda Studio ..."
    Push-Location (Join-Path $RepoRoot "koda-ide\frontend")
    if (-not (Test-Path "node_modules")) {
        npm install
    }
    npm run build
    Pop-Location

    Push-Location (Join-Path $RepoRoot "koda-ide")
    $wails = Get-Command wails -ErrorAction SilentlyContinue
    if (-not $wails) {
        Write-Host "Installing Wails CLI ..."
        go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0
    }
    wails build -s -m -nopackage -tags native_webview2loader
    Pop-Location
}

$studioCandidates = @(
    (Join-Path $RepoRoot "koda-ide\build\bin\Koda Studio.exe"),
    (Join-Path $RepoRoot "koda-ide\build\bin\koda-ide.exe")
)
$studioExe = $studioCandidates | Where-Object { Test-Path $_ } | Select-Object -First 1

$assembleArgs = @{
    KodaExe    = $kodaExe
    KodawrapExe = $wrapExe
    IncludeStudio = (-not $SkipStudio)
}
if ($studioExe) {
    $assembleArgs["StudioExe"] = $studioExe
}
if ($Zip) {
    $assembleArgs["Zip"] = $true
}

& (Join-Path $PSScriptRoot "assemble-offline-sdk.ps1") @assembleArgs

Write-Host ""
Write-Host "Portable SDK ready. Double-click 'Start Koda Studio.bat' in the dist folder."
