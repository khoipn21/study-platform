package config

import (
	"os"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Payment  PaymentConfig  `json:"payment"`
	Services ServicesConfig `json:"services"`
}

type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

type PaymentConfig struct {
	StripeSecretKey      string `json:"stripe_secret_key"`
	StripePublishableKey string `json:"stripe_publishable_key"`
	StripeWebhookSecret  string `json:"stripe_webhook_secret"`
	PayPalClientID       string `json:"paypal_client_id"`
	PayPalClientSecret   string `json:"paypal_client_secret"`
	PayPalSandbox        bool   `json:"paypal_sandbox"`
	Currency             string `json:"currency"`
}

type ServicesConfig struct {
	ProgressServiceURL string `json:"progress_service_url"`
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PAYMENT_PORT", "8086"),
			Host: getEnv("PAYMENT_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "studyplatform"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Payment: PaymentConfig{
			StripeSecretKey:      getEnv("STRIPE_SECRET_KEY", "sk_test_..."),
			StripePublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", "pk_test_..."),
			StripeWebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", "whsec_..."),
			PayPalClientID:       getEnv("PAYPAL_CLIENT_ID", ""),
			PayPalClientSecret:   getEnv("PAYPAL_CLIENT_SECRET", ""),
			PayPalSandbox:        getEnv("PAYPAL_SANDBOX", "true") == "true",
			Currency:             getEnv("PAYMENT_CURRENCY", "USD"),
		},
		Services: ServicesConfig{
			ProgressServiceURL: getEnv("PROGRESS_SERVICE_URL", "progress-service:8080"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}