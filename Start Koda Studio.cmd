@echo off
REM Do not run "go run" from koda-ide — Wails needs "wails build".
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0koda-ide\run-koda-studio.ps1" %*
