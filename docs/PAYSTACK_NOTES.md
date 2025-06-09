# Paystack Integration - Implementation Notes

## Overview

Nyumbani Care uses Paystack for payment processing, specifically configured to accept payments in Kenyan Shillings (KES). This document provides implementation details and notes for developers.

## Key Components

### 1. Environment Configuration

```
PAYSTACK_SECRET_KEY=your_secret_key
PAYSTACK_PUBLIC_KEY=your_public_key
```

### 2. Webhook Setup

Unlike some other payment processors, Paystack doesn't use a webhook secret for verification. Instead:

1. Set up your webhook URL in the Paystack dashboard:
   `https://yourdomain.com/api/webhooks/paystack`

2. For added security, you can implement IP address filtering to only accept webhooks from Paystack's servers.
   Paystack's webhook IPs can be found in their documentation.

### 3. Payment Flow

1. User initiates payment → Creates pending payment record
2. User completes payment on Paystack → Webhook received
3. We verify the payment with Paystack's API
4. We update the payment and order status

### 4. API Endpoints

- `POST /api/v1/payments/paystack` - Initiate payment
- `GET /api/v1/payments/paystack/callback` - Handle redirect after payment
- `POST /api/webhooks/paystack` - Receive webhooks from Paystack

### 5. Verification Mechanism

Since Paystack doesn't use a webhook secret, we verify payments by:

1. Directly querying Paystack's API to verify the transaction status
2. Cross-checking the transaction details against our records

## Additional Security Considerations

- Consider implementing IP whitelisting for Paystack's webhook endpoints
- Validate all payment amounts against your database records before marking as complete
- Add rate limiting to prevent abuse of your payment endpoints

## Testing

Use Paystack's test keys during development to ensure your integration works smoothly.

See TESTING.md for detailed testing procedures.
