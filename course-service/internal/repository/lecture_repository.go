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

type LectureRepository struct {
	db *database.DB
}

func NewLectureRepository(db *database.DB) *LectureRepository {
	return &LectureRepository{db: db}
}

func (r *LectureRepository) Create(ctx context.Context, lecture *model.Lecture) error {
	query := `
		INSERT INTO lectures (id, course_id, title, description, type, order_number, duration_minutes, video_url, video_id, status, is_free, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	lecture.ID = uuid.New()
	lecture.CreatedAt = time.Now()
	lecture.UpdatedAt = time.Now()

	// Set default type if not specified
	if lecture.Type == "" {
		lecture.Type = "video"
	}

	_, err := r.db.ExecContext(ctx, query,
		lecture.ID, lecture.CourseID, lecture.Title, lecture.Description, lecture.Type, lecture.OrderNumber,
		lecture.DurationMinutes, lecture.VideoURL, lecture.VideoID, lecture.Status,
		lecture.IsFree, lecture.CreatedAt, lecture.UpdatedAt,
	)

	return err
}

func (r *LectureRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Lecture, error) {
	query := `
		SELECT
			l.id, l.course_id, l.title, l.description, l.type, l.order_number, l.duration_minutes,
			COALESCE(v.stream_url, l.video_url, '') as video_url,
			l.video_id, l.status, l.is_free, l.created_at, l.updated_at
		FROM lectures l
		LEFT JOIN videos v ON NULLIF(l.video_id, '')::uuid = v.id
		WHERE l.id = $1
	`

	lecture := &model.Lecture{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&lecture.ID, &lecture.CourseID, &lecture.Title, &lecture.Description, &lecture.Type, &lecture.OrderNumber,
		&lecture.DurationMinutes, &lecture.VideoURL, &lecture.VideoID, &lecture.Status,
		&lecture.IsFree, &lecture.CreatedAt, &lecture.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lecture not found")
		}
		return nil, err
	}

	return lecture, nil
}

func (r *LectureRepository) Update(ctx context.Context, lecture *model.Lecture) error {
	query := `
		UPDATE lectures
		SET title = $2, description = $3, order_number = $4, duration_minutes = $5, video_url = $6, video_id = $7, status = $8, is_free = $9, updated_at = $10
		WHERE id = $1
	`
	
	lecture.UpdatedAt = time.Now()
	
	result, err := r.db.ExecContext(ctx, query,
		lecture.ID, lecture.Title, lecture.Description, lecture.OrderNumber,
		lecture.DurationMinutes, lecture.VideoURL, lecture.VideoID, lecture.Status,
		lecture.IsFree, lecture.UpdatedAt,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("lecture not found")
	}
	
	return nil
}

func (r *LectureRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM lectures WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("lecture not found")
	}
	
	return nil
}

func (r *LectureRepository) List(ctx context.Context, filters model.LectureFilters) (*model.LectureSearchResult, error) {
	var conditions []string
	var args []interface{}
	argCount := 0

	// Modified query to LEFT JOIN with videos table to get stream_url
	baseQuery := `
		SELECT
			l.id, l.course_id, l.title, l.description, l.type, l.order_number, l.duration_minutes,
			COALESCE(v.stream_url, l.video_url, '') as video_url,
			l.video_id, l.status, l.is_free, l.created_at, l.updated_at
		FROM lectures l
		LEFT JOIN videos v ON NULLIF(l.video_id, '')::uuid = v.id
	`

	if filters.CourseID != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("l.course_id = $%d::uuid", argCount))
		args = append(args, filters.CourseID)
	}

	if filters.Status != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("l.status = $%d", argCount))
		args = append(args, filters.Status)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total results - also need to update count query to match the JOIN
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM lectures l LEFT JOIN videos v ON NULLIF(l.video_id, '')::uuid = v.id %s", whereClause)
	var totalCount int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	// Get paginated results
	query := fmt.Sprintf("%s %s ORDER BY l.order_number ASC LIMIT $%d OFFSET $%d",
		baseQuery, whereClause, argCount+1, argCount+2)

	limit := filters.PageSize
	if limit <= 0 {
		limit = 50
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

	var lectures []model.Lecture
	for rows.Next() {
		var lecture model.Lecture
		err := rows.Scan(
			&lecture.ID, &lecture.CourseID, &lecture.Title, &lecture.Description, &lecture.Type, &lecture.OrderNumber,
			&lecture.DurationMinutes, &lecture.VideoURL, &lecture.VideoID, &lecture.Status,
			&lecture.IsFree, &lecture.CreatedAt, &lecture.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		lectures = append(lectures, lecture)
	}

	return &model.LectureSearchResult{
		Lectures:   lectures,
		TotalCount: totalCount,
		Page:       filters.Page,
		PageSize:   limit,
	}, nil
}

func (r *LectureRepository) GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]model.Lecture, error) {
	query := `
		SELECT
			l.id, l.course_id, l.title, l.description, l.type, l.order_number, l.duration_minutes,
			COALESCE(v.stream_url, l.video_url, '') as video_url,
			l.video_id, l.status, l.is_free, l.created_at, l.updated_at
		FROM lectures l
		LEFT JOIN videos v ON NULLIF(l.video_id, '')::uuid = v.id
		WHERE l.course_id = $1
		ORDER BY l.order_number ASC
	`

	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lectures []model.Lecture
	for rows.Next() {
		var lecture model.Lecture
		err := rows.Scan(
			&lecture.ID, &lecture.CourseID, &lecture.Title, &lecture.Description, &lecture.Type, &lecture.OrderNumber,
			&lecture.DurationMinutes, &lecture.VideoURL, &lecture.VideoID, &lecture.Status,
			&lecture.IsFree, &lecture.CreatedAt, &lecture.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		lectures = append(lectures, lecture)
	}

	return lectures, nil
}

func (r *LectureRepository) UpdateCourseDuration(ctx context.Context, courseID uuid.UUID) error {
	query := `
		UPDATE courses
		SET duration_minutes = (
			SELECT COALESCE(SUM(duration_minutes), 0) FROM lectures
			WHERE course_id = $1 AND status = 'published'
		)
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, courseID)
	return err
}