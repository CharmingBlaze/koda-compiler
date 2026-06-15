@echo off
setlocal EnableDelayedExpansion

REM Override search paths: KODA_LLVM_BIN (directory with clang.exe), KODA_MINGW_BIN (MinGW bin).
set "LLVM_BIN=%KODA_LLVM_BIN%"
set "MINGW_BIN=%KODA_MINGW_BIN%"
set "CLANG="
set "EXTRA="

if not defined MINGW_BIN (
  if exist "C:\ProgramData\mingw64\mingw64\bin\gcc.exe" set "MINGW_BIN=C:\ProgramData\mingw64\mingw64\bin"
)
if not defined MINGW_BIN (
  if exist "C:\msys64\mingw64\bin\gcc.exe" set "MINGW_BIN=C:\msys64\mingw64\bin"
)

REM Prefer llvm-mingw (self-contained clang + lld + sysroot). Works with winget LLVM-MinGW.UCRT.
for /f "delims=" %%i in ('where x86_64-w64-mingw32-clang 2^>nul') do (
  set "CLANG=%%i"
  goto :run
)
if defined MINGW_BIN if exist "%MINGW_BIN%\x86_64-w64-mingw32-clang.exe" (
  set "CLANG=%MINGW_BIN%\x86_64-w64-mingw32-clang.exe"
  goto :run
)

if not defined LLVM_BIN (
  if exist "C:\Program Files\LLVM\bin\clang.exe" set "LLVM_BIN=C:\Program Files\LLVM\bin"
)
if not defined LLVM_BIN (
  if exist "C:\Program Files (x86)\LLVM\bin\clang.exe" set "LLVM_BIN=C:\Program Files (x86)\LLVM\bin"
)
if not defined LLVM_BIN (
  for /f "delims=" %%i in ('where clang 2^>nul') do (
    set "LLVM_BIN=%%~dpi"
    goto :found_clang
  )
)
:found_clang

if not exist "%LLVM_BIN%\clang.exe" (
  echo koda: no Windows C toolchain found. Install one of: 1>&2
  echo   - llvm-mingw: winget install MartinStorsjo.LLVM-MinGW.UCRT 1>&2
  echo   - LLVM + MinGW: choco install llvm mingw, then set KODA_LLVM_BIN / KODA_MINGW_BIN 1>&2
  exit /b 1
)

if exist "%LLVM_BIN%\lld.exe" (
  if defined MINGW_BIN (
    set "PATH=%LLVM_BIN%;%MINGW_BIN%;%PATH%"
  ) else (
    set "PATH=%LLVM_BIN%;%PATH%"
  )
) else if defined MINGW_BIN (
  set "PATH=%MINGW_BIN%;%LLVM_BIN%;%PATH%"
) else (
  set "PATH=%LLVM_BIN%;%PATH%"
)

set "CLANG=%LLVM_BIN%\clang.exe"
set "EXTRA=--target=x86_64-w64-windows-gnu -fuse-ld=lld -Wno-override-module"

:run
"%CLANG%" %EXTRA% %*
