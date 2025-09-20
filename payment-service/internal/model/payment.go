package model

import (
	"time"

	"github.com/google/uuid"
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
	ID                     string                 `json:"id" db:"id"`
	UserID                 string                 `json:"user_id" db:"user_id"`
	CourseID               *string                `json:"course_id,omitempty" db:"course_id"`
	PaymentMethodID        *string                `json:"payment_method_id,omitempty" db:"payment_method_id"`
	Amount                 float64                `json:"amount" db:"amount"`
	Currency               string                 `json:"currency" db:"currency"`
	Status                 string                 `json:"status" db:"status"`
	TransactionReference   *string                `json:"transaction_reference,omitempty" db:"transaction_reference"`
	LemonSqueezyOrderID    *string                `json:"lemon_squeezy_order_id,omitempty" db:"lemon_squeezy_order_id"`
	LemonSqueezyCheckoutID *string                `json:"lemon_squeezy_checkout_id,omitempty" db:"lemon_squeezy_checkout_id"`
	WebhookEventID         *string                `json:"webhook_event_id,omitempty" db:"webhook_event_id"`
	PaymentVerifiedAt      *time.Time             `json:"payment_verified_at,omitempty" db:"payment_verified_at"`
	PaymentProvider        string                 `json:"payment_provider" db:"payment_provider"`
	RefundedAt             *time.Time             `json:"refunded_at,omitempty" db:"refunded_at"`
	RefundAmount           *float64               `json:"refund_amount,omitempty" db:"refund_amount"`
	ExpiresAt              *time.Time             `json:"expires_at,omitempty" db:"expires_at"`
	CustomData             map[string]interface{} `json:"custom_data,omitempty" db:"custom_data"`
	CreatedAt              time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at" db:"updated_at"`
}

type Subscription struct {
	ID              string     `json:"id" db:"id"`
	UserID          string     `json:"user_id" db:"user_id"`
	PaymentMethodID *string    `json:"payment_method_id,omitempty" db:"payment_method_id"`
	ProviderID      string     `json:"provider_id" db:"provider_id"`
	PlanName        string     `json:"plan_name" db:"plan_name"`
	Status          string     `json:"status" db:"status"`
	BillingPeriod   string     `json:"billing_period" db:"billing_period"`
	NextBillingDate *time.Time `json:"next_billing_date,omitempty" db:"next_billing_date"`
	CancelledAt     *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`
	ExpiredAt       *time.Time `json:"expired_at,omitempty" db:"expired_at"`
	Price           float64    `json:"price" db:"price"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type CheckoutSession struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	CourseID  string    `json:"course_id" db:"course_id"`
	Amount    float64   `json:"amount" db:"amount"`
	Currency  string    `json:"currency" db:"currency"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
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
	VariantID   string                 `json:"variant_id" binding:"required"`
	Amount      float64                `json:"amount" binding:"required"`
	Currency    string                 `json:"currency"`
	CustomData  map[string]interface{} `json:"custom_data,omitempty"`
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

// Lemon Squeezy specific models
type LemonSqueezyWebhookEvent struct {
	ID          string                 `json:"id" db:"id"`
	EventID     string                 `json:"event_id" db:"event_id"`
	EventName   string                 `json:"event_name" db:"event_name"`
	ProcessedAt *time.Time             `json:"processed_at,omitempty" db:"processed_at"`
	Payload     map[string]interface{} `json:"payload" db:"payload"`
	Signature   string                 `json:"signature" db:"signature"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
}

type LemonSqueezyCheckoutRequest struct {
	VariantID   string                 `json:"variant_id" binding:"required"`
	CustomData  map[string]interface{} `json:"custom_data,omitempty"`
	SuccessURL  string                 `json:"success_url,omitempty"`
	CancelURL   string                 `json:"cancel_url,omitempty"`
}

type LemonSqueezyCheckoutResponse struct {
	CheckoutURL string `json:"checkout_url"`
	CheckoutID  string `json:"checkout_id"`
}

type LemonSqueezyOrderData struct {
	ID               string                 `json:"id"`
	Type             string                 `json:"type"`
	Attributes       LemonSqueezyOrderAttrs `json:"attributes"`
	Relationships    map[string]interface{} `json:"relationships"`
}

