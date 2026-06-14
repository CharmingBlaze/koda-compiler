# Assemble a self-contained Koda SDK folder (and optional .zip) matching GitHub Release SDK layout
# for Windows amd64: koda.exe, kodawrap.exe, stdlib/, docs/, root *.md, wrappers/, examples/,
# third_party/raylib_static/stage (raylib 5.0 headers + libraylib.a + raylib.dll).
#
# Prerequisites:
#   • Release builds of koda and kodawrap (embedded Clang + llc + lld + runtime). Build with:
#       powershell -File scripts/build-release.ps1
#     then pass -KodaExe and -KodawrapExe, or use -UseBuildRelease to run that script first.
#
# Usage (from repo root):
#   powershell -File scripts/assemble-offline-sdk.ps1 -UseBuildRelease
#   powershell -File scripts/assemble-offline-sdk.ps1 -KodaExe .\koda-release.exe -KodawrapExe .\kodawrap.exe -Zip
#
# GitHub: push a tag v* — .github/workflows/release.yml builds all platforms and runs
# scripts/vendor-raylib-stage.sh + scripts/package-release-sdk.sh (Linux). This script is the
# Windows-local equivalent for testing or ad-hoc distribution.

param(
    [string]$KodaExe = "",
    [string]$KodawrapExe = "",
    [switch]$UseBuildRelease,
    [string]$Version = "",
    [string]$OutputRoot = "dist",
    [switch]$SkipRaylib,
    [string]$RaylibZipPath = "",
    [switch]$Zip
)

$ErrorActionPreference = "Stop"
$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $RepoRoot

function Get-KodaVersionFromSource {
    $mainGo = Join-Path $RepoRoot "cmd\koda\main.go"
    $t = Get-Content -Raw $mainGo
    $m = [regex]::Match($t, 'var version = "([^"]+)"')
    if (-not $m.Success) { throw "Could not parse version from $mainGo" }
    return $m.Groups[1].Value
}

if ([string]::IsNullOrWhiteSpace($Version)) {
    $Version = Get-KodaVersionFromSource
}

if ($UseBuildRelease) {
    Write-Host "Running scripts/build-release.ps1 ..."
    & (Join-Path $PSScriptRoot "build-release.ps1")
    $KodaExe = Join-Path $RepoRoot "koda-release.exe"
    $KodawrapExe = Join-Path $RepoRoot "kodawrap.exe"
}

if ([string]::IsNullOrWhiteSpace($KodaExe) -or [string]::IsNullOrWhiteSpace($KodawrapExe)) {
    throw "Specify -KodaExe and -KodawrapExe (release builds), or -UseBuildRelease."
}

$KodaExe = (Resolve-Path $KodaExe).Path
$KodawrapExe = (Resolve-Path $KodawrapExe).Path

$folderName = "koda-$Version-sdk-windows-amd64"
$stageParent = Join-Path $RepoRoot $OutputRoot
$outDir = Join-Path $stageParent $folderName

if (Test-Path $outDir) {
    Remove-Item -Recurse -Force $outDir
}
New-Item -ItemType Directory -Force -Path $outDir | Out-Null

Copy-Item -Force $KodaExe (Join-Path $outDir "koda.exe")
Copy-Item -Force $KodawrapExe (Join-Path $outDir "kodawrap.exe")

foreach ($d in @("stdlib", "docs", "wrappers", "examples")) {
    $src = Join-Path $RepoRoot $d
    if (Test-Path $src) {
        Copy-Item -Recurse -Force $src (Join-Path $outDir $d)
    }
}

Get-ChildItem -Path $RepoRoot -Filter "*.md" -File | ForEach-Object {
    Copy-Item -Force $_.FullName (Join-Path $outDir $_.Name)
}

