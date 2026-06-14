param(
    [string]$OutputDir = "dist\koda-mario64-offline"
)

$ErrorActionPreference = "Stop"
$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $repoRoot

Write-Host "Preparing offline Koda + Raylib + Mario64 package..."

# 1) Build offline compiler/wrapper bundle.
powershell -ExecutionPolicy Bypass -File (Join-Path $PSScriptRoot "build-release.ps1")

# 2) Build mario64 playable exe (dev build path with raylib shim).
powershell -ExecutionPolicy Bypass -File (Join-Path $PSScriptRoot "build-mario64.ps1")

$outAbs = (Join-Path $repoRoot $OutputDir)
New-Item -ItemType Directory -Force $outAbs | Out-Null

# Clean prior package contents.
if (Test-Path $outAbs) {
    Get-ChildItem -Force $outAbs | Remove-Item -Recurse -Force
}

# Folder layout
$binDir = Join-Path $outAbs "bin"
$compilerDir = Join-Path $outAbs "compiler"
$raylibDir = Join-Path $outAbs "raylib"
$raylibLibDir = Join-Path $raylibDir "lib"
$raylibIncludeDir = Join-Path $raylibDir "include"
$wrappersDir = Join-Path $outAbs "wrappers"
$examplesDir = Join-Path $outAbs "examples\games"
$scriptsDir = Join-Path $outAbs "scripts"

New-Item -ItemType Directory -Force $binDir,$compilerDir,$raylibDir,$raylibLibDir,$raylibIncludeDir,$wrappersDir,$examplesDir,$scriptsDir | Out-Null

# 3) Copy compiler + wrappers from offline bundle (see scripts/build-release.ps1 layout).
$offlineBundle = Join-Path $repoRoot "dist\offline-release\bin"
Copy-Item -Force (Join-Path $offlineBundle "koda.exe") (Join-Path $compilerDir "koda.exe")
Copy-Item -Force (Join-Path $offlineBundle "kodawrap.exe") (Join-Path $compilerDir "kodawrap.exe")
Copy-Item -Recurse -Force (Join-Path $offlineBundle "stdlib") (Join-Path $compilerDir "stdlib")
if (Test-Path (Join-Path $offlineBundle "runtime")) {
    Copy-Item -Recurse -Force (Join-Path $offlineBundle "runtime") (Join-Path $compilerDir "runtime")
} else {
    Copy-Item -Recurse -Force (Join-Path $repoRoot "runtime") (Join-Path $compilerDir "runtime")
}

# 4) Copy raylib SDK bits and DLL (no download required).
$raylibRoot = Join-Path $repoRoot "raylib_lib\raylib-5.0_win64_mingw-w64"
$raylibDll = Join-Path $raylibRoot "lib\raylib.dll"
if (!(Test-Path $raylibDll)) {
    $raylibZip = Join-Path $repoRoot "raylib_win.zip"
    if (!(Test-Path $raylibZip)) {
        throw "Missing raylib.dll and raylib_win.zip. Cannot ship package without raylib.dll."
    }
    Write-Host "raylib.dll not found in raylib_lib; extracting from raylib_win.zip..."
    $extractDir = Join-Path $repoRoot ".tmp-raylib-extract"
    if (Test-Path $extractDir) {
        Remove-Item -Recurse -Force $extractDir
    }
    New-Item -ItemType Directory -Force $extractDir | Out-Null
    Expand-Archive -Path $raylibZip -DestinationPath $extractDir -Force
    $extractedDll = Join-Path $extractDir "raylib-5.0_win64_mingw-w64\lib\raylib.dll"
    if (!(Test-Path $extractedDll)) {
        throw "raylib.dll not found inside raylib_win.zip"
    }
    Copy-Item -Force $extractedDll $raylibDll
}

Copy-Item -Force (Join-Path $raylibRoot "lib\libraylib.a") (Join-Path $raylibLibDir "libraylib.a")
Copy-Item -Force $raylibDll (Join-Path $binDir "raylib.dll")
Copy-Item -Force $raylibDll (Join-Path $raylibLibDir "raylib.dll")
Copy-Item -Force (Join-Path $raylibRoot "include\raylib.h") (Join-Path $raylibIncludeDir "raylib.h")
Copy-Item -Force (Join-Path $raylibRoot "include\raymath.h") (Join-Path $raylibIncludeDir "raymath.h")
Copy-Item -Force (Join-Path $raylibRoot "include\rlgl.h") (Join-Path $raylibIncludeDir "rlgl.h")

# 5) Copy wrappers + game sources.
Copy-Item -Recurse -Force (Join-Path $repoRoot "wrappers\raylib_shim") (Join-Path $wrappersDir "raylib_shim")
Copy-Item -Force (Join-Path $repoRoot "examples\games\mario64_style.koda") (Join-Path $examplesDir "mario64_style.koda")

# 6) Copy prebuilt game exe and launch script.
Copy-Item -Force (Join-Path $repoRoot "mario64_style.exe") (Join-Path $binDir "mario64_style.exe")

$runBat = @"
@echo off
setlocal
cd /d "%~dp0\..\bin"
start "" "mario64_style.exe"
"@
Set-Content -Path (Join-Path $scriptsDir "play-mario64.bat") -Value $runBat -Encoding ASCII

# 7) Add optional compile script using packaged compiler/raylib paths.
$buildBat = @"
@echo off
setlocal
set ROOT=%~dp0\..
set KODA_NATIVE_SOURCES=%ROOT%\wrappers\raylib_shim\wrapper.c
set KODA_LINKFLAGS=-I%ROOT%\raylib\include -L%ROOT%\raylib\lib -lraylib -lopengl32 -lgdi32 -lwinmm
"%ROOT%\compiler\koda.exe" build "%ROOT%\examples\games\mario64_style.koda" -o "%ROOT%\bin\mario64_style.exe"
if errorlevel 1 exit /b 1
copy /Y "%ROOT%\raylib\lib\raylib.dll" "%ROOT%\bin\raylib.dll" >nul
echo Built %ROOT%\bin\mario64_style.exe
"@
Set-Content -Path (Join-Path $scriptsDir "build-mario64.bat") -Value $buildBat -Encoding ASCII

$readme = @"
Koda Mario64 Offline Package
============================

Everything needed is included in this folder:
- compiler\koda.exe
- compiler\kodawrap.exe
- raylib SDK files (include + static lib + raylib.dll)
- wrappers\raylib_shim
- examples\games\mario64_style.koda
- bin\mario64_style.exe (playable build)

Play now:
- run scripts\play-mario64.bat

Rebuild locally (offline):
- run scripts\build-mario64.bat
"@
Set-Content -Path (Join-Path $outAbs "README.txt") -Value $readme -Encoding ASCII

Write-Host "Done:"
Write-Host "  $outAbs"
