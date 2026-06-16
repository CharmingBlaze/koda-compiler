param(
    [string]$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
)

$raylibRoot = & (Join-Path $PSScriptRoot "resolve-raylib-stage.ps1") -RepoRoot $RepoRoot
if ($LASTEXITCODE -ne 0) { throw "Raylib stage resolution failed" }

$raylibInclude = Join-Path $raylibRoot "include"
$raylibLib = Join-Path $raylibRoot "lib"
$shim = Join-Path $RepoRoot "wrappers\raylib_shim\wrapper.c"
$clangShim = Join-Path $RepoRoot "scripts\clang-gnu.cmd"
$llcPath = "C:\Program Files\LLVM\bin\llc.exe"

if (!(Test-Path $shim)) {
    throw "Raylib shim not found at '$shim'."
}

$env:KODA_CLANG = $clangShim
$env:KODA_LLC = $llcPath
$env:KODA_NATIVE_SOURCES = $shim
$env:KODA_LINKFLAGS = "-I$raylibInclude -L$raylibLib -lraylib -lopengl32 -lgdi32 -lwinmm"

Write-Output "Raylib stage: $raylibRoot"
Write-Output "KODA_CLANG=$env:KODA_CLANG"
Write-Output "KODA_LLC=$env:KODA_LLC"
Write-Output "KODA_NATIVE_SOURCES=$env:KODA_NATIVE_SOURCES"
Write-Output "KODA_LINKFLAGS=$env:KODA_LINKFLAGS"
