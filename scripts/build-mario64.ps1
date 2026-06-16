param(
    [string]$Source = "examples/games/mario64_style.koda",
    [string]$Output = "mario64_style.exe"
)

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $repoRoot

. (Join-Path $PSScriptRoot "use-raylib-env.ps1") -RepoRoot $repoRoot

& ".\koda.exe" build --no-opt $Source -o $Output
if ($LASTEXITCODE -ne 0) {
    throw "Build failed for $Source"
}

$raylibRoot = & (Join-Path $PSScriptRoot "resolve-raylib-stage.ps1") -RepoRoot $repoRoot
$raylibDll = Join-Path $raylibRoot "lib\raylib.dll"
if (Test-Path $raylibDll) {
    Copy-Item -Force $raylibDll (Join-Path $repoRoot "raylib.dll")
}

Write-Output "Built $Output"
