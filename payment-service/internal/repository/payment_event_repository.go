package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/study-platform/payment-service/internal/model"
)

type PaymentEventRepository struct {
	db *sql.DB
}

func NewPaymentEventRepository(db *sql.DB) *PaymentEventRepository {
	return &PaymentEventRepository{db: db}
}

func (r *PaymentEventRepository) Create(ctx context.Context, event *model.PaymentEvent) error {
	payloadJSON, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO payment_events (
			id, event_type, provider, provider_event_id, payload, processed,
			processed_at, error_message, retry_count, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = r.db.ExecContext(ctx, query,
		event.ID, event.EventType, event.Provider, event.ProviderEventID,
		payloadJSON, event.Processed, event.ProcessedAt, event.ErrorMessage,
		event.RetryCount, event.CreatedAt, event.UpdatedAt)
	return err
}

func (r *PaymentEventRepository) Update(ctx context.Context, event *model.PaymentEvent) error {
	payloadJSON, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	query := `
		UPDATE payment_events
		SET event_type = $2, provider = $3, provider_event_id = $4, payload = $5,
		    processed = $6, processed_at = $7, error_message = $8, retry_count = $9,
		    updated_at = $10
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query,
		event.ID, event.EventType, event.Provider, event.ProviderEventID,
		payloadJSON, event.Processed, event.ProcessedAt, event.ErrorMessage,
		event.RetryCount, time.Now())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment event not found")
	}

	return nil
}

func (r *PaymentEventRepository) GetByID(ctx context.Context, id string) (*model.PaymentEvent, error) {
	event := &model.PaymentEvent{}
	var payloadJSON []byte

	query := `
		SELECT id, event_type, provider, provider_event_id, payload, processed,
		       processed_at, error_message, retry_count, created_at, updated_at
		FROM payment_events WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID, &event.EventType, &event.Provider, &event.ProviderEventID,
		&payloadJSON, &event.Processed, &event.ProcessedAt, &event.ErrorMessage,
		&event.RetryCount, &event.CreatedAt, &event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if payloadJSON != nil {
		if err := json.Unmarshal(payloadJSON, &event.Payload); err != nil {
			return nil, err
		}
	}

	return event, nil
}

func (r *PaymentEventRepository) GetByProviderEventID(ctx context.Context, providerEventID string) (*model.PaymentEvent, error) {
	event := &model.PaymentEvent{}
	var payloadJSON []byte

	query := `
		SELECT id, event_type, provider, provider_event_id, payload, processed,
		       processed_at, error_message, retry_count, created_at, updated_at
		FROM payment_events WHERE provider_event_id = $1`

	err := r.db.QueryRowContext(ctx, query, providerEventID).Scan(
		&event.ID, &event.EventType, &event.Provider, &event.ProviderEventID,
		&payloadJSON, &event.Processed, &event.ProcessedAt, &event.ErrorMessage,
		&event.RetryCount, &event.CreatedAt, &event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if payloadJSON != nil {
		if err := json.Unmarshal(payloadJSON, &event.Payload); err != nil {
			return nil, err
		}
	}

	return event, nil
}

func (r *PaymentEventRepository) GetUnprocessed(ctx context.Context, limit int) ([]*model.PaymentEvent, error) {
	query := `
		SELECT id, event_type, provider, provider_event_id, payload, processed,
		       processed_at, error_message, retry_count, created_at, updated_at
		FROM payment_events WHERE processed = false ORDER BY created_at ASC LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.PaymentEvent
	for rows.Next() {
		event := &model.PaymentEvent{}
		var payloadJSON []byte

		err := rows.Scan(
			&event.ID, &event.EventType, &event.Provider, &event.ProviderEventID,
			&payloadJSON, &event.Processed, &event.ProcessedAt, &event.ErrorMessage,
			&event.RetryCount, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if payloadJSON != nil {
			if err := json.Unmarshal(payloadJSON, &event.Payload); err != nil {
				return nil, err
			}
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *PaymentEventRepository) GetFailedEvents(ctx context.Context, maxRetries int, limit int) ([]*model.PaymentEvent, error) {
	query := `
		SELECT id, event_type, provider, provider_event_id, payload, processed,
		       processed_at, error_message, retry_count, created_at, updated_at
		FROM payment_events
		WHERE processed = false AND retry_count < $1 AND error_message != ''
		ORDER BY retry_count ASC, created_at ASC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, maxRetries, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.PaymentEvent
	for rows.Next() {
		event := &model.PaymentEvent{}
		var payloadJSON []byte

		err := rows.Scan(
			&event.ID, &event.EventType, &event.Provider, &event.ProviderEventID,
			&payloadJSON, &event.Processed, &event.ProcessedAt, &event.ErrorMessage,
			&event.RetryCount, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if payloadJSON != nil {
			if err := json.Unmarshal(payloadJSON, &event.Payload); err != nil {
				return nil, err
			}
		}

		events = append(events, event)
	}

	return events, nil
}