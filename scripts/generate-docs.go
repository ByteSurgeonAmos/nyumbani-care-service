package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// OpenAPI 3.0 specification structure
type OpenAPISpec struct {
	OpenAPI string                 `json:"openapi"`
	Info    Info                   `json:"info"`
	Servers []Server               `json:"servers"`
	Paths   map[string]interface{} `json:"paths"`
	Components Components           `json:"components"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Contact     Contact `json:"contact"`
}

type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type Components struct {
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes"`
	Schemas         map[string]interface{}    `json:"schemas"`
}

type SecurityScheme struct {
	Type   string `json:"type"`
	Scheme string `json:"scheme"`
}

func generateAPIDoc() {
	spec := OpenAPISpec{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:       "Nyumbani Care API",
			Description: "Digital-first healthcare platform for Africa providing test kit e-commerce, prescription management, lab work booking, telehealth consultations, health education, AI symptom checker, and CareSense analytics.",
			Version:     "1.0.0",
			Contact: Contact{
				Name:  "Nyumbani Care Team",
				Email: "api@nyumbanicare.com",
				URL:   "https://nyumbanicare.com",
			},
		},
		Servers: []Server{
			{
				URL:         "http://localhost:8080",
				Description: "Development server",
			},
			{
				URL:         "https://api.nyumbanicare.com",
				Description: "Production server",
			},
		},
		Paths: map[string]interface{}{
			"/health": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Health check endpoint",
					"description": "Returns the health status of the API",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "API is healthy",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"status": map[string]interface{}{
												"type": "string",
											},
											"service": map[string]interface{}{
												"type": "string",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"/api/v1/auth/register": map[string]interface{}{
				"post": map[string]interface{}{
					"summary":     "Register a new user",
					"description": "Create a new user account",
					"tags":        []string{"Authentication"},
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/RegisterRequest",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"201": map[string]interface{}{
							"description": "User registered successfully",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/AuthResponse",
									},
								},
							},
						},
						"400": map[string]interface{}{
							"description": "Bad request",
						},
						"409": map[string]interface{}{
							"description": "User already exists",
						},
					},
				},
			},
			"/api/v1/auth/login": map[string]interface{}{
				"post": map[string]interface{}{
					"summary":     "Login user",
					"description": "Authenticate user and return JWT token",
					"tags":        []string{"Authentication"},
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/LoginRequest",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Login successful",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/AuthResponse",
									},
								},
							},
						},
						"401": map[string]interface{}{
							"description": "Invalid credentials",
						},
					},
				},
			},
			"/api/v1/test-kits": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "List test kits",
					"description": "Get all available test kits",
					"tags":        []string{"Test Kits"},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "List of test kits",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "array",
										"items": map[string]interface{}{
											"$ref": "#/components/schemas/TestKit",
										},
									},
								},
							},
						},
					},
				},
			},
			"/api/v1/orders": map[string]interface{}{
				"post": map[string]interface{}{
					"summary":     "Create test kit order",
					"description": "Place an order for test kits",
					"tags":        []string{"Orders"},
					"security": []map[string]interface{}{
						{"bearerAuth": []string{}},
					},
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/OrderRequest",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"201": map[string]interface{}{
							"description": "Order created successfully",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/Order",
									},
								},
							},
						},
						"401": map[string]interface{}{
							"description": "Unauthorized",
						},
					},
				},
			},
		},
		Components: Components{
			SecuritySchemes: map[string]SecurityScheme{
				"bearerAuth": {
					Type:   "http",
					Scheme: "bearer",
				},
			},
			Schemas: map[string]interface{}{
				"RegisterRequest": map[string]interface{}{
					"type": "object",
					"required": []string{"email", "password", "first_name", "last_name", "phone_number", "date_of_birth", "gender", "address"},
					"properties": map[string]interface{}{
						"email": map[string]interface{}{
							"type":   "string",
							"format": "email",
						},
						"password": map[string]interface{}{
							"type":      "string",
							"minLength": 6,
						},
						"first_name": map[string]interface{}{
							"type": "string",
						},
						"last_name": map[string]interface{}{
							"type": "string",
						},
						"phone_number": map[string]interface{}{
							"type": "string",
						},
						"date_of_birth": map[string]interface{}{
							"type":   "string",
							"format": "date",
						},
						"gender": map[string]interface{}{
							"type": "string",
							"enum": []string{"male", "female", "other"},
						},
						"address": map[string]interface{}{
							"type": "string",
						},
					},
				},
				"LoginRequest": map[string]interface{}{
					"type": "object",
					"required": []string{"email", "password"},
					"properties": map[string]interface{}{
						"email": map[string]interface{}{
							"type":   "string",
							"format": "email",
						},
						"password": map[string]interface{}{
							"type": "string",
						},
					},
				},
				"AuthResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type": "string",
						},
						"user": map[string]interface{}{
							"$ref": "#/components/schemas/User",
						},
					},
				},
				"User": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":   "string",
							"format": "uuid",
						},
						"email": map[string]interface{}{
							"type": "string",
						},
						"first_name": map[string]interface{}{
							"type": "string",
						},
						"last_name": map[string]interface{}{
							"type": "string",
						},
						"role": map[string]interface{}{
							"type": "string",
						},
					},
				},
				"TestKit": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":   "string",
							"format": "uuid",
						},
						"name": map[string]interface{}{
							"type": "string",
						},
						"description": map[string]interface{}{
							"type": "string",
						},
						"price": map[string]interface{}{
							"type": "number",
						},
						"category": map[string]interface{}{
							"type": "string",
						},
						"stock_quantity": map[string]interface{}{
							"type": "integer",
						},
					},
				},
				"OrderRequest": map[string]interface{}{
					"type": "object",
					"required": []string{"test_kit_id", "quantity", "shipping_address"},
					"properties": map[string]interface{}{
						"test_kit_id": map[string]interface{}{
							"type":   "string",
							"format": "uuid",
						},
						"quantity": map[string]interface{}{
							"type":    "integer",
							"minimum": 1,
							"maximum": 10,
						},
						"shipping_address": map[string]interface{}{
							"type": "string",
						},
					},
				},
				"Order": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":   "string",
							"format": "uuid",
						},
						"test_kit_id": map[string]interface{}{
							"type":   "string",
							"format": "uuid",
						},
						"quantity": map[string]interface{}{
							"type": "integer",
						},
						"total_price": map[string]interface{}{
							"type": "number",
						},
						"status": map[string]interface{}{
							"type": "string",
						},
						"created_at": map[string]interface{}{
							"type":   "string",
							"format": "date-time",
						},
					},
				},
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling OpenAPI spec: %v", err)
	}

	// Write to file
	err = os.WriteFile("docs/api.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing API documentation: %v", err)
	}

	fmt.Println("API documentation generated at docs/api.json")
}

func main() {
	// Create docs directory if it doesn't exist
	err := os.MkdirAll("docs", 0755)
	if err != nil {
		log.Fatalf("Error creating docs directory: %v", err)
	}

	generateAPIDoc()
}
