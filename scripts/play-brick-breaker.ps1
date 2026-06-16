param(
    [string]$Source = "examples/games/brick_breaker_new.koda",
    [string]$Output = "brick_breaker_new.exe"
)

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $repoRoot

. (Join-Path $PSScriptRoot "use-raylib-env.ps1") -RepoRoot $repoRoot

$raylibRoot = & (Join-Path $PSScriptRoot "resolve-raylib-stage.ps1") -RepoRoot $repoRoot
$raylibDll = Join-Path $raylibRoot "lib\raylib.dll"
if (!(Test-Path $raylibDll)) {
    throw "raylib.dll not found at '$raylibDll'"
}

& ".\koda.exe" build --no-opt $Source -o $Output
if ($LASTEXITCODE -ne 0) {
    throw "Build failed for $Source"
}

Copy-Item -Force $raylibDll (Join-Path $repoRoot "raylib.dll")
Start-Process -FilePath (Join-Path $repoRoot $Output)
Write-Output "Launched $Output"
