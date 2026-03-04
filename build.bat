@echo off
echo Building scaffold...
go mod tidy
if %errorlevel% neq 0 (
    echo Failed to download dependencies. Is Go installed?
    exit /b %errorlevel%
)
go build -o scaffold.exe .
if %errorlevel% neq 0 (
    echo Build failed.
    exit /b %errorlevel%
)
echo.
echo Build successful: scaffold.exe
echo.
echo Quick start:
echo   scaffold init
echo   scaffold list
echo   scaffold use general-dev
