# scripts/build-release.ps1
# Builds a local release `koda.exe` + offline SDK under dist\offline-release\bin\,
# then a shippable zip dist\koda-<version>-windows-amd64-offline.zip (omit with -NoZip).
# Offline-only: this script never downloads dependencies.

param(
    [string]$Platform = "windows",
    [string]$LLVMBin = "C:\Program Files\LLVM\bin",
    [string]$MinGWBin = "C:\ProgramData\mingw64\mingw64\bin",
    [string]$OutputDir = "dist\offline-release",
    [switch]$NoZip,
    [switch]$PackageSdk,
    [switch]$PackageSdkZip
)

$ErrorActionPreference = "Stop"
$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $RepoRoot

Write-Host "Building Koda release binary for $Platform..."
Write-Host "Mode: offline (no downloads)"

if ($Platform -ne "windows") {
    throw "This script currently supports Platform=windows only."
}

$clangExe = Join-Path $LLVMBin "clang.exe"
$lldExe = Join-Path $LLVMBin "lld.exe"
$clangDriver = Join-Path $RepoRoot "scripts\clang-gnu.cmd"
$mingwArExe = Join-Path $MinGWBin "ar.exe"

$llcExe = $null
foreach ($name in @("llc.exe", "llc-18.exe", "llc-14.exe")) {
    $cand = Join-Path $LLVMBin $name
    if (Test-Path $cand) { $llcExe = $cand; break }
}
if (!$llcExe) {
    $onPath = (Get-Command "llc.exe" -ErrorAction SilentlyContinue).Source
    if ($onPath) { $llcExe = $onPath }
}

if (!(Test-Path $clangExe)) { throw "Missing $clangExe" }
if (!(Test-Path $lldExe)) { throw "Missing $lldExe" }
if (!$llcExe -or !(Test-Path $llcExe)) {
    throw @"
Missing LLVM llc (LLVM IR -> object). Koda release builds must embed llc.exe next to clang.
Checked: $LLVMBin (and llc.exe on PATH).

Some Windows LLVM packages omit llc. Install a full LLVM binary distribution, for example:
  https://github.com/llvm/llvm-project/releases  (Pre-built binaries / Windows installer, all tools)
or: choco install llvm -y

Then re-run with -LLVMBin pointing at that LLVM's bin folder (e.g. 'C:\Program Files\LLVM\bin').
"@
}
if (!(Test-Path $mingwArExe)) { throw "Missing $mingwArExe" }
if (!(Test-Path $clangDriver)) { throw "Missing $clangDriver" }

$env:PATH = "$LLVMBin;$MinGWBin;$env:PATH"

$kodaMain = Join-Path $RepoRoot "cmd\koda\main.go"
$kodaMainText = Get-Content -Raw $kodaMain
$versionMatch = [regex]::Match($kodaMainText, 'var version = "([^"]+)"')
if (!$versionMatch.Success) {
    throw "Could not extract Koda version from $kodaMain"
}
$releaseVersion = $versionMatch.Groups[1].Value
Write-Host "Version: $releaseVersion"

Write-Host "Building C runtime..."
New-Item -ItemType Directory -Force runtime\obj_win | Out-Null
& cmd /c """$clangDriver"" -c runtime\src\value.c        -O2 -std=c11 -Iruntime\src -D_AMD64_=1 -o runtime\obj_win\value.o"
if ($LASTEXITCODE -ne 0) { throw "Failed compiling runtime\src\value.c" }
& cmd /c """$clangDriver"" -c runtime\src\object.c       -O2 -std=c11 -Iruntime\src -D_AMD64_=1 -o runtime\obj_win\object.o"
if ($LASTEXITCODE -ne 0) { throw "Failed compiling runtime\src\object.c" }
& cmd /c """$clangDriver"" -c runtime\src\gc.c           -O2 -std=c11 -Iruntime\src -D_AMD64_=1 -o runtime\obj_win\gc.o"
if ($LASTEXITCODE -ne 0) { throw "Failed compiling runtime\src\gc.c" }
& cmd /c """$clangDriver"" -c runtime\src\koda_runtime.c -O2 -std=c11 -Iruntime\src -D_AMD64_=1 -o runtime\obj_win\koda_runtime.o"
if ($LASTEXITCODE -ne 0) { throw "Failed compiling runtime\src\koda_runtime.c" }
& $mingwArExe rcs runtime\libkoda_runtime.a runtime\obj_win\value.o runtime\obj_win\object.o runtime\obj_win\gc.o runtime\obj_win\koda_runtime.o
if ($LASTEXITCODE -ne 0) { throw "Failed creating runtime\libkoda_runtime.a" }

