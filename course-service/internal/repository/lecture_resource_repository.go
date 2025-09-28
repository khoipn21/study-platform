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

type LectureResourceRepository struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewLectureResourceRepository(db *sqlx.DB, logger logger.Logger) *LectureResourceRepository {
	return &LectureResourceRepository{
		db:     db,
		logger: logger,
	}
}

func (r *LectureResourceRepository) Create(ctx context.Context, resource *model.LectureResource) error {
	query := `
		INSERT INTO lecture_resources (
			lecture_id, file_id, resource_type, display_order, is_required
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx, query,
		resource.LectureID,
		resource.FileID,
		resource.ResourceType,
		resource.DisplayOrder,
		resource.IsRequired,
	).Scan(&resource.ID, &resource.CreatedAt, &resource.UpdatedAt)

	if err != nil {
		r.logger.Errorf("Failed to create lecture resource: %v", err)
		return err
	}

	return nil
}

func (r *LectureResourceRepository) GetByLectureID(ctx context.Context, lectureID uuid.UUID) ([]model.LectureResource, error) {
	query := `
		SELECT
			lr.id, lr.lecture_id, lr.file_id, lr.resource_type, lr.display_order,
			lr.is_required, lr.created_at, lr.updated_at,
			f.filename, f.original_filename, f.content_type, f.size_bytes,
			f.is_public, f.created_at as file_created_at,
			CASE
				WHEN f.is_public THEN CONCAT('https://', f.bucket_name, '.s3.', 'ap-southeast-2', '.amazonaws.com/', f.object_key)
				ELSE 'signed-url-required'
			END as download_url
		FROM lecture_resources lr
		JOIN files f ON lr.file_id = f.id
		WHERE lr.lecture_id = $1 AND f.deleted_at IS NULL
		ORDER BY lr.display_order ASC`

	var resources []model.LectureResource
	err := r.db.SelectContext(ctx, &resources, query, lectureID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.LectureResource{}, nil
		}
		r.logger.Errorf("Failed to get lecture resources: %v", err)
		return nil, err
	}

	return resources, nil
}

func (r *LectureResourceRepository) GetByCourseID(ctx context.Context, courseID uuid.UUID) (map[uuid.UUID][]model.LectureResource, error) {
	r.logger.Infof("DEBUG: GetByCourseID called for course %s", courseID)

	query := `
		SELECT
			lr.id, lr.lecture_id, lr.file_id, lr.resource_type, lr.display_order,
			lr.is_required, lr.created_at, lr.updated_at,
			f.filename, f.original_filename, f.content_type, f.size_bytes,
			f.is_public, f.created_at as file_created_at,
			CASE
				WHEN f.is_public THEN CONCAT('https://', f.bucket_name, '.s3.', 'ap-southeast-2', '.amazonaws.com/', f.object_key)
				ELSE 'signed-url-required'
			END as download_url
		FROM lecture_resources lr
		JOIN files f ON lr.file_id = f.id
		JOIN lectures l ON lr.lecture_id = l.id
		WHERE l.course_id = $1 AND f.deleted_at IS NULL AND l.deleted_at IS NULL
		ORDER BY l.order_number, lr.display_order ASC`

	var resources []model.LectureResource
	err := r.db.SelectContext(ctx, &resources, query, courseID)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Infof("DEBUG: No resources found for course %s", courseID)
			return make(map[uuid.UUID][]model.LectureResource), nil
		}
		r.logger.Errorf("Failed to get course lecture resources: %v", err)
		return nil, err
	}

	r.logger.Infof("DEBUG: Found %d total resources for course %s", len(resources), courseID)

	// Group resources by lecture_id
	resourcesMap := make(map[uuid.UUID][]model.LectureResource)
	for _, resource := range resources {
		resourcesMap[resource.LectureID] = append(resourcesMap[resource.LectureID], resource)
		r.logger.Infof("DEBUG: Resource - LectureID: %s, Filename: %s, Type: %s", resource.LectureID, resource.Filename, resource.ResourceType)
	}

	r.logger.Infof("DEBUG: Grouped resources into %d lecture groups", len(resourcesMap))
	for lectureID, lectureResources := range resourcesMap {
		r.logger.Infof("DEBUG: Lecture %s has %d resources", lectureID, len(lectureResources))
	}

	return resourcesMap, nil
}

func (r *LectureResourceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.LectureResource, error) {
	query := `
		SELECT
			lr.id, lr.lecture_id, lr.file_id, lr.resource_type, lr.display_order,
			lr.is_required, lr.created_at, lr.updated_at,
			f.filename, f.original_filename, f.content_type, f.size_bytes,
			f.is_public, f.created_at as file_created_at, f.bucket_name, f.object_key,
			CASE
				WHEN f.is_public THEN CONCAT('https://', f.bucket_name, '.s3.', 'ap-southeast-2', '.amazonaws.com/', f.object_key)
				ELSE 'signed-url-required'
			END as download_url
		FROM lecture_resources lr
		JOIN files f ON lr.file_id = f.id
		WHERE lr.id = $1 AND f.deleted_at IS NULL`

	var resource model.LectureResource
	err := r.db.GetContext(ctx, &resource, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lecture resource not found")
		}
		r.logger.Errorf("Failed to get lecture resource: %v", err)
		return nil, err
	}

	return &resource, nil
}

func (r *LectureResourceRepository) Update(ctx context.Context, resource *model.LectureResource) error {
	query := `
		UPDATE lecture_resources
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
		r.logger.Errorf("Failed to update lecture resource: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("lecture resource not found")
	}

	return nil
}

func (r *LectureResourceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM lecture_resources WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Errorf("Failed to delete lecture resource: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("lecture resource not found")
	}

	return nil
}

func (r *LectureResourceRepository) DeleteByLectureID(ctx context.Context, lectureID uuid.UUID) error {
	query := `DELETE FROM lecture_resources WHERE lecture_id = $1`

	_, err := r.db.ExecContext(ctx, query, lectureID)
	if err != nil {
		r.logger.Errorf("Failed to delete lecture resources: %v", err)
		return err
	}

	return nil
}

// BulkCreate creates multiple resources for a lecture
func (r *LectureResourceRepository) BulkCreate(ctx context.Context, lectureID uuid.UUID, resources []model.LectureResource) error {
	if len(resources) == 0 {
		return nil
	}

	query := `
		INSERT INTO lecture_resources (
			lecture_id, file_id, resource_type, display_order, is_required
		) VALUES ($1, $2, $3, $4, $5)`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, resource := range resources {
		resource.LectureID = lectureID // Ensure lecture ID is set
		_, err := tx.ExecContext(
			ctx, query,
			resource.LectureID,
			resource.FileID,
			resource.ResourceType,
			resource.DisplayOrder,
			resource.IsRequired,
		)
		if err != nil {
			r.logger.Errorf("Failed to bulk create lecture resource: %v", err)
			return err
		}
	}

	return tx.Commit()
}

// UpdateDisplayOrder updates the display order of multiple resources
func (r *LectureResourceRepository) UpdateDisplayOrder(ctx context.Context, updates []struct {
	ID           uuid.UUID
	DisplayOrder int32
}) error {
	if len(updates) == 0 {
		return nil
	}

	query := `UPDATE lecture_resources SET display_order = $1, updated_at = now() WHERE id = $2`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, update := range updates {
		_, err := tx.ExecContext(ctx, query, update.DisplayOrder, update.ID)
		if err != nil {
			r.logger.Errorf("Failed to update resource display order: %v", err)
			return err
		}
	}

	return tx.Commit()
}