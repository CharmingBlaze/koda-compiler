param(
    [string]$ProjectDir = "examples/games/brick-breaker",
    [string]$Output = "brick-breaker.exe"
)

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $repoRoot

. (Join-Path $PSScriptRoot "use-raylib-env.ps1") -RepoRoot $repoRoot

$raylibRoot = & (Join-Path $PSScriptRoot "resolve-raylib-stage.ps1") -RepoRoot $repoRoot
$raylibDll = Join-Path $raylibRoot "lib\raylib.dll"
if (!(Test-Path $raylibDll)) {
    throw "raylib.dll not found at '$raylibDll'"
}

Push-Location $ProjectDir
try {
    & (Join-Path $repoRoot "koda.exe") build --release src\main.koda -o $Output
    if ($LASTEXITCODE -ne 0) {
        throw "Build failed in $ProjectDir"
    }
    $built = Join-Path (Get-Location) $Output
} finally {
    Pop-Location
}

Copy-Item -Force $raylibDll (Join-Path $repoRoot "raylib.dll")
Start-Process -FilePath $built
Write-Output "Launched $built"
