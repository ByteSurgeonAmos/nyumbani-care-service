# Nyumbani Care Healthcare Platform - Setup Guide

## Prerequisites

1. **Go 1.21 or later**

   - Download from https://golang.org/dl/
   - Add Go to your PATH environment variable

2. **PostgreSQL Database**
   - Local installation or cloud database (we're using Neon)
   - Database credentials configured in `.env` file
3. **External Service Accounts**
   - Paystack account for payment processing
   - Cloudinary account for file storage
   - OpenAI API key for ChatGPT integration
   - Dreamhost SMTP credentials for email communications

## Installation Steps

### 1. Install Go (if not already installed)

```bash
# Windows (using Chocolatey)
choco install golang

# Or download from https://golang.org/dl/
```

### 2. Clone and Setup Project

```bash
# Navigate to project directory
cd c:\Users\amo$\nyumbanicare

# Download dependencies
go mod download

# Verify dependencies
go mod tidy
```

### 3. Environment Configuration

The `.env` file has been created with the following configurations:

```
# Server settings
PORT=8080
ENV=development
JWT_SECRET=your_jwt_secret_here
JWT_EXPIRATION=24h

# Database connection to Neon PostgreSQL
DB_HOST=your-db-host.neon.tech
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=nyumbanicare
DB_SSLMODE=require

# Cloudinary configuration for file storage
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret
CLOUDINARY_UPLOAD_FOLDER=nyumbanicare

# Dreamhost SMTP for email communications
SMTP_HOST=smtp.dreamhost.com
SMTP_PORT=587
SMTP_USER=notifications@yourdomain.com
SMTP_PASSWORD=your_smtp_password
EMAIL_FROM=Nyumbani Care <notifications@yourdomain.com>

# Paystack payment integration
PAYSTACK_SECRET_KEY=sk_test_your_paystack_secret_key
PAYSTACK_PUBLIC_KEY=pk_test_your_paystack_public_key
PAYSTACK_CALLBACK_URL=https://yourdomain.com/api/payments/paystack/callback

# M-Pesa configuration (backup payment method)
MPESA_CONSUMER_KEY=your_consumer_key
MPESA_CONSUMER_SECRET=your_consumer_secret
MPESA_PASSKEY=your_passkey
MPESA_SHORT_CODE=your_short_code
MPESA_CALLBACK_URL=https://yourdomain.com/api/payments/mpesa/callback

# OpenAI/ChatGPT configuration
OPENAI_API_KEY=your_openai_api_key
OPENAI_MODEL=gpt-4-turbo-preview
```

Make sure to replace the placeholder values with your actual credentials.

### 4. Build and Run

```bash
# Build the application
go build -o nyumbanicare.exe cmd/api/main.go

# Run the application
go run cmd/api/main.go

# Or run the built executable
./nyumbanicare.exe
```

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh JWT token

### Test Kits (Public)

- `GET /api/v1/test-kits` - List available test kits
- `GET /api/v1/test-kits/:id` - Get test kit details

### Protected Endpoints (Require Authentication)

#### User Management

- `GET /api/v1/users/me` - Get current user profile
- `PUT /api/v1/users/me` - Update user profile

#### Orders

- `POST /api/v1/orders` - Create test kit order
- `GET /api/v1/orders` - List user orders
- `GET /api/v1/orders/:id` - Get order details
- `PUT /api/v1/orders/:id/status` - Update order status

#### Medical Records

- `GET /api/v1/medical-records` - List medical records
- `POST /api/v1/medical-records` - Create medical record
- `GET /api/v1/medical-records/:id` - Get medical record
- `PUT /api/v1/medical-records/:id` - Update medical record

#### Prescriptions

- `POST /api/v1/prescriptions` - Upload prescription
- `GET /api/v1/prescriptions` - List prescriptions
- `PUT /api/v1/prescriptions/:id/status` - Update prescription status

#### Lab Work

- `GET /api/v1/lab-tests` - List available lab tests
- `POST /api/v1/lab-bookings` - Book lab test
- `GET /api/v1/lab-bookings` - List lab bookings
- `PUT /api/v1/lab-bookings/:id/status` - Update booking status

#### Telehealth

- `POST /api/v1/telehealth/sessions` - Schedule telehealth session
- `GET /api/v1/telehealth/sessions` - List telehealth sessions

#### Health Education

- `GET /api/v1/health-education/articles` - List health articles
- `GET /api/v1/health-education/articles/:id` - Get article details

#### AI Symptom Checker

- `POST /api/v1/symptoms/check` - Check symptoms
- `GET /api/v1/symptoms/history` - Get symptom check history

#### CareSense Analytics

- `POST /api/v1/caresense/analytics` - Generate analytics
- `GET /api/v1/caresense/analytics` - Get analytics

#### Payments

- `POST /api/v1/payments/paystack` - Process Paystack payment (primary)
- `GET /api/v1/payments/paystack/callback` - Paystack payment callback handler
- `POST /api/webhooks/paystack` - Paystack webhook handler
- `POST /api/v1/payments/mpesa` - Process M-Pesa payment (backup)
- `GET /api/v1/payments` - List user payments
- `GET /api/v1/payments/:id` - Get payment status

#### File Uploads

- `POST /api/v1/uploads/file` - Upload file to Cloudinary
- `POST /api/v1/uploads/file?folder=test_results` - Upload file to specific folder

### Admin Endpoints (Require Admin Role)

- `GET /api/v1/admin/users` - List all users
- `GET /api/v1/admin/orders` - List all orders
- `POST /api/v1/admin/test-kits` - Create test kit
- `PUT /api/v1/admin/test-kits/:id` - Update test kit
- `DELETE /api/v1/admin/test-kits/:id` - Delete test kit
- `POST /api/v1/admin/lab-tests` - Create lab test
- `PUT /api/v1/admin/lab-tests/:id` - Update lab test
- `DELETE /api/v1/admin/lab-tests/:id` - Delete lab test
- `POST /api/v1/admin/health-articles` - Create health article
- `PUT /api/v1/admin/health-articles/:id` - Update health article
- `DELETE /api/v1/admin/health-articles/:id` - Delete health article

### Webhooks

- `POST /api/v1/webhooks/mpesa` - M-Pesa payment callback

## Features Implemented

### Core Features

✅ **Test Kit E-commerce**

- Test kit catalog with search and filtering
- Shopping cart and order management
- Inventory tracking
- Order status updates

✅ **Prescription Management**

- Prescription image upload
- Pharmacy processing workflow
- Medication dispensing tracking
- Status notifications

✅ **Lab Work Booking**

- Lab test catalog
- Online booking system
- Sample collection scheduling
- Results delivery

✅ **Telehealth Consultations**

- Doctor/specialist booking
- Session scheduling
- Virtual consultation support
- Session history

✅ **Health Education**

- Health articles and resources
- Educational content management
- Search and categorization

✅ **AI Symptom Checker**

- Symptom analysis
- Health recommendations
- Risk assessment
- History tracking

✅ **CareSense Analytics**

- Health data analysis
- Trend identification
- Personalized insights
- Risk scoring

### Technical Features

✅ **Authentication & Authorization**

- JWT-based authentication
- Role-based access control (patient, admin)
- Token refresh mechanism

✅ **Payment Integration**

- M-Pesa payment processing
- Stripe integration
- Payment status tracking
- Webhook handling

✅ **File Management**

- File upload functionality
- Image handling for prescriptions
- Document storage

✅ **Email Notifications**

- Appointment confirmations
- Test result notifications
- Prescription updates
- Order confirmations

✅ **Data Validation**

- Input validation middleware
- Business rule enforcement
- Error handling

✅ **Logging & Monitoring**

- Request/response logging
- Security event logging
- Error tracking
- Performance monitoring

✅ **API Documentation**

- OpenAPI 3.0 specification
- Comprehensive endpoint documentation
- Request/response schemas

## Database Schema

The application includes comprehensive data models for:

- User management and authentication
- Medical records and health data
- Test kits and laboratory services
- Orders and payments
- Prescriptions and medications
- Telehealth sessions
- Health education content
- AI-powered features and analytics

## Security Features

- JWT authentication with configurable expiration
- Role-based access control
- Input validation and sanitization
- CORS configuration
- Security event logging
- Rate limiting preparation
- Secure password hashing

## Development Status

### Completed ✅

- Complete API backend implementation
- Database models and migrations
- Authentication and authorization
- All core feature endpoints
- Payment integration framework
- Email notification service
- AI service integration
- File upload handling
- Comprehensive logging
- API documentation generation

### Next Steps 🔄

1. **Frontend Development**: Build React/Vue.js web application
2. **Mobile Apps**: Develop iOS and Android applications
3. **Testing**: Implement comprehensive test suite
4. **Deployment**: Set up CI/CD and production deployment
5. **External Integrations**: Connect real payment gateways and AI services
6. **Performance Optimization**: Implement caching and optimization
7. **Security Hardening**: Add rate limiting and additional security measures

## Getting Started

1. Ensure Go 1.21+ is installed
2. Set up PostgreSQL database or use the provided Neon connection
3. Copy `.env.example` to `.env` and configure settings
4. Run `go mod download` to install dependencies
5. Run `go run cmd/api/main.go` to start the server
6. Access the API at `http://localhost:8080`
7. Use the health check endpoint: `GET /health`

The API is now ready for frontend integration and testing!
