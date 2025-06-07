package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	Email     string  `json:"email"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
	Callback  string  `json:"callback_url,omitempty"`
	Currency  string  `json:"currency,omitempty"`
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
	Status  bool `json:"status"`
	Data    struct {
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
func (ps *PaymentService) InitiatePayment(order *models.TestKitOrder, email string, callbackURL string) (*models.Payment, error) {
	// Create a new payment record
	payment := &models.Payment{
		ID:          uuid.New(),
		OrderID:     order.ID,
		UserID:      order.UserID,
		Amount:      order.TotalPrice,
		Currency:    "NGN", // Nigerian Naira for Paystack
		Method:      "paystack",
		Status:      "pending",
		PhoneNumber: "", // Not used for Paystack
	}
	
	// Convert amount to kobo (Paystack uses the smallest currency unit)
	amountInKobo := int(payment.Amount * 100)
	
	// Generate reference
	reference := fmt.Sprintf("NC-%s", payment.ID.String()[:8])
	
	// Create Paystack request
	paystackReq := PaystackInitiateRequest{
		Email:     email,
		Amount:    float64(amountInKobo),
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
		return "", "", err
	}
	
	// Create HTTP request
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", err
	}
	
	httpReq.Header.Set("Authorization", "Bearer "+ps.config.PaystackSecretKey)
	httpReq.Header.Set("Content-Type", "application/json")
	
	// Make the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	
	// Parse response
	var paystackResp struct {
		Status bool `json:"status"`
		Data   struct {
			AuthorizationURL string `json:"authorization_url"`
			AccessCode       string `json:"access_code"`
			Reference        string `json:"reference"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(body, &paystackResp); err != nil {
		return "", "", err
	}
	
	if !paystackResp.Status {
		return "", "", fmt.Errorf("paystack initialization failed")
	}
	
	return paystackResp.Data.AuthorizationURL, paystackResp.Data.Reference, nil
}

// VerifyPayment verifies a payment transaction with Paystack
func (ps *PaymentService) VerifyPayment(reference string) (*PaystackVerifyResponse, error) {
	if ps.config.PaystackSecretKey == "" {
		// Mock verification for development
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
				Currency:  "NGN",
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
		return nil, err
	}
	
	httpReq.Header.Set("Authorization", "Bearer "+ps.config.PaystackSecretKey)
	
	// Make the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	// Parse response
	var verifyResp PaystackVerifyResponse
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return nil, err
	}
	
	return &verifyResp, nil
}

// HandlePaystackWebhook processes Paystack webhook events
func (ps *PaymentService) HandlePaystackWebhook(requestBody []byte) (string, error) {
	// Parse the webhook payload
	var payload struct {
		Event string `json:"event"`
		Data struct {
			Reference string `json:"reference"`
			Status    string `json:"status"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(requestBody, &payload); err != nil {
		return "", fmt.Errorf("invalid webhook payload: %v", err)
	}
	
	// We only care about successful charges
	if payload.Event != "charge.success" {
		return payload.Data.Reference, nil
	}
	
	// In a real implementation, you would verify the transaction
	// and update the database record accordingly
	return payload.Data.Reference, nil
}
