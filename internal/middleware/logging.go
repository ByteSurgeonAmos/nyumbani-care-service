package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	StatusCode int       `json:"status_code"`
	Duration   string    `json:"duration"`
	UserID     string    `json:"user_id,omitempty"`
	Role       string    `json:"role,omitempty"`
	IP         string    `json:"ip"`
	UserAgent  string    `json:"user_agent"`
	Error      string    `json:"error,omitempty"`
}

// Custom response writer to capture response data
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logging middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Capture request body for POST/PUT requests
		var requestBody []byte
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create custom response writer
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:          bytes.NewBufferString(""),
		}
		c.Writer = w

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get user info from context
		userID := ""
		role := ""
		if uid, exists := c.Get("user_id"); exists {
			userID = uid.(string)
		}
		if r, exists := c.Get("role"); exists {
			role = r.(string)
		}

		// Capture error if any
		errorMsg := ""
		if len(c.Errors) > 0 {
			errorMsg = c.Errors.String()
		}

		// Create log entry
		logEntry := LogEntry{
			Timestamp:  start,
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			StatusCode: c.Writer.Status(),
			Duration:   duration.String(),
			UserID:     userID,
			Role:       role,
			IP:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Error:      errorMsg,
		}

		// Log as JSON
		logJSON, _ := json.Marshal(logEntry)
		log.Printf("API_LOG: %s", string(logJSON))

		// Log detailed info for errors
		if c.Writer.Status() >= 400 {
			log.Printf("ERROR_DETAILS: %s %s - Status: %d - Error: %s - User: %s", 
				c.Request.Method, c.Request.URL.Path, c.Writer.Status(), errorMsg, userID)
		}
	}
}

// Security logging for sensitive operations
func SecurityLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log sensitive operations
		sensitiveEndpoints := map[string]bool{
			"/api/v1/auth/login":    true,
			"/api/v1/auth/register": true,
			"/api/v1/payments/":     true,
			"/api/v1/admin/":        true,
		}

		path := c.Request.URL.Path
		for endpoint := range sensitiveEndpoints {
			if len(path) >= len(endpoint) && path[:len(endpoint)] == endpoint {
				userID := ""
				if uid, exists := c.Get("user_id"); exists {
					userID = uid.(string)
				}

				log.Printf("SECURITY_LOG: %s %s - User: %s - IP: %s - UA: %s", 
					c.Request.Method, path, userID, c.ClientIP(), c.Request.UserAgent())
				break
			}
		}

		c.Next()
	}
}

// Rate limiting placeholder
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement actual rate limiting
		// For now, just log high-frequency requests
		
		userID := ""
		if uid, exists := c.Get("user_id"); exists {
			userID = uid.(string)
		}

		// Simple IP-based logging
		ip := c.ClientIP()
		log.Printf("RATE_CHECK: IP=%s User=%s Path=%s", ip, userID, c.Request.URL.Path)
		
		c.Next()
	}
}
