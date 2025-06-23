package services

import (
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/nyumbanicare/internal/config"
)

// ExtendedEmailService handles sending emails with additional functionality
type ExtendedEmailService struct {
	config *config.EmailConfig
}

// NewExtendedEmailService creates a new extended email service instance
func NewExtendedEmailService(config *config.EmailConfig) (*ExtendedEmailService, error) {
	if config == nil {
		return nil, fmt.Errorf("email configuration is required")
	}

	return &ExtendedEmailService{
		config: config,
	}, nil
}

// SendEmail sends an email to the specified recipient
func (s *ExtendedEmailService) SendEmail(to, subject, body string) error {
	// Email server configuration
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPass, s.config.SMTPHost)
	// Prepare email content with headers
	fromHeader := fmt.Sprintf("From: %s <%s>\r\n", s.config.FromName, s.config.FromEmail)
	toHeader := fmt.Sprintf("To: %s\r\n", to)
	subjectHeader := fmt.Sprintf("Subject: %s\r\n", subject)
	mimeHeader := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	message := fromHeader + toHeader + subjectHeader + mimeHeader + formatEmailBody(body)

	// Convert port string to int
	port, err := strconv.Atoi(s.config.SMTPPort)
	if err != nil {
		port = 587 // Default port if conversion fails
	}

	// Send email
	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, port),
		auth,
		s.config.SMTPUser,
		[]string{to},
		[]byte(message),
	)

	return err
}

// Helper function to format email body as HTML
func formatEmailBody(content string) string {
	// Replace newlines with <br> for HTML
	htmlContent := strings.ReplaceAll(content, "\n", "<br>")

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Nyumbani Care</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      line-height: 1.6;
      color: #333;
      max-width: 600px;
      margin: 0 auto;
      padding: 20px;
    }
    .header {
      background-color: #4f46e5;
      padding: 20px;
      text-align: center;
      color: white;
      border-radius: 5px 5px 0 0;
    }
    .content {
      padding: 20px;
      background-color: #f9fafb;
      border: 1px solid #e5e7eb;
      border-top: none;
      border-radius: 0 0 5px 5px;
    }
    .footer {
      margin-top: 20px;
      text-align: center;
      font-size: 12px;
      color: #6b7280;
    }
  </style>
</head>
<body>
  <div class="header">
    <h2>Nyumbani Care</h2>
  </div>
  <div class="content">
    %s
  </div>
  <div class="footer">
    <p>&copy; %d Nyumbani Care. All rights reserved.</p>
    <p>This is an automated message, please do not reply to this email.</p>
  </div>
</body>
</html>
`, htmlContent, getCurrentYear())
}

// Get current year for copyright notice
func getCurrentYear() int {
	return time.Now().Year()
}
