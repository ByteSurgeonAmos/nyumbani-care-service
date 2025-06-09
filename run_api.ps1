# Run Nyumbani Care API
Write-Host "Starting Nyumbani Care API..." -ForegroundColor Cyan

# Check if executable exists
if (!(Test-Path -Path ".\nyumbanicare.exe")) {
    Write-Host "Error: nyumbanicare.exe not found. Please build the application first." -ForegroundColor Red
    Write-Host "Run .\build.ps1 to build the application." -ForegroundColor Yellow
    exit 1
}

# Run the application
try {
    Write-Host "API is running. Press Ctrl+C to stop." -ForegroundColor Green
    .\nyumbanicare.exe
}
catch {
    Write-Host "Error running the application: $_" -ForegroundColor Red
    exit 1
}
