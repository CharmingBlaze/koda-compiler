# Build koda + kodawrap (C header → .koda + wrapper.c) into ./bin — one script for Windows authors.
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

New-Item -ItemType Directory -Force -Path "bin" | Out-Null

Write-Host "Building koda..." -ForegroundColor Cyan
go build -trimpath -ldflags "-s -w" -o "bin/koda.exe" ./cmd/koda

Write-Host "Building kodawrap (cmd/wrapgen)..." -ForegroundColor Cyan
go build -trimpath -ldflags "-s -w" -o "bin/kodawrap.exe" ./cmd/wrapgen

Write-Host "Building wrapgen.exe (same tool, legacy name)..." -ForegroundColor DarkGray
go build -trimpath -ldflags "-s -w" -o "bin/wrapgen.exe" ./cmd/wrapgen

Write-Host "Building C runtime (runtime/libkoda_runtime.a)..." -ForegroundColor Cyan
& "$PSScriptRoot\build-runtime.ps1"
if ($LASTEXITCODE -ne 0) {
    Write-Host "  runtime build failed (install MinGW gcc/ar or run scripts/build-runtime.ps1 — see CONTRIBUTING.md)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Done. Add .\bin to PATH, or run:" -ForegroundColor Green
Write-Host "  .\bin\koda.exe help" -ForegroundColor Gray
Write-Host "  .\bin\kodawrap.exe -name mylib -headers .\mylib.h -out .\wrappers\mylib" -ForegroundColor Gray
