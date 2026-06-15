# Build runtime/libkoda_runtime.a for local development on Windows.
# Matches CI/release: prefer LLVM clang + llvm-ar with -D_AMD64_=1; falls back to gcc + ar.
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$src = Join-Path $root "runtime\src"
$outLib = Join-Path $root "runtime\libkoda_runtime.a"
$objDir = Join-Path $root "runtime\obj_win"

New-Item -ItemType Directory -Force -Path $objDir | Out-Null

function Find-LLVMTool([string] $BaseName) {
    if ($env:KODA_LLVM_BIN) {
        $p = Join-Path $env:KODA_LLVM_BIN "$BaseName.exe"
        if (Test-Path -LiteralPath $p) { return $p }
    }
    foreach ($dir in @(
        "C:\Program Files\LLVM\bin",
        "C:\Program Files (x86)\LLVM\bin"
    )) {
        $p = Join-Path $dir "$BaseName.exe"
        if (Test-Path -LiteralPath $p) { return $p }
    }
    $found = Get-Command $BaseName -ErrorAction SilentlyContinue
    if ($found) { return $found.Source }
    return $null
}

# Prefer llvm-mingw prefix compiler (matches scripts/clang-gnu.cmd link driver).
$cc = $env:CC
if (-not $cc) {
    $mingwClang = Get-Command x86_64-w64-mingw32-clang -ErrorAction SilentlyContinue
    if ($mingwClang) { $cc = $mingwClang.Source }
}
if (-not $cc) { $cc = Find-LLVMTool "clang" }
if (-not $cc) {
    $gcc = Get-Command gcc -ErrorAction SilentlyContinue
    if ($gcc) {
        $cc = $gcc.Source
    } else {
        Write-Error @"
No C compiler found. Install LLVM (e.g. choco install llvm) or MinGW gcc, then re-run.
  Optional: set CC to your compiler and AR to llvm-ar or ar.
  Optional: set KODA_LLVM_BIN to the directory containing clang.exe and llvm-ar.exe.
"@
    }
}

$ar = $env:AR
if (-not $ar) {
    $llvmAr = Find-LLVMTool "llvm-ar"
    if ($llvmAr) {
        $ar = $llvmAr
    } else {
        $arCmd = Get-Command ar -ErrorAction SilentlyContinue
        if ($arCmd) {
            $ar = $arCmd.Source
        } else {
            Write-Error "No llvm-ar or ar found. Install LLVM or set AR."
        }
    }
}

$useClang = $cc -match "clang"
$cflags = @("-O2", "-std=c11")
if ($useClang) { $cflags += "-D_AMD64_=1" }

$sources = @("value.c", "object.c", "gc.c", "koda_runtime.c")
$objs = @()
foreach ($s in $sources) {
    $objName = [System.IO.Path]::ChangeExtension($s, "o")
    $objPath = Join-Path $objDir $objName
    $srcPath = Join-Path $src $s
    & $cc @cflags -c $srcPath "-I$src" -o $objPath
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    $objs += $objPath
}

if (Test-Path -LiteralPath $outLib) { Remove-Item -LiteralPath $outLib }
& $ar rcs $outLib @objs
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Wrote $outLib (CC=$cc AR=$ar)"
