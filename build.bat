@echo off
echo Building Nyumbani Care Healthcare Platform...

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

REM Navigate to project directory
cd /d "%~dp0"

REM Download dependencies
echo Downloading dependencies...
go mod download
if %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to download dependencies
    pause
    exit /b 1
)

REM Tidy up dependencies
echo Tidying dependencies...
go mod tidy
if %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to tidy dependencies
    pause
    exit /b 1
)

REM Build the application
echo Building application...
go build -o nyumbanicare.exe cmd/api/main.go
if %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to build application
    pause
    exit /b 1
)

echo Build successful! 
echo.
echo To run the application:
echo   1. Ensure your database is running
echo   2. Configure .env file with your settings
echo   3. Run: nyumbanicare.exe
echo.
echo Or run directly with: go run cmd/api/main.go
pause
