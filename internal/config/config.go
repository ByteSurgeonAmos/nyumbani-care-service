package config

import (
	"os"
	"strconv"
)

type Config struct {
	Database    DatabaseConfig
	JWT         JWTConfig
	Server      ServerConfig
	Storage     StorageConfig
	Email       EmailConfig
	Payment     PaymentConfig
	External    ExternalConfig
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	SecretKey string
	ExpiresIn int 
}

type ServerConfig struct {
	Port string
	Mode string
	Host string
}

type StorageConfig struct {
	Provider          string // cloudinary, local, azure, aws
	CloudName         string // Cloudinary cloud name
	APIKey            string // Cloudinary API key
	APISecret         string // Cloudinary API secret
	UploadFolder      string // Cloudinary folder for uploads
	BucketName        string // For backward compatibility
	Region            string // For backward compatibility
	AccessKey         string // For backward compatibility
	SecretKey         string // For backward compatibility
}

type EmailConfig struct {
	Provider    string // smtp, sendgrid, etc
	SMTPHost    string // Dreamhost SMTP server
	SMTPPort    string // SMTP port
	SMTPUser    string // SMTP username
	SMTPPass    string // SMTP password
	FromName    string // Sender name
	FromEmail   string // Sender email
	APIKey      string // For backward compatibility
}

type PaymentConfig struct {
	// Paystack configuration
	PaystackSecretKey    string
	PaystackPublicKey    string
	PaystackWebhookSecret string
	
	// M-Pesa configuration (kept for backup)
	MPesaConsumerKey     string
	MPesaConsumerSecret  string
	MPesaShortcode       string
	MPesaPasskey         string
	
	// Removed Stripe configuration
}

type ExternalConfig struct {
	// ChatGPT configuration
	ChatGPTAPIKey    string
	ChatGPTModel     string
	ChatGPTEndpoint  string
	
	// Telehealth configuration
	TelehealthAPIURL string
	TelehealthAPIKey string
	
	// Legacy fields (kept for backward compatibility)
	AIServiceURL     string
	AIServiceAPIKey  string
}

func Load() (*Config, error) {	return &Config{
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "postgres"),
			Host:     getEnv("DB_HOST", "ep-flat-star-a8f5raed-pooler.eastus2.azure.neon.tech"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "datahub_owner"),
			Password: getEnv("DB_PASSWORD", "npg_gs4kl0RAmzBX"),
			DBName:   getEnv("DB_NAME", "datahub"),
			SSLMode:  getEnv("DB_SSLMODE", "require"),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET_KEY", "fghjmkafghjfghjksvsbscscvsvs"),
			ExpiresIn: getEnvAsInt("JWT_EXPIRES_IN", 24),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
			Host: getEnv("HOST", "0.0.0.0"),
		},		Storage: StorageConfig{
			Provider:     getEnv("STORAGE_PROVIDER", "cloudinary"),
			CloudName:    getEnv("CLOUDINARY_CLOUD_NAME", "nyumbanicare"),
			APIKey:       getEnv("CLOUDINARY_API_KEY", ""),
			APISecret:    getEnv("CLOUDINARY_API_SECRET", ""),
			UploadFolder: getEnv("CLOUDINARY_UPLOAD_FOLDER", "nyumbanicare-files"),
			// Keep backwards compatibility
			BucketName:   getEnv("STORAGE_BUCKET", "nyumbanicare-files"),
			Region:       getEnv("STORAGE_REGION", "us-east-1"),
			AccessKey:    getEnv("STORAGE_ACCESS_KEY", ""),
			SecretKey:    getEnv("STORAGE_SECRET_KEY", ""),
		},
		Email: EmailConfig{
			Provider:    getEnv("EMAIL_PROVIDER", "smtp"),
			SMTPHost:    getEnv("EMAIL_SMTP_HOST", "smtp.dreamhost.com"),
			SMTPPort:    getEnv("EMAIL_SMTP_PORT", "587"),
			SMTPUser:    getEnv("EMAIL_SMTP_USER", "noreply@nyumbanicare.com"),
			SMTPPass:    getEnv("EMAIL_SMTP_PASSWORD", ""),
			FromName:    getEnv("EMAIL_FROM_NAME", "Nyumbani Care"),
			FromEmail:   getEnv("EMAIL_FROM_EMAIL", "noreply@nyumbanicare.com"),
			// Keep backwards compatibility
			APIKey:      getEnv("EMAIL_API_KEY", ""),
		},
		Payment: PaymentConfig{
			// Paystack
			PaystackSecretKey:    getEnv("PAYSTACK_SECRET_KEY", ""),
			PaystackPublicKey:    getEnv("PAYSTACK_PUBLIC_KEY", ""),
			PaystackWebhookSecret: getEnv("PAYSTACK_WEBHOOK_SECRET", ""),
			// M-Pesa (backup)
			MPesaConsumerKey:     getEnv("MPESA_CONSUMER_KEY", ""),
			MPesaConsumerSecret:  getEnv("MPESA_CONSUMER_SECRET", ""),
			MPesaShortcode:       getEnv("MPESA_SHORTCODE", ""),
			MPesaPasskey:         getEnv("MPESA_PASSKEY", ""),
		},
		External: ExternalConfig{
			// ChatGPT
			ChatGPTAPIKey:    getEnv("CHATGPT_API_KEY", ""),
			ChatGPTModel:     getEnv("CHATGPT_MODEL", "gpt-4"),
			ChatGPTEndpoint:  getEnv("CHATGPT_ENDPOINT", "https://api.openai.com/v1/chat/completions"),
			// Telehealth
			TelehealthAPIURL: getEnv("TELEHEALTH_API_URL", ""),
			TelehealthAPIKey: getEnv("TELEHEALTH_API_KEY", ""),
			// Legacy AI service
			AIServiceURL:     getEnv("AI_SERVICE_URL", ""),
			AIServiceAPIKey:  getEnv("AI_SERVICE_API_KEY", ""),
		},
	}, nil
}

var globalConfig *Config

func GetConfig() *Config {
	if globalConfig == nil {
		cfg, err := Load()
		if err != nil {
			panic("Failed to load configuration: " + err.Error())
		}
		globalConfig = cfg
	}
	return globalConfig
}




func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}