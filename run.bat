@echo off
setlocal

set CGO_ENABLED=0
set GOEXPERIMENT=greenteagc

echo Building Kuroryuu CLI v0.5...
go build -v -o kuroryuu.exe .
if %errorlevel% neq 0 (
    echo Build failed.
    pause
    exit /b %errorlevel%
)

echo.
echo Starting Kuroryuu...
kuroryuu.exe %*
