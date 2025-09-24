package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// StripeProduct represents a Stripe product record
type StripeProduct struct {
	ID                      uuid.UUID  `db:"id" json:"id"`
	CourseID                *uuid.UUID `db:"course_id" json:"course_id"`
	StripeProductID         string     `db:"stripe_product_id" json:"stripe_product_id"`
	StripePriceID           string     `db:"stripe_price_id" json:"stripe_price_id"`
	ProductName             string     `db:"product_name" json:"product_name"`
	ProductDescription      *string    `db:"product_description" json:"product_description"`
	PriceAmount             int64      `db:"price_amount" json:"price_amount"`
	PriceCurrency           string     `db:"price_currency" json:"price_currency"`
	PriceType               string     `db:"price_type" json:"price_type"`
	RecurringInterval       *string    `db:"recurring_interval" json:"recurring_interval"`
	RecurringIntervalCount  *int       `db:"recurring_interval_count" json:"recurring_interval_count"`
	Active                  *bool      `db:"active" json:"active"`
	CreatedAt               time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt               time.Time  `db:"updated_at" json:"updated_at"`
}

// StripeProductRepository handles database operations for Stripe products
type StripeProductRepository struct {
	db *sqlx.DB
}

// NewStripeProductRepository creates a new repository
func NewStripeProductRepository(db *sqlx.DB) *StripeProductRepository {
	return &StripeProductRepository{db: db}
}

// Create creates a new Stripe product record
func (r *StripeProductRepository) Create(product *StripeProduct) error {
	query := `
		INSERT INTO stripe_products (
			id, course_id, stripe_product_id, stripe_price_id, product_name,
			product_description, price_amount, price_currency, price_type,
			recurring_interval, recurring_interval_count, active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	product.ID = uuid.New()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		product.ID,
		product.CourseID,
		product.StripeProductID,
		product.StripePriceID,
		product.ProductName,
		product.ProductDescription,
		product.PriceAmount,
		product.PriceCurrency,
		product.PriceType,
		product.RecurringInterval,
		product.RecurringIntervalCount,
		product.Active,
		product.CreatedAt,
		product.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create stripe product: %w", err)
	}

	return nil
}

// GetByCourseID retrieves a Stripe product by course ID
func (r *StripeProductRepository) GetByCourseID(courseID uuid.UUID) (*StripeProduct, error) {
	query := `
		SELECT id, course_id, stripe_product_id, stripe_price_id, product_name,
			   product_description, price_amount, price_currency, price_type,
			   recurring_interval, recurring_interval_count, active, created_at, updated_at
		FROM stripe_products
		WHERE course_id = $1 AND active = true
	`

	var product StripeProduct
	err := r.db.Get(&product, query, courseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found, not an error
		}
		return nil, fmt.Errorf("failed to get stripe product by course ID: %w", err)
	}

	return &product, nil
}

// GetByStripeProductID retrieves a Stripe product by Stripe product ID
func (r *StripeProductRepository) GetByStripeProductID(stripeProductID string) (*StripeProduct, error) {
	query := `
		SELECT id, course_id, stripe_product_id, stripe_price_id, product_name,
			   product_description, price_amount, price_currency, price_type,
			   recurring_interval, recurring_interval_count, active, created_at, updated_at
		FROM stripe_products
		WHERE stripe_product_id = $1
	`

	var product StripeProduct
	err := r.db.Get(&product, query, stripeProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found, not an error
		}
		return nil, fmt.Errorf("failed to get stripe product by stripe product ID: %w", err)
	}

	return &product, nil
}

// GetByStripePriceID retrieves a Stripe product by Stripe price ID
func (r *StripeProductRepository) GetByStripePriceID(stripePriceID string) (*StripeProduct, error) {
	query := `
		SELECT id, course_id, stripe_product_id, stripe_price_id, product_name,
			   product_description, price_amount, price_currency, price_type,
			   recurring_interval, recurring_interval_count, active, created_at, updated_at
		FROM stripe_products
		WHERE stripe_price_id = $1
	`

	var product StripeProduct
	err := r.db.Get(&product, query, stripePriceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found, not an error
		}
		return nil, fmt.Errorf("failed to get stripe product by stripe price ID: %w", err)
	}

	return &product, nil
}

// Update updates a Stripe product record
func (r *StripeProductRepository) Update(product *StripeProduct) error {
	query := `
		UPDATE stripe_products
		SET product_name = $1, product_description = $2, price_amount = $3,
			price_currency = $4, active = $5, updated_at = $6
		WHERE id = $7
	`

	product.UpdatedAt = time.Now()

	result, err := r.db.Exec(query,
		product.ProductName,
		product.ProductDescription,
		product.PriceAmount,
		product.PriceCurrency,
		product.Active,
		product.UpdatedAt,
		product.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update stripe product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stripe product not found")
	}

	return nil
}

// Delete deletes a Stripe product record (soft delete by setting active = false)
func (r *StripeProductRepository) Delete(id uuid.UUID) error {
	query := `
		UPDATE stripe_products
		SET active = false, updated_at = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete stripe product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stripe product not found")
	}

	return nil
}

// List retrieves all active Stripe products with pagination
func (r *StripeProductRepository) List(limit, offset int) ([]*StripeProduct, error) {
	query := `
		SELECT id, course_id, stripe_product_id, stripe_price_id, product_name,
			   product_description, price_amount, price_currency, price_type,
			   recurring_interval, recurring_interval_count, active, created_at, updated_at
		FROM stripe_products
		WHERE active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var products []*StripeProduct
	err := r.db.Select(&products, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stripe products: %w", err)
	}

	return products, nil
}