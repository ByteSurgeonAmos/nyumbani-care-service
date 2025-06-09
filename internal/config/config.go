package config

import (
	"os"
	"strconv"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	Storage  StorageConfig
	Email    EmailConfig
	Payment  PaymentConfig
	External ExternalConfig
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
	Provider     string
	CloudName    string
	APIKey       string
	APISecret    string
	UploadFolder string
	BucketName   string
	Region       string
	AccessKey    string
	SecretKey    string
}

type EmailConfig struct {
	Provider  string
	SMTPHost  string
	SMTPPort  string
	SMTPUser  string
	SMTPPass  string
	FromName  string
	FromEmail string
	APIKey    string
}

type PaymentConfig struct {
	PaystackSecretKey string
	PaystackPublicKey string
}

type ExternalConfig struct {
	ChatGPTAPIKey   string
	ChatGPTModel    string
	ChatGPTEndpoint string
}

func Load() (*Config, error) {
	return &Config{
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", ""),
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", ""),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", ""),
			SSLMode:  getEnv("DB_SSLMODE", "require"),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET_KEY", ""),
			ExpiresIn: getEnvAsInt("JWT_EXPIRES_IN", 24),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
			Host: getEnv("HOST", "0.0.0.0"),
		}, Storage: StorageConfig{
			Provider:     getEnv("STORAGE_PROVIDER", "cloudinary"),
			CloudName:    getEnv("CLOUDINARY_CLOUD_NAME", "nyumbanicare"),
			APIKey:       getEnv("CLOUDINARY_API_KEY", ""),
			APISecret:    getEnv("CLOUDINARY_API_SECRET", ""),
			UploadFolder: getEnv("CLOUDINARY_UPLOAD_FOLDER", "nyumbanicare-files"),
			BucketName:   getEnv("STORAGE_BUCKET", "nyumbanicare-files"),
			Region:       getEnv("STORAGE_REGION", "us-east-1"),
			AccessKey:    getEnv("STORAGE_ACCESS_KEY", ""),
			SecretKey:    getEnv("STORAGE_SECRET_KEY", ""),
		},
		Email: EmailConfig{
			Provider:  getEnv("EMAIL_PROVIDER", "smtp"),
			SMTPHost:  getEnv("EMAIL_SMTP_HOST", ""),
			SMTPPort:  getEnv("EMAIL_SMTP_PORT", "587"),
			SMTPUser:  getEnv("EMAIL_SMTP_USER", "noreply@nyumbanicare.com"),
			SMTPPass:  getEnv("EMAIL_SMTP_PASSWORD", ""),
			FromName:  getEnv("EMAIL_FROM_NAME", "Nyumbani Care"),
			FromEmail: getEnv("EMAIL_FROM_EMAIL", "noreply@nyumbanicare.com"),
		}, Payment: PaymentConfig{
			PaystackSecretKey: getEnv("PAYSTACK_SECRET_KEY", ""),
			PaystackPublicKey: getEnv("PAYSTACK_PUBLIC_KEY", ""),
		},
		External: ExternalConfig{
			ChatGPTAPIKey:   getEnv("CHATGPT_API_KEY", ""),
			ChatGPTModel:    getEnv("CHATGPT_MODEL", "gpt-4"),
			ChatGPTEndpoint: getEnv("CHATGPT_ENDPOINT", "https://api.openai.com/v1/chat/completions"),
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
