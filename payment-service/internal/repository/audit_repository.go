package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/study-platform/payment-service/internal/model"
)

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(ctx context.Context, audit *model.AuditLog) error {
	detailsJSON, err := json.Marshal(audit.Details)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO audit_logs (
			id, action, user_id, course_id, transaction_id, details, ip_address, timestamp
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = r.db.ExecContext(ctx, query,
		audit.ID, audit.Action, audit.UserID, audit.CourseID,
		audit.TransactionID, detailsJSON, audit.IPAddress, audit.Timestamp)
	return err
}

func (r *AuditRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.AuditLog, error) {
	query := `
		SELECT id, action, user_id, course_id, transaction_id, details, ip_address, timestamp
		FROM audit_logs WHERE user_id = $1 ORDER BY timestamp DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []*model.AuditLog
	for rows.Next() {
		audit := &model.AuditLog{}
		var detailsJSON []byte

		err := rows.Scan(
			&audit.ID, &audit.Action, &audit.UserID, &audit.CourseID,
			&audit.TransactionID, &detailsJSON, &audit.IPAddress, &audit.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		if detailsJSON != nil {
			if err := json.Unmarshal(detailsJSON, &audit.Details); err != nil {
				return nil, err
			}
		}

		audits = append(audits, audit)
	}

	return audits, nil
}

func (r *AuditRepository) GetByAction(ctx context.Context, action string, limit, offset int) ([]*model.AuditLog, error) {
	query := `
		SELECT id, action, user_id, course_id, transaction_id, details, ip_address, timestamp
		FROM audit_logs WHERE action = $1 ORDER BY timestamp DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, action, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []*model.AuditLog
	for rows.Next() {
		audit := &model.AuditLog{}
		var detailsJSON []byte

		err := rows.Scan(
			&audit.ID, &audit.Action, &audit.UserID, &audit.CourseID,
			&audit.TransactionID, &detailsJSON, &audit.IPAddress, &audit.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		if detailsJSON != nil {
			if err := json.Unmarshal(detailsJSON, &audit.Details); err != nil {
				return nil, err
			}
		}

		audits = append(audits, audit)
	}

	return audits, nil
}

func (r *AuditRepository) GetByTransactionID(ctx context.Context, transactionID string) ([]*model.AuditLog, error) {
	query := `
		SELECT id, action, user_id, course_id, transaction_id, details, ip_address, timestamp
		FROM audit_logs WHERE transaction_id = $1 ORDER BY timestamp DESC`

	rows, err := r.db.QueryContext(ctx, query, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []*model.AuditLog
	for rows.Next() {
		audit := &model.AuditLog{}
		var detailsJSON []byte

		err := rows.Scan(
			&audit.ID, &audit.Action, &audit.UserID, &audit.CourseID,
			&audit.TransactionID, &detailsJSON, &audit.IPAddress, &audit.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		if detailsJSON != nil {
			if err := json.Unmarshal(detailsJSON, &audit.Details); err != nil {
				return nil, err
			}
		}

		audits = append(audits, audit)
	}

	return audits, nil
}

func (r *AuditRepository) GetByCourseID(ctx context.Context, courseID string, limit, offset int) ([]*model.AuditLog, error) {
	query := `
		SELECT id, action, user_id, course_id, transaction_id, details, ip_address, timestamp
		FROM audit_logs WHERE course_id = $1 ORDER BY timestamp DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, courseID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []*model.AuditLog
	for rows.Next() {
		audit := &model.AuditLog{}
		var detailsJSON []byte

		err := rows.Scan(
			&audit.ID, &audit.Action, &audit.UserID, &audit.CourseID,
			&audit.TransactionID, &detailsJSON, &audit.IPAddress, &audit.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		if detailsJSON != nil {
			if err := json.Unmarshal(detailsJSON, &audit.Details); err != nil {
				return nil, err
			}
		}

		audits = append(audits, audit)
	}

	return audits, nil
}

func (r *AuditRepository) GetByTimeRange(ctx context.Context, startTime, endTime time.Time, limit, offset int) ([]*model.AuditLog, error) {
	query := `
		SELECT id, action, user_id, course_id, transaction_id, details, ip_address, timestamp
		FROM audit_logs WHERE timestamp >= $1 AND timestamp <= $2
		ORDER BY timestamp DESC LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, startTime, endTime, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []*model.AuditLog
	for rows.Next() {
		audit := &model.AuditLog{}
		var detailsJSON []byte

		err := rows.Scan(
			&audit.ID, &audit.Action, &audit.UserID, &audit.CourseID,
			&audit.TransactionID, &detailsJSON, &audit.IPAddress, &audit.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		if detailsJSON != nil {
			if err := json.Unmarshal(detailsJSON, &audit.Details); err != nil {
				return nil, err
			}
		}

		audits = append(audits, audit)
	}

	return audits, nil
}
// LogCheckoutCreated logs a checkout creation event
func (r *AuditRepository) LogCheckoutCreated(ctx context.Context, userID, courseID, checkoutID string, details map[string]interface{}) error {
	audit := &model.AuditLog{
		ID:            uuid.New().String(),
		Action:        "checkout_created",
		UserID:        userID,
		CourseID:      &courseID,
		TransactionID: &checkoutID,
		Details:       details,
		IPAddress:     "", // TODO: Get from context if available
		Timestamp:     time.Now(),
	}
	return r.Create(ctx, audit)
}

// LogEnrollmentCreated logs an enrollment creation event
func (r *AuditRepository) LogEnrollmentCreated(ctx context.Context, userID, courseID, enrollmentID string, details map[string]interface{}) error {
	audit := &model.AuditLog{
		ID:            uuid.New().String(),
		Action:        "enrollment_created",
		UserID:        userID,
		CourseID:      &courseID,
		TransactionID: &enrollmentID,
		Details:       details,
		IPAddress:     "", // TODO: Get from context if available
		Timestamp:     time.Now(),
	}
	return r.Create(ctx, audit)
}

// LogCheckoutFailed logs a checkout failure event
func (r *AuditRepository) LogCheckoutFailed(ctx context.Context, userID, courseID, checkoutID string, details map[string]interface{}) error {
	audit := &model.AuditLog{
		ID:            uuid.New().String(),
		Action:        "checkout_failed",
		UserID:        userID,
		CourseID:      &courseID,
		TransactionID: &checkoutID,
		Details:       details,
		IPAddress:     "", // TODO: Get from context if available
		Timestamp:     time.Now(),
	}
	return r.Create(ctx, audit)
}