$raylibStage = Join-Path $outDir "third_party\raylib_static\stage"
if (-not $SkipRaylib) {
    New-Item -ItemType Directory -Force -Path $raylibStage | Out-Null
    $zipLocal = $false
    $zipFile = $null
    $downloadTmp = $null
    if (-not [string]::IsNullOrWhiteSpace($RaylibZipPath)) {
        $zipFile = (Resolve-Path $RaylibZipPath).Path
        $zipLocal = $true
    }
    if (-not $zipLocal) {
        $raylibVer = "5.0"
        $url = "https://github.com/raysan5/raylib/releases/download/$raylibVer/raylib-${raylibVer}_win64_mingw-w64.zip"
        $downloadTmp = Join-Path ([System.IO.Path]::GetTempPath()) ("koda-raylib-" + [Guid]::NewGuid().ToString("n"))
        New-Item -ItemType Directory -Force -Path $downloadTmp | Out-Null
        $zipFile = Join-Path $downloadTmp "raylib.zip"
        Write-Host "Downloading raylib $raylibVer Windows prebuild..."
        try {
            [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
            Invoke-WebRequest -Uri $url -OutFile $zipFile -UseBasicParsing
        } catch {
            Remove-Item -Recurse -Force $downloadTmp -ErrorAction SilentlyContinue
            throw "Raylib download failed: $($_.Exception.Message). Use -RaylibZipPath or -SkipRaylib."
        }
    }

    $extract = Join-Path ([System.IO.Path]::GetTempPath()) ("koda-raylib-extract-" + [Guid]::NewGuid().ToString("n"))
    New-Item -ItemType Directory -Force -Path $extract | Out-Null
    try {
        Expand-Archive -LiteralPath $zipFile -DestinationPath $extract -Force
        $inner = $null
        Get-ChildItem -Path $extract -Directory | ForEach-Object {
            $inc = Join-Path $_.FullName "include"
            $lib = Join-Path $_.FullName "lib"
            if ((Test-Path $inc) -and (Test-Path $lib)) { $inner = $_.FullName }
        }
        if (-not $inner) {
            throw "Raylib zip layout unexpected (expected one top folder with include/ and lib/)."
        }
        Copy-Item -Recurse -Force (Join-Path $inner "include") (Join-Path $raylibStage "include")
        Copy-Item -Recurse -Force (Join-Path $inner "lib") (Join-Path $raylibStage "lib")
        foreach ($x in @("LICENSE", "CHANGELOG", "README.md")) {
            $p = Join-Path $inner $x
            if (Test-Path $p) { Copy-Item -Force $p (Join-Path $raylibStage $x) }
        }
    } finally {
        Remove-Item -Recurse -Force $extract -ErrorAction SilentlyContinue
        if ($downloadTmp) {
            Remove-Item -Recurse -Force $downloadTmp -ErrorAction SilentlyContinue
        }
    }

    $need = @(
        (Join-Path $raylibStage "include\raylib.h"),
        (Join-Path $raylibStage "lib\libraylib.a"),
        (Join-Path $raylibStage "lib\raylib.dll")
    )
    foreach ($n in $need) {
        if (-not (Test-Path $n)) { throw "Raylib stage incomplete after extract: missing $n" }
    }
    Write-Host "Raylib stage OK under third_party\raylib_static\stage"
} else {
    Write-Host "Skipped raylib vendoring (-SkipRaylib). Set KODA_RAYLIB_STAGE or add stage/ manually for raylib builds."
}

$sdkReadme = @"
Koda SDK $Version (windows-amd64)
==============================

Offline layout: compiler, kodawrap (C header -> .koda + wrapper.c), stdlib, docs, examples,
wrappers (raylib + glue), and (unless -SkipRaylib) third_party/raylib_static/stage with raylib 5.0.

End users: koda does not download LLVM, Raylib, or anything else at compile time. Embedded Clang/llc/runtime unpack to a local temp directory on first use only.

Use
---
  Keep this folder together so stdlib sits next to koda.exe:

    .\koda.exe version
    .\koda.exe run examples\games\brick_breaker.koda

  Raylib static linking: when third_party/raylib_static/stage exists next to koda.exe, the
  compiler adds include + libraylib.a automatically (see docs/guides/raylib.md).

  For distribution of YOUR game: use koda bundle, and on Windows copy raylib.dll next to the
  .exe if your link uses the DLL (official prebuild ships both .a and .dll).

Wrapper tool
------------
  .\kodawrap.exe -help
  .\koda.exe wrap -help

GitHub releases
---------------
  Tag v* on the default branch to run CI: per-OS SDK zips are attached automatically.
  This folder matches that zip contents for Windows amd64.

"@

Set-Content -Path (Join-Path $outDir "SDK_README.txt") -Value $sdkReadme -Encoding UTF8

Write-Host ""
Write-Host "Assembled SDK folder:"
Write-Host "  $outDir"

if ($Zip) {
    $zipOut = Join-Path $stageParent "$folderName.zip"
    if (Test-Path $zipOut) { Remove-Item -Force $zipOut }
    Compress-Archive -Path $outDir -DestinationPath $zipOut -Force
    Write-Host "Wrote zip: $zipOut"
    Write-Host "(Upload this to a GitHub Release asset alongside binaries, or share the folder as-is.)"
}
