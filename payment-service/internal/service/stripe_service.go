package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/study-platform/payment-service/internal/config"
	"github.com/study-platform/payment-service/internal/model"
	"github.com/study-platform/payment-service/internal/repository"
	"github.com/study-platform/payment-service/internal/stripe"
)

// StripeService handles Stripe payment operations
type StripeService struct {
	client               *stripe.Client
	config               *config.Config
	stripeCustomerRepo   *repository.StripeCustomerRepository
	stripeProductRepo    *repository.StripeProductRepository
	stripeWebhookRepo    *repository.StripeWebhookRepository
	transactionRepo      *repository.TransactionRepository
	courseRepo           *repository.CourseRepository
	enrollmentRepo       *repository.EnrollmentRepository
}

// NewStripeService creates a new Stripe service
func NewStripeService(
	client *stripe.Client,
	config *config.Config,
	stripeCustomerRepo *repository.StripeCustomerRepository,
	stripeProductRepo *repository.StripeProductRepository,
	stripeWebhookRepo *repository.StripeWebhookRepository,
	transactionRepo *repository.TransactionRepository,
	courseRepo *repository.CourseRepository,
	enrollmentRepo *repository.EnrollmentRepository,
) *StripeService {
	return &StripeService{
		client:               client,
		config:               config,
		stripeCustomerRepo:   stripeCustomerRepo,
		stripeProductRepo:    stripeProductRepo,
		stripeWebhookRepo:    stripeWebhookRepo,
		transactionRepo:      transactionRepo,
		courseRepo:           courseRepo,
		enrollmentRepo:       enrollmentRepo,
	}
}

