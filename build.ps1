# Nyumbani Care Healthcare Platform Build Script

Write-Host "Building Nyumbani Care Healthcare Platform..." -ForegroundColor Green

# Check if Go is installed
try {
    $goVersion = go version
    Write-Host "Found Go: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "Error: Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Please install Go from https://golang.org/dl/" -ForegroundColor Yellow
    Read-Host "Press Enter to exit"
    exit 1
}

# Navigate to project directory
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Definition
Set-Location $scriptPath

# Check if .env file exists
if (-not (Test-Path ".env")) {
    Write-Host "Warning: .env file not found. Creating from template..." -ForegroundColor Yellow
    if (Test-Path ".env.example") {
        Copy-Item ".env.example" ".env"
        Write-Host "Please configure .env file with your database settings" -ForegroundColor Yellow
    } else {
        Write-Host "Please create .env file with your configuration" -ForegroundColor Red
    }
}

# Download dependencies
Write-Host "Downloading dependencies..." -ForegroundColor Blue
try {
    go mod download
    Write-Host "Dependencies downloaded successfully" -ForegroundColor Green
} catch {
    Write-Host "Error: Failed to download dependencies" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# Tidy up dependencies
Write-Host "Tidying dependencies..." -ForegroundColor Blue
try {
    go mod tidy
    Write-Host "Dependencies tidied successfully" -ForegroundColor Green
} catch {
    Write-Host "Error: Failed to tidy dependencies" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# Build the application
Write-Host "Building application..." -ForegroundColor Blue
try {
    go build -o nyumbanicare.exe cmd/api/main.go
    Write-Host "Build successful!" -ForegroundColor Green
} catch {
    Write-Host "Error: Failed to build application" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# Generate API documentation
Write-Host "Generating API documentation..." -ForegroundColor Blue
try {
    if (-not (Test-Path "docs")) {
        New-Item -ItemType Directory -Path "docs"
    }
    go run scripts/generate-docs.go
    Write-Host "API documentation generated at docs/api.json" -ForegroundColor Green
} catch {
    Write-Host "Warning: Failed to generate API documentation" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Build completed successfully!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Ensure your PostgreSQL database is running" -ForegroundColor White
Write-Host "2. Configure .env file with your database credentials" -ForegroundColor White
Write-Host "3. Run the application:" -ForegroundColor White
Write-Host "   - Option 1: .\nyumbanicare.exe" -ForegroundColor Cyan
Write-Host "   - Option 2: go run cmd/api/main.go" -ForegroundColor Cyan
Write-Host ""
Write-Host "The API will be available at: http://localhost:8080" -ForegroundColor Green
Write-Host "Health check: http://localhost:8080/health" -ForegroundColor Green
Write-Host ""

Read-Host "Press Enter to exit"
