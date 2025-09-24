package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// StripeWebhookEvent represents a Stripe webhook event record
type StripeWebhookEvent struct {
	ID                  uuid.UUID  `db:"id" json:"id"`
	StripeEventID       string     `db:"stripe_event_id" json:"stripe_event_id"`
	EventType           string     `db:"event_type" json:"event_type"`
	Processed           *bool      `db:"processed" json:"processed"`
	ProcessingAttempts  *int       `db:"processing_attempts" json:"processing_attempts"`
	EventData           []byte     `db:"event_data" json:"event_data"` // JSONB data
	CreatedAt           time.Time  `db:"created_at" json:"created_at"`
	ProcessedAt         *time.Time `db:"processed_at" json:"processed_at"`
	ErrorMessage        *string    `db:"error_message" json:"error_message"`
}

// StripeWebhookRepository handles database operations for Stripe webhook events
type StripeWebhookRepository struct {
	db *sqlx.DB
}

// NewStripeWebhookRepository creates a new repository
func NewStripeWebhookRepository(db *sqlx.DB) *StripeWebhookRepository {
	return &StripeWebhookRepository{db: db}
}

// Create creates a new Stripe webhook event record
func (r *StripeWebhookRepository) Create(event *StripeWebhookEvent) error {
	query := `
		INSERT INTO stripe_webhook_events (
			id, stripe_event_id, event_type, processed, processing_attempts,
			event_data, created_at, processed_at, error_message
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	event.ID = uuid.New()
	event.CreatedAt = time.Now()

	// Set defaults
	if event.Processed == nil {
		processed := false
		event.Processed = &processed
	}
	if event.ProcessingAttempts == nil {
		attempts := 0
		event.ProcessingAttempts = &attempts
	}

	_, err := r.db.Exec(query,
		event.ID,
		event.StripeEventID,
		event.EventType,
		event.Processed,
		event.ProcessingAttempts,
		event.EventData,
		event.CreatedAt,
		event.ProcessedAt,
		event.ErrorMessage,
	)

	if err != nil {
		return fmt.Errorf("failed to create stripe webhook event: %w", err)
	}

	return nil
}

// GetByStripeEventID retrieves a webhook event by Stripe event ID
func (r *StripeWebhookRepository) GetByStripeEventID(stripeEventID string) (*StripeWebhookEvent, error) {
	query := `
		SELECT id, stripe_event_id, event_type, processed, processing_attempts,
			   event_data, created_at, processed_at, error_message
		FROM stripe_webhook_events
		WHERE stripe_event_id = $1
	`

	var event StripeWebhookEvent
	err := r.db.Get(&event, query, stripeEventID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found, not an error
		}
		return nil, fmt.Errorf("failed to get stripe webhook event: %w", err)
	}

	return &event, nil
}

// MarkAsProcessed marks a webhook event as processed
func (r *StripeWebhookRepository) MarkAsProcessed(id uuid.UUID) error {
	query := `
		UPDATE stripe_webhook_events
		SET processed = true, processed_at = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to mark stripe webhook event as processed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stripe webhook event not found")
	}

	return nil
}

// UpdateProcessingAttempt increments processing attempts and updates error message
func (r *StripeWebhookRepository) UpdateProcessingAttempt(id uuid.UUID, errorMessage string) error {
	query := `
		UPDATE stripe_webhook_events
		SET processing_attempts = processing_attempts + 1,
			error_message = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(query, errorMessage, id)
	if err != nil {
		return fmt.Errorf("failed to update stripe webhook event processing attempt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stripe webhook event not found")
	}

	return nil
}

// GetUnprocessedEvents retrieves unprocessed webhook events
func (r *StripeWebhookRepository) GetUnprocessedEvents(limit int) ([]*StripeWebhookEvent, error) {
	query := `
		SELECT id, stripe_event_id, event_type, processed, processing_attempts,
			   event_data, created_at, processed_at, error_message
		FROM stripe_webhook_events
		WHERE processed = false AND processing_attempts < 3
		ORDER BY created_at ASC
		LIMIT $1
	`

	var events []*StripeWebhookEvent
	err := r.db.Select(&events, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get unprocessed stripe webhook events: %w", err)
	}

	return events, nil
}

// GetFailedEvents retrieves failed webhook events (3+ processing attempts)
func (r *StripeWebhookRepository) GetFailedEvents(limit, offset int) ([]*StripeWebhookEvent, error) {
	query := `
		SELECT id, stripe_event_id, event_type, processed, processing_attempts,
			   event_data, created_at, processed_at, error_message
		FROM stripe_webhook_events
		WHERE processed = false AND processing_attempts >= 3
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var events []*StripeWebhookEvent
	err := r.db.Select(&events, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get failed stripe webhook events: %w", err)
	}

	return events, nil
}

// List retrieves webhook events with pagination
func (r *StripeWebhookRepository) List(limit, offset int) ([]*StripeWebhookEvent, error) {
	query := `
		SELECT id, stripe_event_id, event_type, processed, processing_attempts,
			   event_data, created_at, processed_at, error_message
		FROM stripe_webhook_events
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var events []*StripeWebhookEvent
	err := r.db.Select(&events, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stripe webhook events: %w", err)
	}

	return events, nil
}

// Delete deletes a webhook event record
func (r *StripeWebhookRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM stripe_webhook_events WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete stripe webhook event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stripe webhook event not found")
	}

	return nil
}