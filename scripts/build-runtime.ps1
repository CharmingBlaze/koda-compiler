# Build runtime/libkoda_runtime.a (MinGW or MSVC-style gcc + ar on PATH).
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$src = Join-Path $root "runtime\src"
$outLib = Join-Path $root "runtime\libkoda_runtime.a"
$objDir = Join-Path $root "runtime\obj_win"

New-Item -ItemType Directory -Force -Path $objDir | Out-Null

$cc = $env:CC
if (-not $cc) { $cc = "gcc" }
$ar = $env:AR
if (-not $ar) { $ar = "ar" }

$sources = @("value.c", "object.c", "gc.c", "koda_runtime.c")
$objs = @()
foreach ($s in $sources) {
    $objName = [System.IO.Path]::ChangeExtension($s, "o")
    $objPath = Join-Path $objDir $objName
    & $cc -O3 -fPIC -std=c11 -c (Join-Path $src $s) "-I$src" -o $objPath
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    $objs += $objPath
}

if (Test-Path $outLib) { Remove-Item $outLib }
& $ar rcs $outLib @objs
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Wrote $outLib"
