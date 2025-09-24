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
	// Lemon Squeezy Configuration
	LemonSqueezyAPIKey       string `json:"lemon_squeezy_api_key"`
	LemonSqueezyStoreID      string `json:"lemon_squeezy_store_id"`
	LemonSqueezyProductID    string `json:"lemon_squeezy_product_id"`
	LemonSqueezyVariantID    string `json:"lemon_squeezy_variant_id"`
	LemonSqueezyWebhookSecret string `json:"lemon_squeezy_webhook_secret"`
	LemonSqueezyWebhookURL   string `json:"lemon_squeezy_webhook_url"`
	LemonSqueezyBaseURL      string `json:"lemon_squeezy_base_url"`

	// Stripe Configuration
	StripeSecretKey      string `json:"stripe_secret_key"`
	StripePublishableKey string `json:"stripe_publishable_key"`
	StripeWebhookSecret  string `json:"stripe_webhook_secret"`
	StripeSuccessURL     string `json:"stripe_success_url"`
	StripeCancelURL      string `json:"stripe_cancel_url"`
	StripeWebhookURL     string `json:"stripe_webhook_url"`

	// General Payment Settings
	Currency        string `json:"currency"`
	PaymentProvider string `json:"payment_provider"`
}

type ServicesConfig struct {
	ProgressServiceURL string `json:"progress_service_url"`
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PAYMENT_PORT", "8088"),
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
			LemonSqueezyAPIKey:       getEnv("LEMON_SQUEEZY_API_KEY", ""),
			LemonSqueezyStoreID:      getEnv("LEMON_SQUEEZY_STORE_ID", ""),
			LemonSqueezyProductID:    getEnv("LEMON_SQUEEZY_PRODUCT_ID", ""),
			LemonSqueezyVariantID:    getEnv("LEMON_SQUEEZY_VARIANT_ID", ""),
			LemonSqueezyWebhookSecret: getEnv("LEMON_SQUEEZY_WEBHOOK_SECRET", ""),
			LemonSqueezyWebhookURL:   getEnv("LEMON_SQUEEZY_WEBHOOK_URL", ""),
			LemonSqueezyBaseURL:      getEnv("LEMON_SQUEEZY_BASE_URL", "https://api.lemonsqueezy.com/v1"),
			StripeSecretKey:          getEnv("STRIPE_SECRET_KEY", ""),
			StripePublishableKey:     getEnv("STRIPE_PUBLISHABLE_KEY", ""),
			StripeWebhookSecret:      getEnv("STRIPE_WEBHOOK_SECRET", ""),
			StripeSuccessURL:         getEnv("STRIPE_SUCCESS_URL", "http://localhost:3000/payment/success"),
			StripeCancelURL:          getEnv("STRIPE_CANCEL_URL", "http://localhost:3000/payment/cancel"),
			StripeWebhookURL:         getEnv("STRIPE_WEBHOOK_URL", "http://localhost:8080/api/v1/payments/stripe/webhook"),
			Currency:                 getEnv("PAYMENT_CURRENCY", "VND"),
			PaymentProvider:          getEnv("PAYMENT_PROVIDER", "lemonsqueezy"),
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