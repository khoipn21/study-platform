package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/study-platform/course-service/internal/model"
	"github.com/study-platform/pkg/database"
)

type CourseRepository struct {
	db *database.DB
}

func NewCourseRepository(db *database.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) Create(ctx context.Context, course *model.Course) error {
	query := `
		INSERT INTO courses (
			id, title, description, instructor_id, instructor_name, category, level, price, currency,
			thumbnail_url, status, tags, difficulty_level, language, learning_outcomes, requirements,
			estimated_duration_hours, auto_approve_enrollment, allow_previews, has_certificate,
			mobile_access, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
	`

	course.ID = uuid.New()
	course.CreatedAt = time.Now()
	course.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		course.ID, course.Title, course.Description, course.InstructorID, course.InstructorName,
		course.Category, course.Level, course.Price, course.Currency, course.ThumbnailURL,
		course.Status, pq.Array(course.Tags), course.DifficultyLevel, course.Language,
		pq.Array(course.LearningOutcomes), pq.Array(course.Requirements), course.EstimatedDurationHours,
		course.AutoApproveEnrollment, course.AllowPreviews, course.HasCertificate, course.MobileAccess,
		course.CreatedAt, course.UpdatedAt,
	)

	return err
}

func (r *CourseRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Course, error) {
	query := `
		SELECT id, title, description, instructor_id, instructor_name, category, level, price, currency,
			   thumbnail_url, status, duration_minutes, enrollment_count, rating, rating_count, tags,
			   difficulty_level, language, learning_outcomes, requirements, estimated_duration_hours,
			   auto_approve_enrollment, allow_previews, has_certificate, mobile_access,
			   created_at, updated_at, deleted_at
		FROM courses
		WHERE id = $1 AND deleted_at IS NULL
	`

	course := &model.Course{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&course.ID, &course.Title, &course.Description, &course.InstructorID, &course.InstructorName,
		&course.Category, &course.Level, &course.Price, &course.Currency, &course.ThumbnailURL,
		&course.Status, &course.DurationMinutes, &course.EnrollmentCount, &course.Rating,
		&course.RatingCount, pq.Array(&course.Tags), &course.DifficultyLevel, &course.Language,
		pq.Array(&course.LearningOutcomes), pq.Array(&course.Requirements), &course.EstimatedDurationHours,
		&course.AutoApproveEnrollment, &course.AllowPreviews, &course.HasCertificate, &course.MobileAccess,
		&course.CreatedAt, &course.UpdatedAt, &course.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("course not found")
		}
		return nil, err
	}

	return course, nil
}

func (r *CourseRepository) Update(ctx context.Context, course *model.Course) error {
	query := `
		UPDATE courses
		SET title = $2, description = $3, category = $4, level = $5, price = $6, currency = $7,
			thumbnail_url = $8, status = $9, tags = $10, updated_at = $11
		WHERE id = $1
	`
	
	course.UpdatedAt = time.Now()
	
	result, err := r.db.ExecContext(ctx, query,
		course.ID, course.Title, course.Description, course.Category, course.Level,
		course.Price, course.Currency, course.ThumbnailURL, course.Status,
		pq.Array(course.Tags), course.UpdatedAt,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("course not found")
	}
	
	return nil
}