Write-Host "Populating embed directory..."
New-Item -ItemType Directory -Force internal\embed\windows\amd64 | Out-Null
Copy-Item $clangExe internal\embed\windows\amd64\clang.exe
Copy-Item $lldExe   internal\embed\windows\amd64\lld.exe
Copy-Item $llcExe   internal\embed\windows\amd64\llc.exe
Copy-Item runtime\libkoda_runtime.a  internal\embed\windows\amd64\

Write-Host "Building koda.exe..."
go build -trimpath -tags release -ldflags="-s -w -X main.version=$releaseVersion" -o koda-release.exe .\cmd\koda
if ($LASTEXITCODE -ne 0) { throw "Failed building koda-release.exe" }

Write-Host "Building kodawrap (embedded Clang like koda -- no system compiler needed at runtime)..."
go build -trimpath -tags release -ldflags="-s -w -X main.WrapgenVersion=$releaseVersion" -o kodawrap.exe .\cmd\wrapgen
if ($LASTEXITCODE -ne 0) { throw "Failed building kodawrap.exe" }

Write-Host "Assembling offline distribution (everything under bin\)..."
$bundleParent = Join-Path $RepoRoot $OutputDir
$BinDir = Join-Path $bundleParent "bin"

# Remove stale layout (flat bundle or old bin\) so we never mix old and new trees.
if (Test-Path $BinDir) {
    Remove-Item -Recurse -Force $BinDir
}
foreach ($name in @("koda.exe", "kodawrap.exe", "README_OFFLINE.txt")) {
    $p = Join-Path $bundleParent $name
    if (Test-Path $p) { Remove-Item -Force $p }
}
foreach ($d in @("stdlib", "wrappers", "runtime", "docs", "examples")) {
    $p = Join-Path $bundleParent $d
    if (Test-Path $p) { Remove-Item -Recurse -Force $p }
}
foreach ($md in Get-ChildItem -Path $bundleParent -Filter "*.md" -File -ErrorAction SilentlyContinue) {
    Remove-Item -Force $md.FullName
}

New-Item -ItemType Directory -Force $BinDir | Out-Null

Copy-Item -Force .\koda-release.exe (Join-Path $BinDir "koda.exe")
Copy-Item -Force .\kodawrap.exe (Join-Path $BinDir "kodawrap.exe")

if (Test-Path .\stdlib)   { Copy-Item -Recurse -Force .\stdlib   (Join-Path $BinDir "stdlib") }
if (Test-Path .\wrappers) { Copy-Item -Recurse -Force .\wrappers  (Join-Path $BinDir "wrappers") }
if (Test-Path .\docs)     { Copy-Item -Recurse -Force .\docs      (Join-Path $BinDir "docs") }
if (Test-Path .\examples){ Copy-Item -Recurse -Force .\examples  (Join-Path $BinDir "examples") }

# koda build links with -I <cwd>/runtime/src ΓÇö cwd must contain this tree (see internal/nativebuild/build.go).
$rtSrc = Join-Path $RepoRoot "runtime\src"
if (Test-Path $rtSrc) {
    $dstRt = Join-Path $BinDir "runtime\src"
    New-Item -ItemType Directory -Force $dstRt | Out-Null
    Copy-Item -Force (Join-Path $rtSrc "*") $dstRt
} else {
    throw "Missing runtime\src (required in the bundle for koda build / koda run)"
}

# MinGW-w64 runtime DLLs often needed next to game .exe when linking with the same GNU toolchain.
foreach ($dll in @("libgcc_s_seh-1.dll", "libstdc++-6.dll", "libwinpthread-1.dll")) {
    $p = Join-Path $MinGWBin $dll
    if (Test-Path $p) {
        Copy-Item -Force $p (Join-Path $BinDir $dll)
    }
}

$samplesDir = Join-Path $BinDir "samples"
New-Item -ItemType Directory -Force $samplesDir | Out-Null
$helloTest = Join-Path $RepoRoot "tests\hello.koda"
if (Test-Path $helloTest) {
    Copy-Item -Force $helloTest (Join-Path $samplesDir "hello.koda")
}

foreach ($md in Get-ChildItem -Path $RepoRoot -Filter "*.md" -File) {
    Copy-Item -Force $md.FullName (Join-Path $BinDir $md.Name)
}

