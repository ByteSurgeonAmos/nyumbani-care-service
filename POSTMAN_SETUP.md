# Nyumbani Care API - Postman Collection

This directory contains a comprehensive Postman collection for testing all endpoints of the Nyumbani Care Healthcare Platform API.

## Files Included

1. **`Nyumbani_Care_API.postman_collection.json`** - Complete API collection with all endpoints
2. **`Nyumbani_Care_Local_Environment.postman_environment.json`** - Environment variables for local development
3. **`POSTMAN_SETUP.md`** - This setup guide

## Quick Setup

### 1. Import into Postman

1. Open Postman
2. Click **Import** button
3. Select **Files** tab
4. Choose both files:
   - `Nyumbani_Care_API.postman_collection.json`
   - `Nyumbani_Care_Local_Environment.postman_environment.json`
5. Click **Import**

### 2. Configure Environment

1. Select **"Nyumbani Care - Local Development"** environment from the dropdown
2. Update the `baseUrl` variable if your API runs on a different port
3. Other variables will be automatically populated during testing

### 3. Start Testing

1. Ensure your Nyumbani Care API is running (`go run cmd/api/main.go`)
2. Start with the **Health Check** request to verify connectivity
3. Use **Register User** or **Login User** to get authentication tokens
4. Explore the protected endpoints using the populated JWT tokens

## Collection Structure

### 📋 **Health Check**

- Health status endpoint
- API documentation endpoint

### 🔐 **Authentication**

- User registration with automatic token extraction
- User login with automatic token and user ID extraction
- Token refresh functionality

### 👤 **User Management**

- Get current user profile
- Update user profile information

### 🧪 **Test Kits**

- List available test kits (public)
- Get test kit details (public)

### 📦 **Orders**

- Create test kit orders
- List user orders
- Get order details
- Update order status (admin)

### 🔬 **Test Kit Results**

- Upload and analyze test kit images with AI
- List test results with pagination
- Get individual test results
- Update results (healthcare professionals)

### 📋 **Medical Records**

- CRUD operations for medical records
- Support for medications, test results, and consultations

### 🔔 **Notifications**

- Get user notifications
- Create notifications
- Mark as read (individual and bulk)
- Get unread count
- Email notifications

### 💊 **Prescriptions**

- Upload prescription images
- List prescriptions
- Update prescription status

### 🏥 **Lab Tests & Bookings**

- List available lab tests (public)
- Create lab bookings
- List user bookings
- Update booking status

### 📞 **Telehealth**

- Schedule telehealth sessions
- List sessions for users/providers

### 🤖 **AI Symptom Checker**

- Submit symptoms for AI analysis
- View symptom check history

### 📊 **CareSense Analytics**

- Generate health trends analytics
- Generate risk assessment analytics
- Generate wellness score analytics
- Generate comprehensive analytics
- Retrieve stored analytics

### 📚 **Health Education**

- List health articles (public)
- Get article details (public)

### 💳 **Payments**

- Process Paystack payments
- Handle payment callbacks
- List user payments
- Get payment status

### 📁 **File Uploads**

- Upload files to Cloudinary with folder organization

### 🔗 **Webhooks**

- Paystack webhook handler for payment notifications

### ⚙️ **Admin Endpoints**

- User management
- Order management
- Test kit CRUD operations
- Lab test CRUD operations
- Health article CRUD operations

## Environment Variables

The collection uses these environment variables:

| Variable          | Description         | Auto-populated        |
| ----------------- | ------------------- | --------------------- |
| `baseUrl`         | API base URL        | ❌                    |
| `jwt_token`       | User JWT token      | ✅ (from login)       |
| `admin_token`     | Admin JWT token     | ✅ (from admin login) |
| `user_id`         | Current user ID     | ✅ (from login)       |
| `test_kit_id`     | Test kit UUID       | ❌                    |
| `order_id`        | Order UUID          | ❌                    |
| `result_id`       | Test result UUID    | ❌                    |
| `record_id`       | Medical record UUID | ❌                    |
| `notification_id` | Notification UUID   | ❌                    |
| `prescription_id` | Prescription UUID   | ❌                    |
| `booking_id`      | Lab booking UUID    | ❌                    |
| `lab_test_id`     | Lab test UUID       | ❌                    |
| `article_id`      | Health article UUID | ❌                    |
| `payment_id`      | Payment UUID        | ❌                    |

## Authentication Flow

1. **Register** or **Login** to get JWT token
2. Token is automatically stored in `jwt_token` variable
3. Protected endpoints automatically use this token
4. Admin endpoints require admin role (use `admin_token`)

## Testing Workflow

### Basic User Flow

1. Health Check → Register User → Login User
2. Get Current User → Update User Profile
3. List Test Kits → Create Order
4. Upload Test Result → Get Results

### Healthcare Professional Flow

1. Login (with healthcare professional role)
2. List Test Kit Results → Update Test Result
3. List Prescriptions → Update Prescription Status

### Admin Flow

1. Login (with admin role)
2. Create Test Kit → Update Test Kit
3. List All Users → List All Orders
4. Create Health Article

### Analytics Flow

1. Login User
2. Create some test data (symptoms, test results)
3. Generate Analytics (various types)
4. Get Analytics Results

## Tips for Testing

### 🔄 **Automatic Token Management**

- JWT tokens are automatically extracted and stored after login
- No need to manually copy-paste tokens
- Tokens are used automatically in Authorization headers

### 📝 **Sample Data**

- All requests include realistic sample data
- Modify JSON bodies as needed for your testing
- UUID placeholders are used for relationship fields

### 📋 **Pagination**

- List endpoints support pagination parameters
- Use `page` and `limit` query parameters

### 📁 **File Uploads**

- Test kit results and prescriptions support file uploads
- Select actual image files for realistic testing

### 🔍 **Error Testing**

- Try requests without authentication to test middleware
- Use invalid UUIDs to test error handling
- Submit malformed JSON to test validation

## Common Issues

### ❌ **Connection Errors**

- Verify API is running on correct port
- Check `baseUrl` in environment settings
- Ensure no firewall blocking local connections

### 🔐 **Authentication Errors**

- Ensure you've logged in to get JWT token
- Check token hasn't expired (refresh if needed)
- Verify user has correct role for admin endpoints

### 📊 **Empty Analytics**

- Analytics require existing user data
- Create test results, symptoms, etc. first
- Allow some time for data to accumulate

## API Documentation

For detailed API documentation, visit:

- Swagger UI: `http://localhost:8080/swagger/index.html`
- API docs: `http://localhost:8080/api-docs`

## Support

If you encounter issues:

1. Check the API logs for error details
2. Verify database connection and migrations
3. Ensure all required environment variables are set
4. Review the API documentation for endpoint requirements

---

**Happy Testing! 🚀**
