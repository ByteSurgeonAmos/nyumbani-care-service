# Nyumbani Care Payment System Documentation

## Overview

The payment system in Nyumbani Care is fully integrated with Paystack, enabling secure and reliable payment processing for test kit orders and other services.

## Currency

All payments are processed in Kenyan Shillings (KES).

## Paystack Integration

### Configuration

The following environment variables need to be set for Paystack integration:

```
PAYSTACK_SECRET_KEY=your_paystack_secret_key
PAYSTACK_PUBLIC_KEY=your_paystack_public_key
```

You must also register your webhook URL in your Paystack dashboard:

```
Webhook URL: https://yourdomain.com/api/webhooks/paystack
```

### API Endpoints

#### 1. Process Payment

- URL: `/api/v1/payments/paystack`
- Method: `POST`
- Authentication: Required
- Request Body:
  ```json
  {
    "email": "customer@example.com",
    "amount": 1000,
    "order_id": "uuid-of-order"
  }
  ```
- Response:
  ```json
  {
    "payment_id": "uuid",
    "reference": "NC-12345678",
    "status": "pending",
    "message": "Payment initialized successfully"
  }
  ```

#### 2. Payment Callback

- URL: `/api/v1/payments/paystack/callback`
- Method: `GET`
- Query Parameters:
  - `reference`: Payment reference from Paystack
- Description: This endpoint is called by Paystack when a user completes payment on Paystack's page. It redirects users to a success or failure page.

#### 3. Webhook Handler

- URL: `/api/webhooks/paystack`
- Method: `POST`
- Description: This endpoint receives webhook notifications from Paystack about payment events.

### Webhook Events

The system handles the following Paystack webhook events:

- `charge.success`: When a payment is successful
- `transfer.success`: When a transfer is successful
- `transfer.reversed`: When a transfer is reversed
- `transfer.failed`: When a transfer fails

### Webhook Security

For security, it's recommended to use IP whitelisting to only accept webhook events from Paystack's official IP addresses. See Paystack's documentation for the current list of IP addresses to whitelist.

## Payment Flow

1. User initiates payment by sending order details to `/api/v1/payments/paystack`
2. The API returns a payment URL where the user is redirected to complete payment
3. After payment, user is redirected to the callback URL
4. In addition, Paystack sends webhook events to our webhook endpoint
5. Upon successful payment, the order status is updated to "confirmed"

## Error Handling

The payment system includes comprehensive error handling for:

- API request failures
- Network timeouts
- Non-success HTTP status codes
- Invalid requests
- Invalid currency
- Missing payment records

## Testing

For testing the integration, use Paystack's test keys and their testing dashboard. To test webhooks, you can use Paystack's webhook simulation tools in their dashboard or services like ngrok to expose your local development environment.
