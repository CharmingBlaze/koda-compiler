# Write double-click launchers into a portable Koda SDK folder (Windows).
param(
    [Parameter(Mandatory = $true)]
    [string]$SdkDir
)

$ErrorActionPreference = "Stop"
$SdkDir = (Resolve-Path $SdkDir).Path
$Launchers = Join-Path $PSScriptRoot "launchers"

Copy-Item -Force (Join-Path $Launchers "start-koda-studio.bat") (Join-Path $SdkDir "Start Koda Studio.bat")
Copy-Item -Force (Join-Path $Launchers "open-pong-example.bat") (Join-Path $SdkDir "Open Pong Example.bat")
Copy-Item -Force (Join-Path $Launchers "check-sdk.bat") (Join-Path $SdkDir "Check SDK (doctor).bat")

# Unix/mac launchers ship in the same zip for cross-platform docs / WSL users.
$sh = Join-Path $PSScriptRoot "write-portable-launchers.sh"
if (Test-Path $sh) {
  bash $sh $SdkDir unix 2>$null
  bash $sh $SdkDir macos 2>$null
}

Write-Host "Wrote launchers in $SdkDir"
