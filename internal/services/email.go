package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/nyumbanicare/internal/config"
	"github.com/nyumbanicare/internal/models"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config *config.EmailConfig
}

type EmailData struct {
	To      string
	Subject string
	Content string
	Type    string // html or text
}

func NewEmailService(cfg *config.EmailConfig) *EmailService {
	return &EmailService{config: cfg}
}

func (es *EmailService) SendEmail(data EmailData) error {
	switch es.config.Provider {
	case "smtp":
		return es.sendViaSMTP(data)
	case "sendgrid":
		return es.sendViaSendGrid(data)
	default:
		// For development, just log the email
		fmt.Printf("Email would be sent to %s: %s\n", data.To, data.Subject)
		return nil
	}
}

func (es *EmailService) sendViaSMTP(data EmailData) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", es.config.FromName, es.config.FromEmail))
	m.SetHeader("To", data.To)
	m.SetHeader("Subject", data.Subject)

	if data.Type == "html" {
		m.SetBody("text/html", data.Content)
	} else {
		m.SetBody("text/plain", data.Content)
	}

	port, err := strconv.Atoi(es.config.SMTPPort)
	if err != nil {
		port = 587 // Default to 587 if port is invalid
	}

	d := gomail.NewDialer(es.config.SMTPHost, port, es.config.SMTPUser, es.config.SMTPPass)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email via SMTP: %v", err)
	}

	return nil
}

func (es *EmailService) sendViaSendGrid(data EmailData) error {
	if es.config.APIKey == "" {
		fmt.Printf("SendGrid API key not configured. Email to %s: %s\n", data.To, data.Subject)
		return nil
	}

	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]interface{}{
					{"email": data.To},
				},
			},
		},
		"from": map[string]interface{}{
			"email": es.config.FromEmail,
			"name":  es.config.FromName,
		},
		"subject": data.Subject,
		"content": []map[string]interface{}{
			{
				"type":  data.Type,
				"value": data.Content,
			},
		},
	}

	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+es.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("sendgrid API error: %d", resp.StatusCode)
	}

	return nil
}

// Email templates
func (es *EmailService) SendAppointmentConfirmation(user models.User, session models.TelehealthSession) error {
	content := fmt.Sprintf(`
		<h2>Appointment Confirmed</h2>
		<p>Dear %s,</p>
		<p>Your telehealth appointment has been confirmed for %s.</p>
		<p><strong>Session Details:</strong></p>
		<ul>
			<li>Date: %s</li>
			<li>Type: %s</li>
			<li>Duration: %d minutes</li>
		</ul>
		<p>You will receive a link to join the session 15 minutes before the appointment.</p>
		<p>Best regards,<br>Nyumbani Care Team</p>
	`, user.FirstName, session.ScheduledAt.Format("January 2, 2006 at 3:04 PM"),
		session.ScheduledAt.Format("January 2, 2006 at 3:04 PM"),
		session.SessionType, session.Duration)

	return es.SendEmail(EmailData{
		To:      user.Email,
		Subject: "Telehealth Appointment Confirmed",
		Content: content,
		Type:    "text/html",
	})
}

func (es *EmailService) SendTestResultsReady(user models.User, testResult models.TestKitResult) error {
	content := fmt.Sprintf(`
		<h2>Test Results Ready</h2>
		<p>Dear %s,</p>
		<p>Your test results are now available in your Nyumbani Care account.</p>
		<p><strong>Test Details:</strong></p>
		<ul>
			<li>Result ID: %s</li>
			<li>Date Processed: %s</li>
		</ul>
		<p>Please log in to your account to view your complete results.</p>
		<p>Best regards,<br>Nyumbani Care Team</p>
	`, user.FirstName, testResult.ID.String(), testResult.UpdatedAt.Format("January 2, 2006"))

	return es.SendEmail(EmailData{
		To:      user.Email,
		Subject: "Your Test Results Are Ready",
		Content: content,
		Type:    "text/html",
	})
}

func (es *EmailService) SendPrescriptionUpdate(user models.User, prescription models.Prescription) error {
	statusMessage := ""
	switch prescription.Status {
	case "approved":
		statusMessage = "Your prescription has been approved and is ready for pickup or delivery."
	case "rejected":
		statusMessage = "Your prescription could not be processed. Please contact us for more information."
	case "dispensed":
		statusMessage = "Your prescription has been dispensed and is on its way to you."
	}

	content := fmt.Sprintf(`
		<h2>Prescription Update</h2>
		<p>Dear %s,</p>
		<p>%s</p>
		<p><strong>Prescription ID:</strong> %s</p>
		<p>Please log in to your account for more details.</p>
		<p>Best regards,<br>Nyumbani Care Team</p>
	`, user.FirstName, statusMessage, prescription.ID.String())

	return es.SendEmail(EmailData{
		To:      user.Email,
		Subject: "Prescription Status Update",
		Content: content,
		Type:    "text/html",
	})
}

func (es *EmailService) SendOrderConfirmation(user models.User, order models.TestKitOrder) error {
	content := fmt.Sprintf(`
		<h2>Order Confirmation</h2>
		<p>Dear %s,</p>
		<p>Thank you for your order. We have received your request and will process it shortly.</p>
		<p><strong>Order Details:</strong></p>
		<ul>
			<li>Order ID: %s</li>
			<li>Total: KES %.2f</li>
			<li>Status: %s</li>
		</ul>
		<p>You will receive updates as your order is processed and shipped.</p>
		<p>Best regards,<br>Nyumbani Care Team</p>
	`, user.FirstName, order.ID.String(), order.TotalPrice, order.Status)

	return es.SendEmail(EmailData{
		To:      user.Email,
		Subject: "Order Confirmation - " + order.ID.String(),
		Content: content,
		Type:    "text/html",
	})
}
