# Generate Swagger documentation for Nyumbani Care API
# PowerShell script

# Ensure the swag CLI is installed
if (-not (Get-Command swag -ErrorAction SilentlyContinue)) {
    Write-Host "Installing swag CLI..."
    go install github.com/swaggo/swag/cmd/swag@latest
}

# Create docs directory if it doesn't exist
if (-not (Test-Path -Path "./docs")) {
    New-Item -Path "./docs" -ItemType Directory
}

# Generate Swagger documentation
Write-Host "Generating Swagger documentation..."
swag init --generalInfo cmd/api/main.go --output ./docs --parseVendor --parseDependency

if ($LASTEXITCODE -eq 0) {
    Write-Host "Swagger documentation generated successfully!" -ForegroundColor Green
    Write-Host "Access the Swagger UI at: http://localhost:8080/swagger/index.html"
} else {
    Write-Host "Failed to generate Swagger documentation." -ForegroundColor Red
}
