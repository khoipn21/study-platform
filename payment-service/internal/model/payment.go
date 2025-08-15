package model

import (
	"time"
)

type PaymentMethod struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	Provider     string    `json:"provider" db:"provider"`
	Token        string    `json:"token" db:"token"`
	CardLastFour string    `json:"card_last_four,omitempty" db:"card_last_four"`
	CardExpiry   string    `json:"card_expiry,omitempty" db:"card_expiry"`
	IsDefault    bool      `json:"is_default" db:"is_default"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Transaction struct {
	ID                   string     `json:"id" db:"id"`
	UserID               string     `json:"user_id" db:"user_id"`
	CourseID             *string    `json:"course_id,omitempty" db:"course_id"`
	PaymentMethodID      *string    `json:"payment_method_id,omitempty" db:"payment_method_id"`
	Amount               float64    `json:"amount" db:"amount"`
	Currency             string     `json:"currency" db:"currency"`
	Status               string     `json:"status" db:"status"`
	TransactionReference *string    `json:"transaction_reference,omitempty" db:"transaction_reference"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

type Subscription struct {
	ID              string     `json:"id" db:"id"`
	UserID          string     `json:"user_id" db:"user_id"`
	PaymentMethodID *string    `json:"payment_method_id,omitempty" db:"payment_method_id"`
	PlanName        string     `json:"plan_name" db:"plan_name"`
	Status          string     `json:"status" db:"status"`
	BillingPeriod   string     `json:"billing_period" db:"billing_period"`
	NextBillingDate *time.Time `json:"next_billing_date,omitempty" db:"next_billing_date"`
	Price           float64    `json:"price" db:"price"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// Request/Response models
type CreatePaymentMethodRequest struct {
	Provider     string `json:"provider" binding:"required"`
	Token        string `json:"token" binding:"required"`
	CardLastFour string `json:"card_last_four,omitempty"`
	CardExpiry   string `json:"card_expiry,omitempty"`
	IsDefault    bool   `json:"is_default"`
}

type UpdatePaymentMethodRequest struct {
	CardLastFour string `json:"card_last_four,omitempty"`
	CardExpiry   string `json:"card_expiry,omitempty"`
	IsDefault    bool   `json:"is_default"`
}

type PurchaseCourseRequest struct {
	PaymentMethodID string  `json:"payment_method_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
	Currency        string  `json:"currency"`
}

type ValidatePaymentRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
	ProviderToken string `json:"provider_token" binding:"required"`
}

type CreateSubscriptionRequest struct {
	PaymentMethodID string  `json:"payment_method_id" binding:"required"`
	PlanName        string  `json:"plan_name" binding:"required"`
	BillingPeriod   string  `json:"billing_period" binding:"required"`
	Price           float64 `json:"price" binding:"required"`
}

type UpdateSubscriptionRequest struct {
	PaymentMethodID string `json:"payment_method_id,omitempty"`
	Status          string `json:"status,omitempty"`
}

type RefundRequest struct {
	Amount string `json:"amount,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// Constants
const (
	// Transaction statuses
	TransactionStatusPending   = "pending"
	TransactionStatusCompleted = "completed"
	TransactionStatusFailed    = "failed"
	TransactionStatusRefunded  = "refunded"
	TransactionStatusCancelled = "cancelled"

	// Subscription statuses
	SubscriptionStatusActive    = "active"
	SubscriptionStatusInactive  = "inactive"
	SubscriptionStatusCancelled = "cancelled"
	SubscriptionStatusExpired   = "expired"

	// Payment providers
	ProviderStripe = "stripe"
	ProviderPayPal = "paypal"

	// Billing periods
	BillingPeriodMonthly = "monthly"
	BillingPeriodYearly  = "yearly"
)