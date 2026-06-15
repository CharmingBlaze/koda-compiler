@echo off
cd /d "%~dp0"
if not exist "Koda Studio.exe" (
  echo Koda Studio.exe not found in this folder.
  echo Unzip the full Koda SDK zip so Koda Studio.exe sits next to koda.exe and stdlib\
  pause
  exit /b 1
)
start "" "%~dp0Koda Studio.exe" %*