type LemonSqueezyOrderAttrs struct {
	StoreID         int                    `json:"store_id"`
	CustomerID      int                    `json:"customer_id"`
	Identifier      string                 `json:"identifier"`
	OrderNumber     int                    `json:"order_number"`
	UserName        string                 `json:"user_name"`
	UserEmail       string                 `json:"user_email"`
	Currency        string                 `json:"currency"`
	CurrencyRate    string                 `json:"currency_rate"`
	Subtotal        int                    `json:"subtotal"`
	DiscountTotal   int                    `json:"discount_total"`
	Tax             int                    `json:"tax"`
	Total           int                    `json:"total"`
	SubtotalUSD     int                    `json:"subtotal_usd"`
	DiscountTotalUSD int                   `json:"discount_total_usd"`
	TaxUSD          int                    `json:"tax_usd"`
	TotalUSD        int                    `json:"total_usd"`
	TaxName         *string                `json:"tax_name"`
	TaxRate         string                 `json:"tax_rate"`
	Status          string                 `json:"status"`
	StatusFormatted string                 `json:"status_formatted"`
	Refunded        bool                   `json:"refunded"`
	RefundedAt      *time.Time             `json:"refunded_at"`
	SubtotalFormatted string               `json:"subtotal_formatted"`
	DiscountTotalFormatted string          `json:"discount_total_formatted"`
	TaxFormatted    string                 `json:"tax_formatted"`
	TotalFormatted  string                 `json:"total_formatted"`
	CustomData      map[string]interface{} `json:"custom_data,omitempty"`
	ProductName     string                 `json:"product_name,omitempty"`
	OrderItems      []LemonSqueezyOrderItem `json:"first_order_item"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type LemonSqueezyOrderItem struct {
	ID              int     `json:"id"`
	OrderID         int     `json:"order_id"`
	ProductID       int     `json:"product_id"`
	VariantID       int     `json:"variant_id"`
	ProductName     string  `json:"product_name"`
	VariantName     string  `json:"variant_name"`
	Price           int     `json:"price"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type LemonSqueezyWebhookPayload struct {
	Meta struct {
		EventName  string                 `json:"event_name"`
		CustomData map[string]interface{} `json:"custom_data,omitempty"`
	} `json:"meta"`
	Data LemonSqueezyOrderData `json:"data"`
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
	ProviderLemonSqueezy = "lemonsqueezy"

	// Billing periods
	BillingPeriodMonthly = "monthly"
	BillingPeriodYearly  = "yearly"
)

// PaymentEvent represents a payment event for audit and tracking
type PaymentEvent struct {
	ID              string                 `json:"id" db:"id"`
	EventType       string                 `json:"event_type" db:"event_type"`
	Provider        string                 `json:"provider" db:"provider"`
	ProviderEventID string                 `json:"provider_event_id" db:"provider_event_id"`
	Payload         interface{}            `json:"payload" db:"payload"`
	Processed       bool                   `json:"processed" db:"processed"`
	ProcessedAt     *time.Time             `json:"processed_at,omitempty" db:"processed_at"`
	ErrorMessage    string                 `json:"error_message,omitempty" db:"error_message"`
	RetryCount      int                    `json:"retry_count" db:"retry_count"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
}

// Course represents course information needed for payment processing
type Course struct {
	ID                    string     `json:"id" db:"id"`
	Title                 string     `json:"title" db:"title"`
	Description           string     `json:"description" db:"description"`
	InstructorID          string     `json:"instructor_id" db:"instructor_id"`
	Price                 float64    `json:"price" db:"price"`
	Currency              string     `json:"currency" db:"currency"`
	Status                string     `json:"status" db:"status"`
	IsPaid                bool       `json:"is_paid" db:"is_paid"`
	LemonSqueezyProductID *string    `json:"lemon_squeezy_product_id,omitempty" db:"lemon_squeezy_product_id"`
	LemonSqueezyVariantID *string    `json:"lemon_squeezy_variant_id,omitempty" db:"lemon_squeezy_variant_id"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
}

// Enrollment represents user enrollment in a course
type Enrollment struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	UserID            string     `json:"user_id" db:"user_id"`
	CourseID          string     `json:"course_id" db:"course_id"`
	Status            string     `json:"status" db:"status"`
	PaymentStatus     string     `json:"payment_status" db:"payment_status"`
	PaymentVerifiedAt *time.Time `json:"payment_verified_at,omitempty" db:"payment_verified_at"`
	TransactionID     *string    `json:"transaction_id,omitempty" db:"transaction_id"`
	EnrolledAt        time.Time  `json:"enrolled_at" db:"enrolled_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// AuditLog represents audit logging for important actions
type AuditLog struct {
	ID            string                 `json:"id" db:"id"`
	Action        string                 `json:"action" db:"action"`
	UserID        string                 `json:"user_id" db:"user_id"`
	CourseID      *string                `json:"course_id,omitempty" db:"course_id"`
	TransactionID *string                `json:"transaction_id,omitempty" db:"transaction_id"`
	Details       map[string]interface{} `json:"details,omitempty" db:"details"`
	IPAddress     string                 `json:"ip_address,omitempty" db:"ip_address"`
	Timestamp     time.Time              `json:"timestamp" db:"timestamp"`
}