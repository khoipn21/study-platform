package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/study-platform/payment-service/internal/model"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) GetByID(ctx context.Context, id string) (*model.Course, error) {
	course := &model.Course{}
	query := `
		SELECT id, title, description, instructor_id, price, currency, status, is_paid,
		       lemon_squeezy_product_id, lemon_squeezy_variant_id, created_at, updated_at
		FROM courses WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&course.ID, &course.Title, &course.Description, &course.InstructorID,
		&course.Price, &course.Currency, &course.Status, &course.IsPaid,
		&course.LemonSqueezyProductID, &course.LemonSqueezyVariantID,
		&course.CreatedAt, &course.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (r *CourseRepository) Create(ctx context.Context, course *model.Course) error {
	query := `
		INSERT INTO courses (
			id, title, description, instructor_id, price, currency, status, is_paid,
			lemon_squeezy_product_id, lemon_squeezy_variant_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.ExecContext(ctx, query,
		course.ID, course.Title, course.Description, course.InstructorID,
		course.Price, course.Currency, course.Status, course.IsPaid,
		course.LemonSqueezyProductID, course.LemonSqueezyVariantID,
		course.CreatedAt, course.UpdatedAt)
	return err
}

func (r *CourseRepository) Update(ctx context.Context, course *model.Course) error {
	query := `
		UPDATE courses
		SET title = $2, description = $3, instructor_id = $4, price = $5, currency = $6,
		    status = $7, is_paid = $8, lemon_squeezy_product_id = $9,
		    lemon_squeezy_variant_id = $10, updated_at = $11
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query,
		course.ID, course.Title, course.Description, course.InstructorID,
		course.Price, course.Currency, course.Status, course.IsPaid,
		course.LemonSqueezyProductID, course.LemonSqueezyVariantID,
		course.UpdatedAt)
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

func (r *CourseRepository) GetByInstructorID(ctx context.Context, instructorID string) ([]*model.Course, error) {
	query := `
		SELECT id, title, description, instructor_id, price, currency, status, is_paid,
		       lemon_squeezy_product_id, lemon_squeezy_variant_id, created_at, updated_at
		FROM courses WHERE instructor_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, instructorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*model.Course
	for rows.Next() {
		course := &model.Course{}
		err := rows.Scan(
			&course.ID, &course.Title, &course.Description, &course.InstructorID,
			&course.Price, &course.Currency, &course.Status, &course.IsPaid,
			&course.LemonSqueezyProductID, &course.LemonSqueezyVariantID,
			&course.CreatedAt, &course.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func (r *CourseRepository) GetPaidCourses(ctx context.Context) ([]*model.Course, error) {
	query := `
		SELECT id, title, description, instructor_id, price, currency, status, is_paid,
		       lemon_squeezy_product_id, lemon_squeezy_variant_id, created_at, updated_at
		FROM courses WHERE is_paid = true AND status = 'published' ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*model.Course
	for rows.Next() {
		course := &model.Course{}
		err := rows.Scan(
			&course.ID, &course.Title, &course.Description, &course.InstructorID,
			&course.Price, &course.Currency, &course.Status, &course.IsPaid,
			&course.LemonSqueezyProductID, &course.LemonSqueezyVariantID,
			&course.CreatedAt, &course.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func (r *CourseRepository) GetByLemonSqueezyVariantID(ctx context.Context, variantID string) (*model.Course, error) {
	course := &model.Course{}
	query := `
		SELECT id, title, description, instructor_id, price, currency, status, is_paid,
		       lemon_squeezy_product_id, lemon_squeezy_variant_id, created_at, updated_at
		FROM courses WHERE lemon_squeezy_variant_id = $1`

	err := r.db.QueryRowContext(ctx, query, variantID).Scan(
		&course.ID, &course.Title, &course.Description, &course.InstructorID,
		&course.Price, &course.Currency, &course.Status, &course.IsPaid,
		&course.LemonSqueezyProductID, &course.LemonSqueezyVariantID,
		&course.CreatedAt, &course.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return course, nil
}

// ExecContext provides direct database access for complex queries
func (r *CourseRepository) ExecContext(ctx context.Context, query string, args ...interface{}) error {
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}