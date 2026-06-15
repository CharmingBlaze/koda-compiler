# Install Koda on Windows (add SDK folder to PATH)
# Run from the root of an extracted SDK zip (where koda.exe lives).
$ErrorActionPreference = "Stop"

$here = $PSScriptRoot
if (Test-Path (Join-Path $here "koda.exe")) {
    $sdkRoot = $here
} elseif (Test-Path (Join-Path (Split-Path $here) "koda.exe")) {
    $sdkRoot = Split-Path $here
} else {
    $sdkRoot = (Get-Location).Path
    if (-not (Test-Path (Join-Path $sdkRoot "koda.exe"))) {
        Write-Error "Run this script from your Koda SDK folder (must contain koda.exe)."
    }
}

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$normRoot = $sdkRoot.TrimEnd('\')
if ($userPath -split ';' | Where-Object { $_.TrimEnd('\') -eq $normRoot }) {
    Write-Host "Already on PATH: $normRoot"
    exit 0
}

$newPath = if ([string]::IsNullOrWhiteSpace($userPath)) { $normRoot } else { "$userPath;$normRoot" }
[Environment]::SetEnvironmentVariable("Path", $newPath, "User")
$env:Path = "$env:Path;$normRoot"

Write-Host ""
Write-Host "Koda SDK added to your user PATH:"
Write-Host "  $normRoot"
Write-Host ""
Write-Host "Open a NEW terminal, then run:"
Write-Host "  koda doctor"
Write-Host "  koda new mygame --template graphics"
Write-Host ""
