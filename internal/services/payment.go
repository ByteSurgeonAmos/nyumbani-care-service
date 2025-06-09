package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nyumbanicare/internal/config"
	"github.com/nyumbanicare/internal/models"
)

// PaymentService handles payment processing
type PaymentService struct {
	config *config.PaymentConfig
}

// NewPaymentService creates a new payment service
func NewPaymentService(cfg *config.PaymentConfig) *PaymentService {
	return &PaymentService{
		config: cfg,
	}
}

// PaystackInitiateRequest represents a request to initiate payment via Paystack
type PaystackInitiateRequest struct {
	Email     string                 `json:"email"`
	Amount    float64                `json:"amount"`
	Reference string                 `json:"reference"`
	Callback  string                 `json:"callback_url,omitempty"`
	Currency  string                 `json:"currency,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// PaystackResponse represents a response from Paystack API
type PaystackResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// PaystackPaymentData represents payment data from Paystack
type PaystackPaymentData struct {
	AuthorizationURL string `json:"authorization_url"`
	AccessCode       string `json:"access_code"`
	Reference        string `json:"reference"`
}

// PaystackVerifyResponse represents a verify response from Paystack
type PaystackVerifyResponse struct {
	Status bool `json:"status"`
	Data   struct {
		Status    string `json:"status"`
		Reference string `json:"reference"`
		Amount    int    `json:"amount"`
		Channel   string `json:"channel"`
		Currency  string `json:"currency"`
		Customer  struct {
			Email string `json:"email"`
		} `json:"customer"`
	} `json:"data"`
}

// InitiatePayment initiates a payment transaction
func (ps *PaymentService) InitiatePayment(order *models.TestKitOrder, email string, callbackURL string) (*models.Payment, error) { // Create a new payment record
	payment := &models.Payment{
		ID:          uuid.New(),
		OrderID:     order.ID,
		UserID:      order.UserID,
		Amount:      order.TotalPrice,
		Currency:    "KES", // Kenyan Shillings
		Method:      "paystack",
		Status:      "pending",
		PhoneNumber: "", // Not used for Paystack
	}
	// Convert amount to cents (Paystack uses the smallest currency unit)
	amountInCents := int(payment.Amount * 100)

	// Generate reference
	reference := fmt.Sprintf("NC-%s", payment.ID.String()[:8])

	// Create Paystack request
	paystackReq := PaystackInitiateRequest{
		Email:     email,
		Amount:    float64(amountInCents),
		Reference: reference,
		Callback:  callbackURL,
		Currency:  payment.Currency,
		Metadata: map[string]interface{}{
			"order_id":   order.ID.String(),
			"payment_id": payment.ID.String(),
		},
	}
	// Try to initiate payment with Paystack
	if ps.config.PaystackSecretKey != "" {
		paymentURL, txnRef, err := ps.initiatePaystackPayment(paystackReq)
		if err != nil {
			return nil, fmt.Errorf("failed to initiate Paystack payment: %v", err)
		}

		// Update payment with transaction data
		payment.TransactionID = txnRef
		payment.Status = "pending"
		// Store payment URL in notes field for redirect
		payment.FailureReason = paymentURL // Temporary hack to pass the URL back

		// Return the initialized payment
		return payment, nil
	}

	// Fallback to mock payment for development
	payment.TransactionID = reference
	payment.Status = "pending"

	return payment, nil
}

// initiatePaystackPayment makes API call to Paystack to initialize payment
func (ps *PaymentService) initiatePaystackPayment(req PaystackInitiateRequest) (string, string, error) {
	url := "https://api.paystack.co/transaction/initialize"

	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request data: %v", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+ps.config.PaystackSecretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Make the request with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("API returned non-success status: %s", resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse response
	var paystackResp struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			AuthorizationURL string `json:"authorization_url"`
			AccessCode       string `json:"access_code"`
			Reference        string `json:"reference"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &paystackResp); err != nil {
		return "", "", fmt.Errorf("failed to parse API response: %v", err)
	}

	if !paystackResp.Status {
		return "", "", fmt.Errorf("payment initialization failed: %s", paystackResp.Message)
	}

	return paystackResp.Data.AuthorizationURL, paystackResp.Data.Reference, nil
}

// VerifyPayment verifies a payment transaction with Paystack
func (ps *PaymentService) VerifyPayment(reference string) (*PaystackVerifyResponse, error) {
	if ps.config.PaystackSecretKey == "" { // Mock verification for development
		return &PaystackVerifyResponse{
			Status: true,
			Data: struct {
				Status    string `json:"status"`
				Reference string `json:"reference"`
				Amount    int    `json:"amount"`
				Channel   string `json:"channel"`
				Currency  string `json:"currency"`
				Customer  struct {
					Email string `json:"email"`
				} `json:"customer"`
			}{
				Status:    "success",
				Reference: reference,
				Amount:    1000,
				Channel:   "card",
				Currency:  "KES",
				Customer: struct {
					Email string `json:"email"`
				}{
					Email: "customer@example.com",
				},
			},
		}, nil
	}
	url := fmt.Sprintf("https://api.paystack.co/transaction/verify/%s", reference)

	// Create HTTP request
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+ps.config.PaystackSecretKey)

	// Make the request with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-success status: %s", resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse response
	var verifyResp PaystackVerifyResponse
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %v", err)
	}

	return &verifyResp, nil
}

// Note: Paystack doesn't provide a webhook secret for verification
// Instead, they recommend validating the IP address of the webhook sender
// See https://paystack.com/docs/payments/webhooks/

// HandlePaystackWebhook processes Paystack webhook events
func (ps *PaymentService) HandlePaystackWebhook(requestBody []byte) (string, error) {
	// Parse the webhook payload
	var payload struct {
		Event string `json:"event"`
		Data  struct {
			Reference string `json:"reference"`
			Status    string `json:"status"`
			Amount    int    `json:"amount"`
			Currency  string `json:"currency"`
			Channel   string `json:"channel"`
			Customer  struct {
				Email string `json:"email"`
			} `json:"customer"`
			Metadata struct {
				OrderID   string `json:"order_id"`
				PaymentID string `json:"payment_id"`
			} `json:"metadata"`
		} `json:"data"`
	}

	if err := json.Unmarshal(requestBody, &payload); err != nil {
		return "", fmt.Errorf("invalid webhook payload: %v", err)
	}

	// Handle different webhook events
	switch payload.Event {
	case "charge.success":
		// Verify the transaction with Paystack
		if ps.config.PaystackSecretKey != "" {
			verifyResp, err := ps.VerifyPayment(payload.Data.Reference)
			if err != nil {
				return "", fmt.Errorf("failed to verify payment: %v", err)
			}

			// Ensure the payment was successful
			if !verifyResp.Status || verifyResp.Data.Status != "success" {
				return "", fmt.Errorf("payment verification failed")
			}

			// Ensure currency is KES
			if verifyResp.Data.Currency != "KES" {
				return "", fmt.Errorf("invalid currency: expected KES, got %s", verifyResp.Data.Currency)
			}
		}

		return payload.Data.Reference, nil

	case "transfer.success", "transfer.reversed", "transfer.failed":
		// Handle transfer events
		return payload.Data.Reference, nil

	default:
		// Ignore other events
		return "", nil
	}
}