$raylibStageRel = Join-Path "third_party" "raylib_static\stage"
$stageDst = Join-Path $BinDir $raylibStageRel
$stageSrcCandidates = @(
    (Join-Path $RepoRoot "third_party\raylib_static\stage"),
    (Join-Path $RepoRoot "raylib_lib\raylib-5.0_win64_mingw-w64")
)
$stageSrc = $null
foreach ($c in $stageSrcCandidates) {
    $inc = Join-Path $c "include"
    $libray = Join-Path $c "lib\libraylib.a"
    if ((Test-Path $inc) -and (Test-Path $libray)) { $stageSrc = $c; break }
}
if (!$stageSrc) {
    throw "Raylib 5.0 Windows stage not found. Expected include/ + lib/libraylib.a under third_party\raylib_static\stage or raylib_lib\raylib-5.0_win64_mingw-w64"
}
New-Item -ItemType Directory -Force $stageDst | Out-Null
Copy-Item -Recurse -Force (Join-Path $stageSrc "include") (Join-Path $stageDst "include")
Copy-Item -Recurse -Force (Join-Path $stageSrc "lib")     (Join-Path $stageDst "lib")
foreach ($x in @("LICENSE", "CHANGELOG", "README.md")) {
    $p = Join-Path $stageSrc $x
    if (Test-Path $p) { Copy-Item -Force $p $stageDst }
}

$raylibDll = Join-Path $stageDst "lib\raylib.dll"
if (Test-Path $raylibDll) {
    Copy-Item -Force $raylibDll (Join-Path $BinDir "raylib.dll")
} else {
    Write-Warning "raylib.dll not found under stage lib\; Windows games may need the DLL next to the .exe"
}

# Upstream raylib 6.x C headers (unwrapped ΓÇö no Koda wrapper); same files as stdlib/sys-include for kodawrap / reference.
$unwrapRoot = Join-Path $BinDir "third_party\raylib6_unwrapped"
$unwrapInc = Join-Path $unwrapRoot "include"
$sysInc = Join-Path $RepoRoot "stdlib\sys-include"
New-Item -ItemType Directory -Force $unwrapInc | Out-Null
if (Test-Path $sysInc) {
    Get-ChildItem -Path $sysInc -Filter "*.h" -File | ForEach-Object {
        Copy-Item -Force $_.FullName (Join-Path $unwrapInc $_.Name)
    }
}
@"
Upstream Raylib 6.x C headers (unwrapped)
==========================================
These are the same headers vendored in stdlib/sys-include/ (raylib 6 API, no Koda bindings).

Use with kodawrap, for C reference, or custom KODA_NATIVE_SOURCES builds, for example:
  kodawrap -name myraylib -headers .\third_party\raylib6_unwrapped\include\raylib.h ...

Note: third_party/raylib_static/stage/ still ships the official raylib 5.0 Windows prebuild
(libraylib.a, raylib.dll). Koda auto-vendored linking targets that tree unless you override
KODA_RAYLIB_STAGE / KODA_USE_VENDORED_RAYLIB. Mixing 6.0 headers with 5.0 libs can break; use
matching raylib 6 binaries if you compile against these headers.
"@ | Set-Content -Path (Join-Path $unwrapRoot "README.txt")

@"
Koda offline SDK (Windows)
----------------------------
No separate install: koda.exe and kodawrap.exe embed Clang, llc, lld, and the Koda runtime library.
You do not install LLVM, MinGW, Go, or Visual Studio on this PC to compile .koda programs ΓÇö only **koda.exe**, **kodawrap.exe**, and this folder (stdlib next to them; write access for first-run unpack).

All tools and libraries live in this folder (bin). Use it as your working directory or add it to PATH.

Contents
--------
  koda.exe, kodawrap.exe  — compiler and wrapper tool
  raylib.dll                            ΓÇö copy next to your game .exe when linking dynamically
  stdlib/, wrappers/, docs/, examples/   ΓÇö language + Raylib docs and samples
  samples/hello.koda                     ΓÇö tiny program to verify the toolchain after unzip
  runtime/src/                          ΓÇö C headers for the native link step (run `koda` with cwd = this bin folder)
  libgcc_s_seh-1.dll, ΓÇª                  ΓÇö optional MinGW runtimes (when present on the build machine)
  third_party/raylib_static/stage/      ΓÇö Raylib 5.0 headers + libraylib.a (+ lib\raylib.dll)
  third_party/raylib6_unwrapped/       ΓÇö upstream Raylib 6.x C headers (include/) + README.txt
  *.md                                  ΓÇö repo root guides (language.md, README, ΓÇª)

