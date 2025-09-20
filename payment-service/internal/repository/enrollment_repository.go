package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/study-platform/payment-service/internal/model"
)

type EnrollmentRepository struct {
	db *sql.DB
}

func NewEnrollmentRepository(db *sql.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) Create(ctx context.Context, enrollment *model.Enrollment) error {
	query := `
		INSERT INTO enrollments (
			id, user_id, course_id, status, payment_status, payment_verified_at,
			transaction_id, enrolled_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.ExecContext(ctx, query,
		enrollment.ID, enrollment.UserID, enrollment.CourseID, enrollment.Status,
		enrollment.PaymentStatus, enrollment.PaymentVerifiedAt, enrollment.TransactionID,
		enrollment.EnrolledAt, enrollment.CreatedAt, enrollment.UpdatedAt)
	return err
}

func (r *EnrollmentRepository) CreateWithTx(ctx context.Context, tx interface{}, enrollment *model.Enrollment) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	query := `
		INSERT INTO enrollments (
			id, user_id, course_id, status, payment_status, payment_verified_at,
			transaction_id, enrolled_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := sqlTx.ExecContext(ctx, query,
		enrollment.ID, enrollment.UserID, enrollment.CourseID, enrollment.Status,
		enrollment.PaymentStatus, enrollment.PaymentVerifiedAt, enrollment.TransactionID,
		enrollment.EnrolledAt, enrollment.CreatedAt, enrollment.UpdatedAt)
	return err
}

func (r *EnrollmentRepository) GetByUserAndCourse(ctx context.Context, userID, courseID string) (*model.Enrollment, error) {
	enrollment := &model.Enrollment{}
	query := `
		SELECT id, user_id, course_id, status, payment_status, payment_verified_at,
		       transaction_id, enrolled_at, created_at, updated_at
		FROM enrollments WHERE user_id = $1 AND course_id = $2`

	err := r.db.QueryRowContext(ctx, query, userID, courseID).Scan(
		&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status,
		&enrollment.PaymentStatus, &enrollment.PaymentVerifiedAt, &enrollment.TransactionID,
		&enrollment.EnrolledAt, &enrollment.CreatedAt, &enrollment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return enrollment, nil
}

func (r *EnrollmentRepository) Update(ctx context.Context, enrollment *model.Enrollment) error {
	query := `
		UPDATE enrollments
		SET status = $3, payment_status = $4, payment_verified_at = $5,
		    transaction_id = $6, updated_at = $7
		WHERE user_id = $1 AND course_id = $2`

	result, err := r.db.ExecContext(ctx, query,
		enrollment.UserID, enrollment.CourseID, enrollment.Status,
		enrollment.PaymentStatus, enrollment.PaymentVerifiedAt,
		enrollment.TransactionID, time.Now())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("enrollment not found")
	}

	return nil
}

func (r *EnrollmentRepository) UpdateWithTx(ctx context.Context, tx interface{}, enrollment *model.Enrollment) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	query := `
		UPDATE enrollments
		SET status = $3, payment_status = $4, payment_verified_at = $5,
		    transaction_id = $6, updated_at = $7
		WHERE user_id = $1 AND course_id = $2`

	result, err := sqlTx.ExecContext(ctx, query,
		enrollment.UserID, enrollment.CourseID, enrollment.Status,
		enrollment.PaymentStatus, enrollment.PaymentVerifiedAt,
		enrollment.TransactionID, time.Now())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("enrollment not found")
	}

	return nil
}

func (r *EnrollmentRepository) GetByUserID(ctx context.Context, userID string) ([]*model.Enrollment, error) {
	query := `
		SELECT id, user_id, course_id, status, payment_status, payment_verified_at,
		       transaction_id, enrolled_at, created_at, updated_at
		FROM enrollments WHERE user_id = $1 ORDER BY enrolled_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []*model.Enrollment
	for rows.Next() {
		enrollment := &model.Enrollment{}
		err := rows.Scan(
			&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status,
			&enrollment.PaymentStatus, &enrollment.PaymentVerifiedAt, &enrollment.TransactionID,
			&enrollment.EnrolledAt, &enrollment.CreatedAt, &enrollment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}

	return enrollments, nil
}

func (r *EnrollmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Enrollment, error) {
	enrollment := &model.Enrollment{}
	query := `
		SELECT id, user_id, course_id, status, payment_status, payment_verified_at,
		       transaction_id, enrolled_at, created_at, updated_at
		FROM enrollments WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status,
		&enrollment.PaymentStatus, &enrollment.PaymentVerifiedAt, &enrollment.TransactionID,
		&enrollment.EnrolledAt, &enrollment.CreatedAt, &enrollment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return enrollment, nil
}

// RemoveEnrollment removes an enrollment (soft delete)
func (r *EnrollmentRepository) RemoveEnrollment(ctx context.Context, userID, courseID string) error {
	query := `
		UPDATE enrollments
		SET status = 'cancelled', updated_at = NOW()
		WHERE user_id = $1 AND course_id = $2`

	result, err := r.db.ExecContext(ctx, query, userID, courseID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("enrollment not found")
	}

	return nil
}