# Nyumbani Care - Testing Guide

## Testing the Paystack Integration

### Prerequisites

- A Paystack test account
- The Nyumbani Care API running on your local machine
- Postman or another API testing tool

## Test Steps

### 1. Initialize a Payment

Send a POST request to `/api/v1/payments/paystack` with the following body:

```json
{
  "email": "test@example.com",
  "amount": 1000,
  "order_id": "<valid-order-id>"
}
```

You should receive a response like:

```json
{
  "payment_id": "uuid",
  "reference": "NC-12345678",
  "status": "pending",
  "message": "Payment initialized successfully"
}
```

### 2. Test the Callback URL

Navigate to the authorization URL provided by Paystack to complete the payment flow.
After completing the payment, you should be redirected to:
`http://localhost:8080/payments/success`

### 3. Test the Webhook

To simulate a Paystack webhook event, you can use tools like ngrok to expose your local server to the internet.

1. Install ngrok: `npm install -g ngrok`
2. Expose your local server: `ngrok http 8080`
3. Configure the webhook URL in your Paystack dashboard to point to: `https://<your-ngrok-url>/api/webhooks/paystack`
4. Create a test payment to trigger the webhook

## Verifying Results

After completing a successful payment:

1. Check the payment status in the database - it should be "completed"
2. Check the order status - it should be "confirmed" with "paid" payment status
3. Check the logs for any errors

## Common Issues

- **Missing Callback URL**: Ensure the callback URL is correctly configured
- **Webhook Not Triggering**: Check if ngrok is running and the webhook URL is properly configured in the Paystack dashboard
- **Payment Status Not Updating**: Check that your database connection is working properly