func (r *CourseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Soft delete implementation
	query := `
		UPDATE courses
		SET deleted_at = $2, updated_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, id, now)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("course not found or already deleted")
	}

	return nil
}

// SoftDeleteCourseWithCascade performs soft delete on course and cascades to lectures and enrollments
func (r *CourseRepository) SoftDeleteCourseWithCascade(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	now := time.Now()

	// First, soft delete the course
	courseQuery := `
		UPDATE courses
		SET deleted_at = $2, updated_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := tx.ExecContext(ctx, courseQuery, id, now)
	if err != nil {
		return fmt.Errorf("failed to soft delete course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check course deletion: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("course not found or already deleted")
	}

	// Soft delete all related lectures
	lectureQuery := `
		UPDATE lectures
		SET deleted_at = $2, updated_at = $2
		WHERE course_id = $1 AND deleted_at IS NULL
	`
	_, err = tx.ExecContext(ctx, lectureQuery, id, now)
	if err != nil {
		return fmt.Errorf("failed to soft delete lectures: %w", err)
	}

	// Soft delete all related enrollments (preserving billing history)
	enrollmentQuery := `
		UPDATE enrollments
		SET deleted_at = $2
		WHERE course_id = $1 AND deleted_at IS NULL
	`
	_, err = tx.ExecContext(ctx, enrollmentQuery, id, now)
	if err != nil {
		return fmt.Errorf("failed to soft delete enrollments: %w", err)
	}

	return tx.Commit()
}

func (r *CourseRepository) List(ctx context.Context, filters model.CourseFilters) (*model.CourseSearchResult, error) {
	var conditions []string
	var args []interface{}
	argCount := 0

	baseQuery := `
		SELECT id, title, description, instructor_id, instructor_name, category, level, price, currency,
			   thumbnail_url, status, duration_minutes, enrollment_count, rating, rating_count, tags, created_at, updated_at, deleted_at
		FROM courses
	`

	// Always exclude soft-deleted courses
	conditions = append(conditions, "deleted_at IS NULL")
	
	if filters.Category != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("category = $%d", argCount))
		args = append(args, filters.Category)
	}
	
	if filters.Level != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("level = $%d", argCount))
		args = append(args, filters.Level)
	}
	
	if filters.Status != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, filters.Status)
	}
	
	if filters.InstructorID != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("instructor_id = $%d", argCount))
		args = append(args, filters.InstructorID)
	}
	
	if filters.MinPrice > 0 {
		argCount++
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argCount))
		args = append(args, filters.MinPrice)
	}
	
	if filters.MaxPrice > 0 {
		argCount++
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argCount))
		args = append(args, filters.MaxPrice)
	}
	
	if filters.MinRating > 0 {
		argCount++
		conditions = append(conditions, fmt.Sprintf("rating >= $%d", argCount))
		args = append(args, filters.MinRating)
	}
	
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	
	// Count total results
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM courses %s", whereClause)
	var totalCount int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, err
	}
	
	// Get paginated results
	query := fmt.Sprintf("%s %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
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
	
	var courses []model.Course
	for rows.Next() {
		var course model.Course
		err := rows.Scan(
			&course.ID, &course.Title, &course.Description, &course.InstructorID, &course.InstructorName,
			&course.Category, &course.Level, &course.Price, &course.Currency, &course.ThumbnailURL,
			&course.Status, &course.DurationMinutes, &course.EnrollmentCount, &course.Rating,
			&course.RatingCount, pq.Array(&course.Tags), &course.CreatedAt, &course.UpdatedAt, &course.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	
	return &model.CourseSearchResult{
		Courses:    courses,
		TotalCount: totalCount,
		Page:       filters.Page,
		PageSize:   limit,
	}, nil
}

func (r *CourseRepository) Search(ctx context.Context, filters model.CourseFilters) (*model.CourseSearchResult, error) {
	var conditions []string
	var args []interface{}
	argCount := 0

	baseQuery := `
		SELECT id, title, description, instructor_id, instructor_name, category, level, price, currency,
			   thumbnail_url, status, duration_minutes, enrollment_count, rating, rating_count, tags, created_at, updated_at, deleted_at
		FROM courses
	`

	// Always exclude soft-deleted courses
	conditions = append(conditions, "deleted_at IS NULL")
	
	if filters.Query != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d OR category ILIKE $%d)", argCount, argCount, argCount))
		args = append(args, "%"+filters.Query+"%")
	}
	
	if filters.Category != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("category = $%d", argCount))
		args = append(args, filters.Category)
	}
	
	if filters.Level != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("level = $%d", argCount))
		args = append(args, filters.Level)
	}
	
	if filters.MinPrice > 0 {
		argCount++
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argCount))
		args = append(args, filters.MinPrice)
	}
	
	if filters.MaxPrice > 0 {
		argCount++
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argCount))
		args = append(args, filters.MaxPrice)
	}
	
	if filters.MinRating > 0 {
		argCount++
		conditions = append(conditions, fmt.Sprintf("rating >= $%d", argCount))
		args = append(args, filters.MinRating)
	}
	
	// Only show published courses in search
	argCount++
	conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
	args = append(args, "published")
	
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	
	// Count total results
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM courses %s", whereClause)
	var totalCount int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, err
	}
	
	// Get paginated results ordered by relevance (rating * enrollment_count)
	query := fmt.Sprintf("%s %s ORDER BY (rating * enrollment_count) DESC, created_at DESC LIMIT $%d OFFSET $%d",
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
	
	var courses []model.Course
	for rows.Next() {
		var course model.Course
		err := rows.Scan(
			&course.ID, &course.Title, &course.Description, &course.InstructorID, &course.InstructorName,
			&course.Category, &course.Level, &course.Price, &course.Currency, &course.ThumbnailURL,
			&course.Status, &course.DurationMinutes, &course.EnrollmentCount, &course.Rating,
			&course.RatingCount, pq.Array(&course.Tags), &course.CreatedAt, &course.UpdatedAt, &course.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	
	return &model.CourseSearchResult{
		Courses:    courses,
		TotalCount: totalCount,
		Page:       filters.Page,
		PageSize:   limit,
	}, nil
}

func (r *CourseRepository) UpdateEnrollmentCount(ctx context.Context, courseID uuid.UUID) error {
	query := `
		UPDATE courses
		SET enrollment_count = (
			SELECT COUNT(*) FROM enrollments
			WHERE course_id = $1 AND status = 'enrolled'
		)
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, courseID)
	return err
}

func (r *CourseRepository) UpdateRating(ctx context.Context, courseID uuid.UUID, rating float64, ratingCount int32) error {
	query := `
		UPDATE courses
		SET rating = $2, rating_count = $3, updated_at = $4
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, courseID, rating, ratingCount, time.Now())
	return err
}