Raylib: see docs/guides/raylib.md. Vendored stage is detected next to koda.exe automatically.
Unwrapped Raylib 6 headers: third_party/raylib6_unwrapped/include/ (see README there).

If native compile fails, run: .\koda.exe doctor
"@ | Set-Content -Path (Join-Path $BinDir "README_OFFLINE.txt")

@"
Koda Windows offline bundle
No system compiler install required ΓÇö open **bin** and run **koda.exe** (toolchain is embedded). See START_HERE.txt.
"@ | Set-Content -Path (Join-Path $bundleParent "README.txt")

@"
Koda ΓÇö ship / unzip / run
==========================

1) Unzip anywhere (a normal folder you can write to is best). **Install nothing else** ΓÇö not Go, not LLVM, not a C compiler. Only use **koda.exe** and **kodawrap.exe** from this SDK (with **stdlib/** beside them). First `koda build` / `koda run` may unpack embedded tools once.

2) Open a terminal in the **bin** folder (this is your "project root" for builds):
     cd bin

3) Smoke test:
     .\koda.exe version
     .\koda.exe doctor

4) Quick build test (from **bin**, no repo checkout needed):
     .\koda.exe build samples\hello.koda -o hello.exe
     .\hello.exe

   Build another program (paths can be absolute, or relative to **bin**):
     .\koda.exe build ..\path\to\yourgame.koda -o mygame.exe

   For the Raylib shim demo, set native glue once per session, then run (see README_OFFLINE.txt and docs\guides\raylib.md):
     `$env:KODA_NATIVE_SOURCES = "wrappers\raylib_shim\wrapper.c"
     .\koda.exe run examples\raylib_shim_demo.koda

5) Ship your game: copy **raylib.dll** (and, if needed, **libgcc_s_seh-1.dll**, **libstdc++-6.dll**, **libwinpthread-1.dll**) next to your **.exe** when you use the Raylib DLL or MinGW-linked output.

Full layout notes: README_OFFLINE.txt (inside bin).
"@ | Set-Content -Path (Join-Path $bundleParent "START_HERE.txt")

if (!$NoZip) {
    $zipFolderName = "koda-$releaseVersion-windows-amd64-offline"
    $zipStageRoot = Join-Path $RepoRoot "dist\_ship_zip_stage"
    if (Test-Path $zipStageRoot) {
        Remove-Item -Recurse -Force $zipStageRoot
    }
    $zipInner = Join-Path $zipStageRoot $zipFolderName
    New-Item -ItemType Directory -Force $zipInner | Out-Null
    Get-ChildItem -Path $bundleParent -Force | ForEach-Object {
        Copy-Item -Recurse -Force $_.FullName (Join-Path $zipInner $_.Name)
    }
    $zipPath = Join-Path $RepoRoot "dist\$zipFolderName.zip"
    if (Test-Path $zipPath) {
        Remove-Item -Force $zipPath
    }
    Write-Host "Writing shippable zip: $zipPath"
    Compress-Archive -Path $zipInner -DestinationPath $zipPath -Force
    Remove-Item -Recurse -Force $zipStageRoot
}

Write-Host ""
Write-Host "Done. Test from bin folder:"
Write-Host "  cd `"$BinDir`""
Write-Host "  .\koda.exe build `"$RepoRoot\tests\hello.koda`" -o hello.exe"
Write-Host "  .\hello.exe"
Write-Host ""
Write-Host "Offline bundle folder:"
Write-Host "  $BinDir"
if (!$NoZip) {
    Write-Host "Ship this zip (unzip = ready):"
    Write-Host "  $(Join-Path $RepoRoot "dist\koda-$releaseVersion-windows-amd64-offline.zip")"
}

if ($PackageSdk -or $PackageSdkZip) {
    $assemble = Join-Path $PSScriptRoot "assemble-offline-sdk.ps1"
    Write-Host ""
    Write-Host "Packaging full SDK (docs, stdlib, wrappers, examples, raylib download unless -SkipRaylib)..."
    $zipArg = @()
    if ($PackageSdkZip) { $zipArg = @("-Zip") }
    & $assemble -KodaExe (Join-Path $RepoRoot "koda-release.exe") -KodawrapExe (Join-Path $RepoRoot "kodawrap.exe") -Version $releaseVersion @zipArg
}
