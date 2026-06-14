# scripts/push-release-tag.ps1
# After CHANGELOG + version bumps land on main, create an annotated tag and push it so
# .github/workflows/release.yml publishes binaries (tag must be v*).
#
# Usage (from repo root, clean tree): powershell -File scripts/push-release-tag.ps1 -Tag v0.3.0
# Then: git push origin main && git push origin v0.3.0

param(
    [Parameter(Mandatory = $true)]
    [string]$Tag
)

$ErrorActionPreference = "Stop"
$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $RepoRoot

if ($Tag -notmatch '^v\d') {
    throw "Tag must look like v0.3.0 (release.yml only runs on v*)"
}

$ver = $Tag.TrimStart('v')
$kodaMain = Join-Path $RepoRoot "cmd\koda\main.go"
$mainText = Get-Content -Raw $kodaMain
if ($mainText -notmatch "var version = `"$([regex]::Escape($ver))`"") {
    throw "cmd/koda/main.go var version must be `"$ver`" before tagging $Tag (see docs/releasing.md)"
}
$wrap = Join-Path $RepoRoot "cmd\wrapgen\wrapgen_version.go"
$wrapText = Get-Content -Raw $wrap
if ($wrapText -notmatch "WrapgenVersion = `"$([regex]::Escape($ver))`"") {
    throw "cmd/wrapgen/wrapgen_version.go WrapgenVersion must be `"$ver`" before tagging $Tag"
}

git rev-parse --is-inside-work-tree 2>$null | Out-Null
if ($LASTEXITCODE -ne 0) {
    throw "Not a git repository. Clone https://github.com/CharmingBlaze/koda-compiler.git, apply your commits, then run this script from the clone root."
}

Write-Host "Running go vet ./..."
go vet ./...
if ($LASTEXITCODE -ne 0) { throw "go vet failed" }

if (git tag -l $Tag) {
    throw "Tag $Tag already exists locally. Delete it first (git tag -d $Tag) if you need to recreate."
}

Write-Host "Creating annotated tag $Tag"
git tag -a $Tag -m "Release $Tag"

Write-Host ""
Write-Host "Tag created locally. Publish the release by pushing:"
Write-Host "  git push origin main"
Write-Host "  git push origin $Tag"
Write-Host ""
Write-Host "GitHub Actions (release.yml) will attach koda/kodawrap per OS and SDK zips."
