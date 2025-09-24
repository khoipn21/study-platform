package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// StripeCustomer represents a Stripe customer record
type StripeCustomer struct {
	ID                 uuid.UUID  `db:"id" json:"id"`
	UserID             *uuid.UUID `db:"user_id" json:"user_id"`
	StripeCustomerID   string     `db:"stripe_customer_id" json:"stripe_customer_id"`
	Email              *string    `db:"email" json:"email"`
	Name               *string    `db:"name" json:"name"`
	Phone              *string    `db:"phone" json:"phone"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
}

// StripeCustomerRepository handles database operations for Stripe customers
type StripeCustomerRepository struct {
	db *sqlx.DB
}

// NewStripeCustomerRepository creates a new repository
func NewStripeCustomerRepository(db *sqlx.DB) *StripeCustomerRepository {
	return &StripeCustomerRepository{db: db}
}

// Create creates a new Stripe customer record
func (r *StripeCustomerRepository) Create(customer *StripeCustomer) error {
	query := `
		INSERT INTO stripe_customers (id, user_id, stripe_customer_id, email, name, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	customer.ID = uuid.New()
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		customer.ID,
		customer.UserID,
		customer.StripeCustomerID,
		customer.Email,
		customer.Name,
		customer.Phone,
		customer.CreatedAt,
		customer.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create stripe customer: %w", err)
	}

	return nil
}

// GetByUserID retrieves a Stripe customer by user ID
func (r *StripeCustomerRepository) GetByUserID(userID uuid.UUID) (*StripeCustomer, error) {
	query := `
		SELECT id, user_id, stripe_customer_id, email, name, phone, created_at, updated_at
		FROM stripe_customers
		WHERE user_id = $1
	`

	var customer StripeCustomer
	err := r.db.Get(&customer, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found, not an error
		}
		return nil, fmt.Errorf("failed to get stripe customer by user ID: %w", err)
	}

	return &customer, nil
}

// GetByStripeCustomerID retrieves a Stripe customer by Stripe customer ID
func (r *StripeCustomerRepository) GetByStripeCustomerID(stripeCustomerID string) (*StripeCustomer, error) {
	query := `
		SELECT id, user_id, stripe_customer_id, email, name, phone, created_at, updated_at
		FROM stripe_customers
		WHERE stripe_customer_id = $1
	`

	var customer StripeCustomer
	err := r.db.Get(&customer, query, stripeCustomerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found, not an error
		}
		return nil, fmt.Errorf("failed to get stripe customer by stripe customer ID: %w", err)
	}

	return &customer, nil
}

// Update updates a Stripe customer record
func (r *StripeCustomerRepository) Update(customer *StripeCustomer) error {
	query := `
		UPDATE stripe_customers
		SET email = $1, name = $2, phone = $3, updated_at = $4
		WHERE id = $5
	`

	customer.UpdatedAt = time.Now()

	result, err := r.db.Exec(query,
		customer.Email,
		customer.Name,
		customer.Phone,
		customer.UpdatedAt,
		customer.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update stripe customer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stripe customer not found")
	}

	return nil
}

// Delete deletes a Stripe customer record
func (r *StripeCustomerRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM stripe_customers WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete stripe customer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stripe customer not found")
	}

	return nil
}

// List retrieves all Stripe customers with pagination
func (r *StripeCustomerRepository) List(limit, offset int) ([]*StripeCustomer, error) {
	query := `
		SELECT id, user_id, stripe_customer_id, email, name, phone, created_at, updated_at
		FROM stripe_customers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var customers []*StripeCustomer
	err := r.db.Select(&customers, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stripe customers: %w", err)
	}

	return customers, nil
}