// CreateOrGetCustomer creates a new Stripe customer or gets an existing one
func (s *StripeService) CreateOrGetCustomer(userID string, email, name string) (*repository.StripeCustomer, error) {
	// Parse user ID
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if customer already exists
	existingCustomer, err := s.stripeCustomerRepo.GetByUserID(parsedUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing customer: %w", err)
	}

	if existingCustomer != nil {
		return existingCustomer, nil
	}

	// Create new Stripe customer
	stripeCustomer, err := s.client.CreateCustomer(email, name, map[string]string{
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	// Save to database
	customer := &repository.StripeCustomer{
		UserID:             &parsedUserID,
		StripeCustomerID:   stripeCustomer.ID,
		Email:              &email,
		Name:               &name,
	}

	err = s.stripeCustomerRepo.Create(customer)
	if err != nil {
		return nil, fmt.Errorf("failed to save customer to database: %w", err)
	}

	return customer, nil
}

// CreateOrGetProductForCourse creates a Stripe product and price for a course
func (s *StripeService) CreateOrGetProductForCourse(courseID string) (*repository.StripeProduct, error) {
	// Parse course ID
	parsedCourseID, err := uuid.Parse(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course ID: %w", err)
	}

	// Check if product already exists
	existingProduct, err := s.stripeProductRepo.GetByCourseID(parsedCourseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing product: %w", err)
	}

	if existingProduct != nil {
		return existingProduct, nil
	}

	// Get course details
	course, err := s.courseRepo.GetByID(context.Background(), courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	// Create Stripe product
	productName := course.Title
	productDescription := course.Description

	stripeProduct, err := s.client.CreateProduct(
		productName,
		productDescription,
		map[string]string{
			"course_id": courseID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe product: %w", err)
	}

	// Convert price to cents
	priceInCents := int64(course.Price * 100)
	currency := strings.ToLower(course.Currency)

	// Create Stripe price
	stripePrice, err := s.client.CreatePrice(
		stripeProduct.ID,
		priceInCents,
		currency,
		"one_time", // For now, only supporting one-time payments
		1,
		"",
		map[string]string{
			"course_id": courseID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe price: %w", err)
	}

	// Save to database
	active := true
	product := &repository.StripeProduct{
		CourseID:           &parsedCourseID,
		StripeProductID:    stripeProduct.ID,
		StripePriceID:      stripePrice.ID,
		ProductName:        productName,
		ProductDescription: &productDescription,
		PriceAmount:        priceInCents,
		PriceCurrency:      currency,
		PriceType:          "one_time",
		Active:             &active,
	}

	err = s.stripeProductRepo.Create(product)
	if err != nil {
		return nil, fmt.Errorf("failed to save product to database: %w", err)
	}

	return product, nil
}

// CreatePaymentIntent creates a payment intent for course purchase
func (s *StripeService) CreatePaymentIntent(userID, courseID, email, userName string) (*PaymentIntentResponse, error) {
	// Create or get customer
	customer, err := s.CreateOrGetCustomer(userID, email, userName)
	if err != nil {
		return nil, fmt.Errorf("failed to create/get customer: %w", err)
	}

	// Create or get product
	product, err := s.CreateOrGetProductForCourse(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to create/get product: %w", err)
	}

	// Create payment intent
	paymentIntent, err := s.client.CreatePaymentIntent(
		product.PriceAmount,
		product.PriceCurrency,
		customer.StripeCustomerID,
		map[string]string{
			"user_id":   userID,
			"course_id": courseID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Create pending transaction record
	transaction := &model.Transaction{
		ID:                    uuid.New().String(),
		UserID:               userID,
		CourseID:             &courseID,
		Amount:               float64(product.PriceAmount) / 100, // Convert cents to dollars
		Currency:             product.PriceCurrency,
		Status:               model.TransactionStatusPending,
		StripePaymentIntentID: &paymentIntent.ID,
		StripeCustomerID:     &customer.StripeCustomerID,
		TransactionReference: &paymentIntent.ID,
		PaymentProvider:      "stripe",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	err = s.transactionRepo.Create(context.Background(), transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	return &PaymentIntentResponse{
		PaymentIntentID:   paymentIntent.ID,
		ClientSecret:      paymentIntent.ClientSecret,
		Amount:            paymentIntent.Amount,
		Currency:          string(paymentIntent.Currency),
		Status:            string(paymentIntent.Status),
		CustomerID:        customer.StripeCustomerID,
		TransactionID:     transaction.ID,
	}, nil
}

// CreatePaymentIntentWithOverrides creates a payment intent with optional amount and currency overrides
func (s *StripeService) CreatePaymentIntentWithOverrides(userID, courseID, email, userName string, amountOverride *int64, currencyOverride *string) (*PaymentIntentResponse, error) {
	// Create or get customer
	customer, err := s.CreateOrGetCustomer(userID, email, userName)
	if err != nil {
		return nil, fmt.Errorf("failed to create/get customer: %w", err)
	}

	// Create or get product
	product, err := s.CreateOrGetProductForCourse(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to create/get product: %w", err)
	}

	// Use overrides if provided, otherwise use product values
	amount := product.PriceAmount
	currency := product.PriceCurrency

	if amountOverride != nil {
		amount = *amountOverride
	}

	if currencyOverride != nil {
		currency = strings.ToLower(*currencyOverride)
	}

	// Create payment intent with enhanced options
	paymentIntent, err := s.client.CreatePaymentIntentWithOptions(stripe.PaymentIntentOptions{
		Amount:                  amount,
		Currency:                currency,
		CustomerID:              customer.StripeCustomerID,
		AutomaticPaymentMethods: true,
		SetupFutureUsage:        "off_session",
		CaptureMethod:           "automatic",
		Metadata: map[string]string{
			"user_id":   userID,
			"course_id": courseID,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Create pending transaction record
	transaction := &model.Transaction{
		ID:                    uuid.New().String(),
		UserID:               userID,
		CourseID:             &courseID,
		Amount:               float64(amount) / 100, // Convert cents to dollars
		Currency:             currency,
		Status:               model.TransactionStatusPending,
		StripePaymentIntentID: &paymentIntent.ID,
		StripeCustomerID:     &customer.StripeCustomerID,
		TransactionReference: &paymentIntent.ID,
		PaymentProvider:      "stripe",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	err = s.transactionRepo.Create(context.Background(), transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	return &PaymentIntentResponse{
		PaymentIntentID:   paymentIntent.ID,
		ClientSecret:      paymentIntent.ClientSecret,
		Amount:            paymentIntent.Amount,
		Currency:          string(paymentIntent.Currency),
		Status:            string(paymentIntent.Status),
		CustomerID:        customer.StripeCustomerID,
		TransactionID:     transaction.ID,
	}, nil
}

// ConfirmPaymentIntent confirms a payment intent
func (s *StripeService) ConfirmPaymentIntent(paymentIntentID, paymentMethodID string) error {
	_, err := s.client.ConfirmPaymentIntent(paymentIntentID, paymentMethodID)
	if err != nil {
		return fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	return nil
}

// ProcessSuccessfulPayment handles successful payment processing
func (s *StripeService) ProcessSuccessfulPayment(paymentIntentID string) error {
	// Get payment intent from Stripe
	paymentIntent, err := s.client.GetPaymentIntent(paymentIntentID)
	if err != nil {
		return fmt.Errorf("failed to get payment intent: %w", err)
	}

	// Get transaction from database
	transaction, err := s.transactionRepo.GetByStripePaymentIntentID(paymentIntentID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if transaction == nil {
		return fmt.Errorf("transaction not found for payment intent %s", paymentIntentID)
	}

	// Update transaction status
	transaction.Status = model.TransactionStatusCompleted
	if paymentIntent.LatestCharge != nil {
		chargeID := paymentIntent.LatestCharge.ID
		transaction.StripeChargeID = &chargeID
	}
	now := time.Now()
	transaction.PaymentVerifiedAt = &now
	transaction.UpdatedAt = now

	err = s.transactionRepo.Update(transaction)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	// Create enrollment if it doesn't exist
	if transaction.CourseID != nil {
		enrollment := &model.Enrollment{
			ID:              uuid.New(),
			UserID:          transaction.UserID,
			CourseID:        *transaction.CourseID,
			Status:          "active",
			PaymentStatus:   "paid",
			TransactionID:   &transaction.ID,
			EnrolledAt:      now,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		err = s.enrollmentRepo.Create(context.Background(), enrollment)
		if err != nil {
			log.Printf("Warning: failed to create enrollment for user %s, course %s: %v",
				transaction.UserID, *transaction.CourseID, err)
			// Don't return error here as payment is already processed
		}
	}

	return nil
}

// ProcessFailedPayment handles failed payment processing
func (s *StripeService) ProcessFailedPayment(paymentIntentID string, errorMessage string) error {
	// Get transaction from database
	transaction, err := s.transactionRepo.GetByStripePaymentIntentID(paymentIntentID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if transaction == nil {
		return fmt.Errorf("transaction not found for payment intent %s", paymentIntentID)
	}

	// Update transaction status
	transaction.Status = model.TransactionStatusFailed
	// Store error message in custom_data
	errorData := map[string]interface{}{
		"error_message": errorMessage,
		"failed_at":     "stripe_processing",
		"failed_at_timestamp": time.Now().Format(time.RFC3339),
	}

	transaction.CustomData = errorData
	transaction.UpdatedAt = time.Now()

	err = s.transactionRepo.Update(transaction)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}

// PaymentIntentResponse represents the response for creating a payment intent
type PaymentIntentResponse struct {
	PaymentIntentID   string `json:"payment_intent_id"`
	ClientSecret      string `json:"client_secret"`
	Amount            int64  `json:"amount"`
	Currency          string `json:"currency"`
	Status            string `json:"status"`
	CustomerID        string `json:"customer_id"`
	TransactionID     string `json:"transaction_id"`
}

// GetPaymentIntent retrieves a payment intent by ID
func (s *StripeService) GetPaymentIntent(paymentIntentID string) (*PaymentIntentResponse, error) {
	paymentIntent, err := s.client.GetPaymentIntent(paymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}

	// Get transaction from database
	transaction, err := s.transactionRepo.GetByStripePaymentIntentID(paymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	var transactionID string
	var customerID string

	if transaction != nil {
		transactionID = transaction.ID
		if transaction.StripeCustomerID != nil {
			customerID = *transaction.StripeCustomerID
		}
	}

	return &PaymentIntentResponse{
		PaymentIntentID:   paymentIntent.ID,
		ClientSecret:      paymentIntent.ClientSecret,
		Amount:            paymentIntent.Amount,
		Currency:          string(paymentIntent.Currency),
		Status:            string(paymentIntent.Status),
		CustomerID:        customerID,
		TransactionID:     transactionID,
	}, nil
}

// CoursePurchaseResponse represents the response for course purchase
type CoursePurchaseResponse struct {
	TransactionID     string    `json:"transaction_id"`
	CourseID          string    `json:"course_id"`
	UserID            string    `json:"user_id"`
	PaymentIntentID   string    `json:"payment_intent_id"`
	Amount            float64   `json:"amount"`
	Currency          string    `json:"currency"`
	Status            string    `json:"status"`
	EnrollmentID      string    `json:"enrollment_id"`
	EnrolledAt        time.Time `json:"enrolled_at"`
}

// ProcessCoursePurchase processes a course purchase using a payment intent
func (s *StripeService) ProcessCoursePurchase(userID, courseID, paymentIntentID string) (*CoursePurchaseResponse, error) {
	// Validate UUIDs
	if _, err := uuid.Parse(userID); err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	if _, err := uuid.Parse(courseID); err != nil {
		return nil, fmt.Errorf("invalid course ID format: %w", err)
	}

	// Get payment intent from Stripe to verify it's successful
	paymentIntent, err := s.client.GetPaymentIntent(paymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("payment intent not found")
	}

	// Check if payment is completed
	if paymentIntent.Status != "succeeded" {
		return nil, fmt.Errorf("payment not completed")
	}

	// Verify the payment intent belongs to the correct user
	userIDFromMetadata := paymentIntent.Metadata["user_id"]
	if userIDFromMetadata != userID {
		return nil, fmt.Errorf("payment intent does not belong to the specified user")
	}

	courseIDFromMetadata := paymentIntent.Metadata["course_id"]
	if courseIDFromMetadata != courseID {
		return nil, fmt.Errorf("payment intent does not belong to the specified course")
	}

	// Check if course exists
	_, err = s.courseRepo.GetByID(context.Background(), courseID)
	if err != nil {
		return nil, fmt.Errorf("course not found")
	}

	// Check if user is already enrolled
	existingEnrollment, err := s.enrollmentRepo.GetByUserAndCourse(context.Background(), userID, courseID)
	if err == nil && existingEnrollment != nil {
		return nil, fmt.Errorf("user already enrolled")
	}

	// Get transaction from database
	transaction, err := s.transactionRepo.GetByStripePaymentIntentID(paymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	if transaction == nil {
		// Create transaction if it doesn't exist (shouldn't happen, but safety measure)
		transaction = &model.Transaction{
			ID:                    uuid.New().String(),
			UserID:               userID,
			CourseID:             &courseID,
			Amount:               float64(paymentIntent.Amount) / 100, // Convert cents to dollars
			Currency:             string(paymentIntent.Currency),
			Status:               model.TransactionStatusCompleted,
			StripePaymentIntentID: &paymentIntent.ID,
			TransactionReference:  &paymentIntent.ID,
			PaymentProvider:      "stripe",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

		if paymentIntent.Customer != nil {
			customerID := paymentIntent.Customer.ID
			transaction.StripeCustomerID = &customerID
		}

		if paymentIntent.LatestCharge != nil {
			chargeID := paymentIntent.LatestCharge.ID
			transaction.StripeChargeID = &chargeID
		}

		now := time.Now()
		transaction.PaymentVerifiedAt = &now

		err = s.transactionRepo.Create(context.Background(), transaction)
		if err != nil {
			return nil, fmt.Errorf("failed to create transaction record: %w", err)
		}
	} else {
		// Update existing transaction to completed if not already
		if transaction.Status != model.TransactionStatusCompleted {
			transaction.Status = model.TransactionStatusCompleted

			if paymentIntent.LatestCharge != nil {
				chargeID := paymentIntent.LatestCharge.ID
				transaction.StripeChargeID = &chargeID
			}

			now := time.Now()
			transaction.PaymentVerifiedAt = &now
			transaction.UpdatedAt = now

			err = s.transactionRepo.Update(transaction)
			if err != nil {
				return nil, fmt.Errorf("failed to update transaction: %w", err)
			}
		}
	}

	// Create enrollment
	enrollmentID := uuid.New()
	now := time.Now()
	enrollment := &model.Enrollment{
		ID:              enrollmentID,
		UserID:          userID,
		CourseID:        courseID,
		Status:          "active",
		PaymentStatus:   "paid",
		TransactionID:   &transaction.ID,
		EnrolledAt:      now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = s.enrollmentRepo.Create(context.Background(), enrollment)
	if err != nil {
		return nil, fmt.Errorf("failed to create enrollment: %w", err)
	}

	// Return success response
	return &CoursePurchaseResponse{
		TransactionID:   transaction.ID,
		CourseID:        courseID,
		UserID:          userID,
		PaymentIntentID: paymentIntentID,
		Amount:          transaction.Amount,
		Currency:        transaction.Currency,
		Status:          "completed",
		EnrollmentID:    enrollmentID.String(),
		EnrolledAt:      now,
	}, nil
}