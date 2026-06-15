@echo off
cd /d "%~dp0"
if not exist "Koda Studio.exe" (
  echo Koda Studio.exe not found.
  pause
  exit /b 1
)
if exist "examples\games\pong" (
  start "" "%~dp0Koda Studio.exe" "%~dp0examples\games\pong"
) else (
  start "" "%~dp0Koda Studio.exe"
)
