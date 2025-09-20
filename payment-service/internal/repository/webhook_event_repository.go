package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/study-platform/payment-service/internal/model"

	"github.com/jmoiron/sqlx"
)

type WebhookEventRepository struct {
	db *sqlx.DB
}

func NewWebhookEventRepository(db *sqlx.DB) *WebhookEventRepository {
	return &WebhookEventRepository{db: db}
}

func (r *WebhookEventRepository) Create(event *model.LemonSqueezyWebhookEvent) error {
	payloadJSON, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO lemon_squeezy_webhook_events (
			id, event_id, event_name, processed_at, payload, signature, created_at
		) VALUES (
			:id, :event_id, :event_name, :processed_at, :payload, :signature, :created_at
		)`

	params := map[string]interface{}{
		"id":           event.ID,
		"event_id":     event.EventID,
		"event_name":   event.EventName,
		"processed_at": event.ProcessedAt,
		"payload":      payloadJSON,
		"signature":    event.Signature,
		"created_at":   event.CreatedAt,
	}

	_, err = r.db.NamedExec(query, params)
	return err
}

func (r *WebhookEventRepository) GetByID(id string) (*model.LemonSqueezyWebhookEvent, error) {
	query := `
		SELECT id, event_id, event_name, processed_at, payload, signature, created_at
		FROM lemon_squeezy_webhook_events
		WHERE id = $1`

	var event model.LemonSqueezyWebhookEvent
	var payloadJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&event.ID,
		&event.EventID,
		&event.EventName,
		&event.ProcessedAt,
		&payloadJSON,
		&event.Signature,
		&event.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(payloadJSON, &event.Payload); err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *WebhookEventRepository) GetByEventID(eventID string) (*model.LemonSqueezyWebhookEvent, error) {
	query := `
		SELECT id, event_id, event_name, processed_at, payload, signature, created_at
		FROM lemon_squeezy_webhook_events
		WHERE event_id = $1`

	var event model.LemonSqueezyWebhookEvent
	var payloadJSON []byte

	err := r.db.QueryRow(query, eventID).Scan(
		&event.ID,
		&event.EventID,
		&event.EventName,
		&event.ProcessedAt,
		&payloadJSON,
		&event.Signature,
		&event.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(payloadJSON, &event.Payload); err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *WebhookEventRepository) EventExists(eventID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM lemon_squeezy_webhook_events WHERE event_id = $1)`

	var exists bool
	err := r.db.QueryRow(query, eventID).Scan(&exists)
	return exists, err
}

func (r *WebhookEventRepository) Update(event *model.LemonSqueezyWebhookEvent) error {
	payloadJSON, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	query := `
		UPDATE lemon_squeezy_webhook_events
		SET event_name = :event_name,
		    processed_at = :processed_at,
		    payload = :payload,
		    signature = :signature
		WHERE id = :id`

	params := map[string]interface{}{
		"id":           event.ID,
		"event_name":   event.EventName,
		"processed_at": event.ProcessedAt,
		"payload":      payloadJSON,
		"signature":    event.Signature,
	}

	_, err = r.db.NamedExec(query, params)
	return err
}

func (r *WebhookEventRepository) ListByEventName(eventName string, limit, offset int) ([]*model.LemonSqueezyWebhookEvent, error) {
	query := `
		SELECT id, event_id, event_name, processed_at, payload, signature, created_at
		FROM lemon_squeezy_webhook_events
		WHERE event_name = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, eventName, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.LemonSqueezyWebhookEvent

	for rows.Next() {
		var event model.LemonSqueezyWebhookEvent
		var payloadJSON []byte

		err := rows.Scan(
			&event.ID,
			&event.EventID,
			&event.EventName,
			&event.ProcessedAt,
			&payloadJSON,
			&event.Signature,
			&event.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(payloadJSON, &event.Payload); err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	return events, rows.Err()
}

func (r *WebhookEventRepository) ListUnprocessed(limit int) ([]*model.LemonSqueezyWebhookEvent, error) {
	query := `
		SELECT id, event_id, event_name, processed_at, payload, signature, created_at
		FROM lemon_squeezy_webhook_events
		WHERE processed_at IS NULL
		ORDER BY created_at ASC
		LIMIT $1`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.LemonSqueezyWebhookEvent

	for rows.Next() {
		var event model.LemonSqueezyWebhookEvent
		var payloadJSON []byte

		err := rows.Scan(
			&event.ID,
			&event.EventID,
			&event.EventName,
			&event.ProcessedAt,
			&payloadJSON,
			&event.Signature,
			&event.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(payloadJSON, &event.Payload); err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	return events, rows.Err()
}

func (r *WebhookEventRepository) MarkAsProcessed(eventID string) error {
	query := `
		UPDATE lemon_squeezy_webhook_events
		SET processed_at = $1
		WHERE event_id = $2`

	_, err := r.db.Exec(query, time.Now(), eventID)
	return err
}

func (r *WebhookEventRepository) DeleteOldEvents(olderThanDays int) error {
	query := `
		DELETE FROM lemon_squeezy_webhook_events
		WHERE created_at < $1`

	cutoffDate := time.Now().AddDate(0, 0, -olderThanDays)
	_, err := r.db.Exec(query, cutoffDate)
	return err
}