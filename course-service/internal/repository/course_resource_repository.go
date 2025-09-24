package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/study-platform/course-service/internal/model"
	"github.com/study-platform/pkg/logger"
)

type CourseResourceRepository struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewCourseResourceRepository(db *sqlx.DB, logger logger.Logger) *CourseResourceRepository {
	return &CourseResourceRepository{
		db:     db,
		logger: logger,
	}
}

func (r *CourseResourceRepository) Create(ctx context.Context, resource *model.CourseResource) error {
	query := `
		INSERT INTO course_resources (
			course_id, file_id, resource_type, display_order, is_required
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx, query,
		resource.CourseID,
		resource.FileID,
		resource.ResourceType,
		resource.DisplayOrder,
		resource.IsRequired,
	).Scan(&resource.ID, &resource.CreatedAt, &resource.UpdatedAt)

	if err != nil {
		r.logger.Errorf("Failed to create course resource: %v", err)
		return err
	}

	return nil
}

func (r *CourseResourceRepository) GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]model.CourseResource, error) {
	query := `
		SELECT
			cr.id, cr.course_id, cr.file_id, cr.resource_type, cr.display_order,
			cr.is_required, cr.created_at, cr.updated_at,
			f.filename, f.original_filename, f.content_type, f.size_bytes,
			f.is_public, f.created_at as file_created_at,
			CASE
				WHEN f.is_public THEN CONCAT('https://', f.bucket_name, '.s3.', 'ap-southeast-2', '.amazonaws.com/', f.object_key)
				ELSE 'signed-url-required'
			END as download_url
		FROM course_resources cr
		JOIN files f ON cr.file_id = f.id
		WHERE cr.course_id = $1 AND f.deleted_at IS NULL
		ORDER BY cr.display_order ASC`

	var resources []model.CourseResource
	err := r.db.SelectContext(ctx, &resources, query, courseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.CourseResource{}, nil
		}
		r.logger.Errorf("Failed to get course resources: %v", err)
		return nil, err
	}

	return resources, nil
}

func (r *CourseResourceRepository) Update(ctx context.Context, resource *model.CourseResource) error {
	query := `
		UPDATE course_resources
		SET resource_type = $1, display_order = $2, is_required = $3, updated_at = now()
		WHERE id = $4`

	result, err := r.db.ExecContext(
		ctx, query,
		resource.ResourceType,
		resource.DisplayOrder,
		resource.IsRequired,
		resource.ID,
	)

	if err != nil {
		r.logger.Errorf("Failed to update course resource: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("course resource not found")
	}

	return nil
}

func (r *CourseResourceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM course_resources WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Errorf("Failed to delete course resource: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("course resource not found")
	}

	return nil
}

func (r *CourseResourceRepository) DeleteByCourseID(ctx context.Context, courseID uuid.UUID) error {
	query := `DELETE FROM course_resources WHERE course_id = $1`

	_, err := r.db.ExecContext(ctx, query, courseID)
	if err != nil {
		r.logger.Errorf("Failed to delete course resources: %v", err)
		return err
	}

	return nil
}