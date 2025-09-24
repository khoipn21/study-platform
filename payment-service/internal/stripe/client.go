package stripe

import (
	"errors"
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/price"
	"github.com/stripe/stripe-go/v82/product"
)

// Client handles Stripe API interactions
type Client struct {
	secretKey string
}

// NewClient creates a new Stripe client
func NewClient(secretKey string) (*Client, error) {
	if secretKey == "" {
		return nil, errors.New("stripe secret key is required")
	}

	// Initialize Stripe with the secret key
	stripe.Key = secretKey

	return &Client{
		secretKey: secretKey,
	}, nil
}

// CreateCustomer creates a new Stripe customer
func (c *Client) CreateCustomer(email, name string, metadata map[string]string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}

	if metadata != nil {
		params.Metadata = metadata
	}

	customer, err := customer.New(params)
	if err != nil {
		log.Printf("Error creating Stripe customer: %v", err)
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

// GetCustomer retrieves a Stripe customer by ID
func (c *Client) GetCustomer(customerID string) (*stripe.Customer, error) {
	customer, err := customer.Get(customerID, nil)
	if err != nil {
		log.Printf("Error getting Stripe customer %s: %v", customerID, err)
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

// CreateProduct creates a new Stripe product
func (c *Client) CreateProduct(name, description string, metadata map[string]string) (*stripe.Product, error) {
	params := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}

	if metadata != nil {
		params.Metadata = metadata
	}

	product, err := product.New(params)
	if err != nil {
		log.Printf("Error creating Stripe product: %v", err)
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// CreatePrice creates a new Stripe price for a product
func (c *Client) CreatePrice(productID string, unitAmount int64, currency string, priceType string, intervalCount int64, interval string, metadata map[string]string) (*stripe.Price, error) {
	params := &stripe.PriceParams{
		UnitAmount: stripe.Int64(unitAmount),
		Currency:   stripe.String(currency),
		Product:    stripe.String(productID),
	}

	// Set up recurring billing if this is a subscription
	if priceType == "recurring" && interval != "" {
		params.Recurring = &stripe.PriceRecurringParams{
			Interval:      stripe.String(interval),
			IntervalCount: stripe.Int64(intervalCount),
		}
	}

	if metadata != nil {
		params.Metadata = metadata
	}

	price, err := price.New(params)
	if err != nil {
		log.Printf("Error creating Stripe price: %v", err)
		return nil, fmt.Errorf("failed to create price: %w", err)
	}

	return price, nil
}

// PaymentIntentOptions contains options for creating a payment intent
type PaymentIntentOptions struct {
	Amount                    int64
	Currency                  string
	CustomerID                string
	Metadata                  map[string]string
	AutomaticPaymentMethods   bool
	SetupFutureUsage          string
	CaptureMethod             string
	ReturnURL                 string
}

// CreatePaymentIntent creates a new payment intent for one-time payments
func (c *Client) CreatePaymentIntent(amount int64, currency, customerID string, metadata map[string]string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		Customer: stripe.String(customerID),
		// Enable automatic payment methods for better frontend integration
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		// Set up for frontend integration
		SetupFutureUsage: stripe.String("off_session"), // Allow saving payment method for future use
		// Capture method - can be "automatic" or "manual"
		CaptureMethod: stripe.String("automatic"),
	}

	if metadata != nil {
		params.Metadata = metadata
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Error creating Stripe payment intent: %v", err)
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return intent, nil
}

// CreatePaymentIntentWithOptions creates a new payment intent with advanced options
func (c *Client) CreatePaymentIntentWithOptions(options PaymentIntentOptions) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(options.Amount),
		Currency: stripe.String(options.Currency),
		Customer: stripe.String(options.CustomerID),
	}

	// Set automatic payment methods if enabled (default: true)
	if options.AutomaticPaymentMethods {
		params.AutomaticPaymentMethods = &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		}
	}

	// Set setup future usage if specified
	if options.SetupFutureUsage != "" {
		params.SetupFutureUsage = stripe.String(options.SetupFutureUsage)
	}

	// Set capture method if specified (default: automatic)
	if options.CaptureMethod != "" {
		params.CaptureMethod = stripe.String(options.CaptureMethod)
	}

	// Set return URL if specified (useful for some payment methods)
	if options.ReturnURL != "" {
		params.ReturnURL = stripe.String(options.ReturnURL)
	}

	if options.Metadata != nil {
		for key, value := range options.Metadata {
			params.AddMetadata(key, value)
		}
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Error creating Stripe payment intent with options: %v", err)
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return intent, nil
}

// GetPaymentIntent retrieves a payment intent by ID
func (c *Client) GetPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	intent, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		log.Printf("Error getting Stripe payment intent %s: %v", paymentIntentID, err)
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}

	return intent, nil
}

// ConfirmPaymentIntent confirms a payment intent
func (c *Client) ConfirmPaymentIntent(paymentIntentID string, paymentMethodID string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String(paymentMethodID),
	}

	intent, err := paymentintent.Confirm(paymentIntentID, params)
	if err != nil {
		log.Printf("Error confirming Stripe payment intent %s: %v", paymentIntentID, err)
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	return intent, nil
}

// GetProduct retrieves a Stripe product by ID
func (c *Client) GetProduct(productID string) (*stripe.Product, error) {
	product, err := product.Get(productID, nil)
	if err != nil {
		log.Printf("Error getting Stripe product %s: %v", productID, err)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// GetPrice retrieves a Stripe price by ID
func (c *Client) GetPrice(priceID string) (*stripe.Price, error) {
	price, err := price.Get(priceID, nil)
	if err != nil {
		log.Printf("Error getting Stripe price %s: %v", priceID, err)
		return nil, fmt.Errorf("failed to get price: %w", err)
	}

	return price, nil
}