@echo off
echo Starting Mock Server...
cd /d "%~dp0"
go run main.go
pause