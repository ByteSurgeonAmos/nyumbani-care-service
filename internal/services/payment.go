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

type PaymentService struct {
	config *config.PaymentConfig
}

func NewPaymentService(cfg *config.PaymentConfig) *PaymentService {
	return &PaymentService{
		config: cfg,
	}
}

type PaystackInitiateRequest struct {
	Email     string                 `json:"email"`
	Amount    float64                `json:"amount"`
	Reference string                 `json:"reference"`
	Callback  string                 `json:"callback_url,omitempty"`
	Currency  string                 `json:"currency,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type PaystackResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaystackPaymentData struct {
	AuthorizationURL string `json:"authorization_url"`
	AccessCode       string `json:"access_code"`
	Reference        string `json:"reference"`
}

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

func (ps *PaymentService) InitiatePayment(order *models.TestKitOrder, email string, callbackURL string) (*models.Payment, error) { // Create a new payment record
	payment := &models.Payment{
		ID:          uuid.New(),
		OrderID:     order.ID,
		UserID:      order.UserID,
		Amount:      order.TotalPrice,
		Currency:    "KES",
		Method:      "paystack",
		Status:      "pending",
		PhoneNumber: "",
	}
	amountInCents := int(payment.Amount * 100)

	reference := fmt.Sprintf("NC-%s", payment.ID.String()[:8])

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
	if ps.config.PaystackSecretKey != "" {
		paymentURL, txnRef, err := ps.initiatePaystackPayment(paystackReq)
		if err != nil {
			return nil, fmt.Errorf("failed to initiate Paystack payment: %v", err)
		}

		payment.TransactionID = txnRef
		payment.Status = "pending"
		payment.FailureReason = paymentURL // Temporary hack to pass the URL back

		return payment, nil
	}

	payment.TransactionID = reference
	payment.Status = "pending"

	return payment, nil
}

func (ps *PaymentService) initiatePaystackPayment(req PaystackInitiateRequest) (string, string, error) {
	url := "https://api.paystack.co/transaction/initialize"

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request data: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+ps.config.PaystackSecretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("API returned non-success status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %v", err)
	}

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

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+ps.config.PaystackSecretKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-success status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var verifyResp PaystackVerifyResponse
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %v", err)
	}

	return &verifyResp, nil
}

func (ps *PaymentService) HandlePaystackWebhook(requestBody []byte) (string, error) {
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

	switch payload.Event {
	case "charge.success":
		if ps.config.PaystackSecretKey != "" {
			verifyResp, err := ps.VerifyPayment(payload.Data.Reference)
			if err != nil {
				return "", fmt.Errorf("failed to verify payment: %v", err)
			}

			if !verifyResp.Status || verifyResp.Data.Status != "success" {
				return "", fmt.Errorf("payment verification failed")
			}

			if verifyResp.Data.Currency != "KES" {
				return "", fmt.Errorf("invalid currency: expected KES, got %s", verifyResp.Data.Currency)
			}
		}

		return payload.Data.Reference, nil

	case "transfer.success", "transfer.reversed", "transfer.failed":
		return payload.Data.Reference, nil

	default:
		// Ignore other events
		return "", nil
	}
}
