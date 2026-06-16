# Resolves the Raylib SDK root (include/ + lib/) for local dev scripts.
# Prefers third_party/raylib_static/stage (canonical); falls back to raylib_lib/ (legacy Windows zip).
param(
    [string]$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
)

$candidates = @(
    (Join-Path $RepoRoot "third_party\raylib_static\stage"),
    (Join-Path $RepoRoot "raylib_lib\raylib-5.0_win64_mingw-w64")
)

foreach ($root in $candidates) {
    $inc = Join-Path $root "include"
    $lib = Join-Path $root "lib"
    if ((Test-Path $inc) -and (Test-Path $lib)) {
        Write-Output $root
        exit 0
    }
}

Write-Error "Raylib stage not found. Build with: make -C third_party/raylib_static  (or place raylib under raylib_lib\)"
exit 1
