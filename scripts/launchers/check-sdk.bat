@echo off
cd /d "%~dp0"
if not exist "koda.exe" (
  echo koda.exe not found in this folder.
  pause
  exit /b 1
)
"%~dp0koda.exe" doctor
echo.
pause
