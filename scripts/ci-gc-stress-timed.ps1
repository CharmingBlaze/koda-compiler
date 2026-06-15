# Time-bounded koda run invocations for Windows CI (no GNU coreutils `timeout` in default PATH).
$ErrorActionPreference = "Stop"
$koda = $env:KODA_CI_BIN
if (-not $koda -or -not (Test-Path -LiteralPath $koda)) {
    throw "KODA_CI_BIN must point to an existing koda executable (got: '$koda')"
}
$repo = $env:GITHUB_WORKSPACE
if (-not $repo) { $repo = (Get-Location).Path }

function Invoke-KodaTimed([int] $Seconds, [string[]] $Arguments) {
    $p = Start-Process -FilePath $koda -ArgumentList $Arguments -WorkingDirectory $repo `
        -PassThru -NoNewWindow
    if (-not $p.WaitForExit($Seconds * 1000)) {
        try { Stop-Process -Id $p.Id -Force -ErrorAction SilentlyContinue } catch { }
        throw "timeout after ${Seconds}s: koda $($Arguments -join ' ')"
    }
    if ($p.ExitCode -ne 0) {
        throw "exit code $($p.ExitCode): koda $($Arguments -join ' ')"
    }
}

Write-Host "==> GC soak (timed)"
Invoke-KodaTimed 90 @("run", "tests/gc_pressure_expr.koda")
Invoke-KodaTimed 90 @("run", "tests/globals_perf.koda")
Invoke-KodaTimed 90 @("run", "tests/gc_soak.koda")
Invoke-KodaTimed 120 @("run", "--no-opt", "tests/nursery_test.koda")
Invoke-KodaTimed 120 @("run", "--no-opt", "tests/incremental_gc_test.koda")

Write-Host "==> Stress smoke (timed)"
Invoke-KodaTimed 90 @("run", "tests/stress/stress_mixed_alloc.koda")
Invoke-KodaTimed 90 @("run", "tests/stress/stress_deep_recursion.koda")
Invoke-KodaTimed 90 @("run", "tests/stress/stress_string_pressure.koda")
Invoke-KodaTimed 120 @("run", "--no-opt", "tests/stress/large_game_sim.koda")

Write-Host "==> Tier-1 regression (timed)"
Invoke-KodaTimed 60 @("run", "tests/struct_methods.koda")
Invoke-KodaTimed 60 @("run", "tests/integer_types.koda")
Invoke-KodaTimed 60 @("run", "tests/intern_clear_test.koda")
Invoke-KodaTimed 60 @("run", "tests/stdlib_modules_test.koda")
Invoke-KodaTimed 60 @("run", "tests/enum_exhaustive.koda")

Write-Host "==> GC / stress timed runs OK"
