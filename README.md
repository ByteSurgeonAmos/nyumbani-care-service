# Nyumbani Care Backend

A robust, scalable backend system for Nyumbani Care's healthcare platform, built with Go.

## Features

- User Authentication & Authorization
- Test Kit E-commerce Management
- Prescription Management
- Lab Work Booking & Tracking
- Telehealth Consultation System
- Health Education Portal
- AI Symptom Checker Integration
- CareSense: Diagnostics Interpretation & EHR System
- Analytics & Reporting

## Tech Stack

- Go 1.21+
- PostgreSQL (Neon)
- Gin Web Framework
- Paystack Payment Gateway
- Cloudinary for File Storage
- GORM ORM
- JWT Authentication
- Paystack Payment Integration (with M-Pesa as backup)
- Cloudinary File Storage for secure document handling
- ChatGPT API for AI-powered health analysis
- Dreamhost SMTP for reliable email communications
- Docker

## Project Structure

```
.
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── api/              # API handlers and routes
│   ├── config/           # Configuration management
│   ├── middleware/       # HTTP middleware
│   ├── models/           # Data models
│   ├── repository/       # Database operations
│   ├── service/          # Business logic
│   └── utils/            # Utility functions
├── pkg/                   # Public libraries
├── migrations/           # Database migrations
├── docs/                 # Documentation
└── scripts/             # Build and deployment scripts
```

## Service Integrations

### Payment Processing

The application uses Paystack as the payment processor, configured to handle payments in Kenyan Shillings (KES). The payment service supports:

- Payment initialization with secure order tracking
- Webhook handling for asynchronous payment updates
- Payment verification and transaction status tracking
- Comprehensive error handling and logging

For detailed implementation notes, see [PAYSTACK_NOTES.md](docs/PAYSTACK_NOTES.md)

### File Storage

All file uploads are managed through Cloudinary, which provides:

- Secure file storage with access control
- Image and document transformations
- Folder organization for medical records, test results, and prescriptions
- Public/private file access controls

### AI-Powered Health Analysis

The system integrates with OpenAI's ChatGPT API to provide:

- AI symptom analysis and recommendations
- Natural language processing for medical queries
- Health education content generation
- Medical record summarization

### Email Communications

Email functionality uses Dreamhost SMTP servers for reliable delivery of:

- Appointment reminders
- Test results notifications
- Order confirmations
- Password reset links
- Service updates

## Getting Started

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Install dependencies: `go mod download`
4. Run the application: `go run cmd/api/main.go`

## Development

- Run tests: `go test ./...`
- Build: `go build -o nyumbanicare cmd/api/main.go`
- Run with hot reload: `air`

## API Documentation

API documentation is available at `/swagger/index.html` when running the server.

## License

MIT License
