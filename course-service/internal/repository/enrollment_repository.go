package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/study-platform/course-service/internal/model"
	"github.com/study-platform/pkg/database"
)

type EnrollmentRepository struct {
	db *database.DB
}

func NewEnrollmentRepository(db *database.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) Create(ctx context.Context, enrollment *model.Enrollment) error {
	query := `
		INSERT INTO enrollments (
			id, user_id, course_id, status, progress_percentage,
			payment_required, payment_status, lemon_squeezy_order_id,
			payment_amount, payment_currency, paid_at,
			enrolled_at, last_accessed
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	if enrollment.ID == uuid.Nil {
		enrollment.ID = uuid.New()
	}
	enrollment.EnrolledAt = time.Now()
	if enrollment.ProgressPercentage == 0 {
		enrollment.ProgressPercentage = 0
	}
	now := time.Now()
	enrollment.LastAccessed = &now

	_, err := r.db.ExecContext(ctx, query,
		enrollment.ID, enrollment.UserID, enrollment.CourseID, enrollment.Status,
		enrollment.ProgressPercentage, enrollment.PaymentRequired, enrollment.PaymentStatus,
		enrollment.LemonSqueezyOrderID, enrollment.PaymentAmount, enrollment.PaymentCurrency,
		enrollment.PaidAt, enrollment.EnrolledAt, enrollment.LastAccessed,
	)

	return err
}

func (r *EnrollmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Enrollment, error) {
	query := `
		SELECT
			id, user_id, course_id, status, progress_percentage,
			payment_required, payment_status, lemon_squeezy_order_id,
			payment_amount, payment_currency, paid_at,
			enrolled_at, completed_at, last_accessed
		FROM enrollments
		WHERE id = $1
	`

	enrollment := &model.Enrollment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status,
		&enrollment.ProgressPercentage, &enrollment.PaymentRequired, &enrollment.PaymentStatus,
		&enrollment.LemonSqueezyOrderID, &enrollment.PaymentAmount, &enrollment.PaymentCurrency,
		&enrollment.PaidAt, &enrollment.EnrolledAt, &enrollment.CompletedAt, &enrollment.LastAccessed,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("enrollment not found")
		}
		return nil, err
	}

	return enrollment, nil
}

func (r *EnrollmentRepository) GetByUserAndCourse(ctx context.Context, userID, courseID uuid.UUID) (*model.Enrollment, error) {
	query := `
		SELECT
			id, user_id, course_id, status, progress_percentage,
			payment_required, payment_status, lemon_squeezy_order_id,
			payment_amount, payment_currency, paid_at,
			enrolled_at, completed_at, last_accessed
		FROM enrollments
		WHERE user_id = $1 AND course_id = $2
	`

	enrollment := &model.Enrollment{}
	err := r.db.QueryRowContext(ctx, query, userID, courseID).Scan(
		&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status,
		&enrollment.ProgressPercentage, &enrollment.PaymentRequired, &enrollment.PaymentStatus,
		&enrollment.LemonSqueezyOrderID, &enrollment.PaymentAmount, &enrollment.PaymentCurrency,
		&enrollment.PaidAt, &enrollment.EnrolledAt, &enrollment.CompletedAt, &enrollment.LastAccessed,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("enrollment not found")
		}
		return nil, err
	}

	return enrollment, nil
}

func (r *EnrollmentRepository) Update(ctx context.Context, enrollment *model.Enrollment) error {
	query := `
		UPDATE enrollments
		SET
			status = $2,
			progress_percentage = $3,
			payment_required = $4,
			payment_status = $5,
			lemon_squeezy_order_id = $6,
			payment_amount = $7,
			payment_currency = $8,
			paid_at = $9,
			completed_at = $10,
			last_accessed = $11
		WHERE id = $1
	`

	now := time.Now()
	enrollment.LastAccessed = &now

	// Set completed_at if status is completed
	if enrollment.Status == model.EnrollmentStatusCompleted && enrollment.CompletedAt == nil {
		enrollment.CompletedAt = &now
	}

	result, err := r.db.ExecContext(ctx, query,
		enrollment.ID, enrollment.Status, enrollment.ProgressPercentage,
		enrollment.PaymentRequired, enrollment.PaymentStatus, enrollment.LemonSqueezyOrderID,
		enrollment.PaymentAmount, enrollment.PaymentCurrency, enrollment.PaidAt,
		enrollment.CompletedAt, enrollment.LastAccessed,
	)

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

func (r *EnrollmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM enrollments WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
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

func (r *EnrollmentRepository) List(ctx context.Context, filters model.EnrollmentFilters) (*model.EnrollmentSearchResult, error) {
	var conditions []string
	var args []interface{}
	argCount := 0
	
	baseQuery := `
		SELECT
			id, user_id, course_id, status, progress_percentage,
			payment_required, payment_status, lemon_squeezy_order_id,
			payment_amount, payment_currency, paid_at,
			enrolled_at, completed_at, last_accessed
		FROM enrollments
	`
	
	if filters.UserID != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argCount))
		args = append(args, filters.UserID)
	}
	
	if filters.CourseID != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("course_id = $%d", argCount))
		args = append(args, filters.CourseID)
	}
	
	if filters.Status != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, filters.Status)
	}
	
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	
	// Count total results
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM enrollments %s", whereClause)
	var totalCount int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, err
	}
	
	// Get paginated results
	query := fmt.Sprintf("%s %s ORDER BY enrolled_at DESC LIMIT $%d OFFSET $%d",
		baseQuery, whereClause, argCount+1, argCount+2)
	
	limit := filters.PageSize
	if limit <= 0 {
		limit = 20
	}
	offset := (filters.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	
	args = append(args, limit, offset)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var enrollments []model.Enrollment
	for rows.Next() {
		var enrollment model.Enrollment
		err := rows.Scan(
			&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status,
			&enrollment.ProgressPercentage, &enrollment.PaymentRequired, &enrollment.PaymentStatus,
			&enrollment.LemonSqueezyOrderID, &enrollment.PaymentAmount, &enrollment.PaymentCurrency,
			&enrollment.PaidAt, &enrollment.EnrolledAt, &enrollment.CompletedAt, &enrollment.LastAccessed,
		)
		if err != nil {
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}
	
	return &model.EnrollmentSearchResult{
		Enrollments: enrollments,
		TotalCount:  totalCount,
		Page:        filters.Page,
		PageSize:    limit,
	}, nil
}

func (r *EnrollmentRepository) UpdateProgress(ctx context.Context, userID, courseID uuid.UUID, progressPercentage float64) error {
	query := `
		UPDATE enrollments
		SET progress_percentage = $3, last_accessed = $4, status = CASE 
			WHEN $3 >= 100 THEN 'completed'
			ELSE status
		END, completed_at = CASE 
			WHEN $3 >= 100 AND completed_at IS NULL THEN $4
			ELSE completed_at
		END
		WHERE user_id = $1 AND course_id = $2
	`
	
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, userID, courseID, progressPercentage, now)
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

func (r *EnrollmentRepository) GetUserEnrollments(ctx context.Context, userID uuid.UUID) ([]model.Enrollment, error) {
	query := `
		SELECT
			id, user_id, course_id, status, progress_percentage,
			payment_required, payment_status, lemon_squeezy_order_id,
			payment_amount, payment_currency, paid_at,
			enrolled_at, completed_at, last_accessed
		FROM enrollments
		WHERE user_id = $1
		ORDER BY enrolled_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var enrollments []model.Enrollment
	for rows.Next() {
		var enrollment model.Enrollment
		err := rows.Scan(
			&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status,
			&enrollment.ProgressPercentage, &enrollment.PaymentRequired, &enrollment.PaymentStatus,
			&enrollment.LemonSqueezyOrderID, &enrollment.PaymentAmount, &enrollment.PaymentCurrency,
			&enrollment.PaidAt, &enrollment.EnrolledAt, &enrollment.CompletedAt, &enrollment.LastAccessed,
		)
		if err != nil {
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}
	
	return enrollments, nil
}

func (r *EnrollmentRepository) GetCourseEnrollments(ctx context.Context, courseID uuid.UUID) ([]model.Enrollment, error) {
	query := `
		SELECT
			id, user_id, course_id, status, progress_percentage,
			payment_required, payment_status, lemon_squeezy_order_id,
			payment_amount, payment_currency, paid_at,
			enrolled_at, completed_at, last_accessed
		FROM enrollments
		WHERE course_id = $1
		ORDER BY enrolled_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var enrollments []model.Enrollment
	for rows.Next() {
		var enrollment model.Enrollment
		err := rows.Scan(
			&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status,
			&enrollment.ProgressPercentage, &enrollment.PaymentRequired, &enrollment.PaymentStatus,
			&enrollment.LemonSqueezyOrderID, &enrollment.PaymentAmount, &enrollment.PaymentCurrency,
			&enrollment.PaidAt, &enrollment.EnrolledAt, &enrollment.CompletedAt, &enrollment.LastAccessed,
		)
		if err != nil {
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}
	
	return enrollments, nil
}