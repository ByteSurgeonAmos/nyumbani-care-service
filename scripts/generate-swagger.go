package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Get the working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// Install swag if not already installed
	cmd := exec.Command("go", "install", "github.com/swaggo/swag/cmd/swag@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to install swag: %v", err)
	}
	
	// Find the swag executable
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		// If GOPATH is not set, try to get it from go env
		getGopathCmd := exec.Command("go", "env", "GOPATH")
		gopathBytes, err := getGopathCmd.Output()
		if err != nil {
			log.Fatalf("Failed to get GOPATH: %v", err)
		}
		gopath = string(gopathBytes)
		// Trim newline
		if len(gopath) > 0 && gopath[len(gopath)-1] == '\n' {
			gopath = gopath[:len(gopath)-1]
		}
		if len(gopath) > 0 && gopath[len(gopath)-1] == '\r' {
			gopath = gopath[:len(gopath)-1]
		}
	}
	
	// Create swag command with full path to executable
	swagPath := filepath.Join(gopath, "bin", "swag")
	if isWindows() {
		swagPath += ".exe"
	}
	
	// Create docs directory if it doesn't exist
	docsDir := filepath.Join(wd, "docs")
	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		if err := os.Mkdir(docsDir, 0755); err != nil {
			log.Fatalf("Failed to create docs directory: %v", err)
		}
	}
	
	// Run swag init with full path
	swagCmd := exec.Command(
		swagPath, 
		"init", 
		"--generalInfo", filepath.Join("cmd", "api", "main.go"),
		"--output", "./docs",
		"--parseVendor",
		"--parseDependency",
	)
	
	swagCmd.Dir = wd
	swagCmd.Stdout = os.Stdout
	swagCmd.Stderr = os.Stderr
	
	fmt.Printf("Running: %s\n", swagCmd.String())
	
	if err := swagCmd.Run(); err != nil {
		log.Fatalf("Failed to generate Swagger documentation: %v", err)
	}
	
	fmt.Println("Swagger documentation generated successfully!")
	fmt.Println("Access the Swagger UI at: http://localhost:8080/swagger/index.html")
}

func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}
