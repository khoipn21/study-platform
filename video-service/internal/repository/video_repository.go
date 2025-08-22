package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"video-service/internal/model"
)

type VideoRepository struct {
	db *sql.DB
}

func NewVideoRepository(db *sql.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

// CreateVideo creates a new video record
func (vr *VideoRepository) CreateVideo(video *model.Video) error {
	query := `
		INSERT INTO videos (
			id, cloudflare_uid, title, description, upload_user_id, 
			course_id, lecture_id, status, visibility, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at
	`

	if video.ID == uuid.Nil {
		video.ID = uuid.New()
	}

	var metadataBytes []byte
	if video.Metadata != nil {
		value, err := video.Metadata.Value()
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		if value != nil {
			metadataBytes = value.([]byte)
		}
	}

	err := vr.db.QueryRow(
		query, video.ID, video.CloudflareUID, video.Title, video.Description,
		video.UploadUserID, video.CourseID, video.LectureID, video.Status,
		video.Visibility, metadataBytes,
	).Scan(&video.CreatedAt, &video.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create video: %w", err)
	}

	return nil
}

// GetVideoByID retrieves a video by its ID
func (vr *VideoRepository) GetVideoByID(id uuid.UUID) (*model.Video, error) {
	query := `
		SELECT id, cloudflare_uid, title, description, duration_seconds, 
			   file_size_bytes, upload_user_id, course_id, lecture_id, 
			   status, visibility, thumbnail_url, stream_url, preview_url,
			   metadata, created_at, updated_at, deleted_at
		FROM videos 
		WHERE id = $1 AND deleted_at IS NULL
	`

	video := &model.Video{}
	var metadataBytes []byte

	var thumbnailURL, streamURL, previewURL sql.NullString
	
	err := vr.db.QueryRow(query, id).Scan(
		&video.ID, &video.CloudflareUID, &video.Title, &video.Description,
		&video.DurationSeconds, &video.FileSizeBytes, &video.UploadUserID,
		&video.CourseID, &video.LectureID, &video.Status, &video.Visibility,
		&thumbnailURL, &streamURL, &previewURL,
		&metadataBytes, &video.CreatedAt, &video.UpdatedAt, &video.DeletedAt,
	)
	
	// Handle nullable fields
	video.ThumbnailURL = thumbnailURL.String
	video.StreamURL = streamURL.String  
	video.PreviewURL = previewURL.String

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("video not found")
		}
		return nil, fmt.Errorf("failed to get video: %w", err)
	}

	if len(metadataBytes) > 0 {
		if err := video.Metadata.Scan(metadataBytes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return video, nil
}

// GetVideoByCloudflareUID retrieves a video by its Cloudflare UID
func (vr *VideoRepository) GetVideoByCloudflareUID(cloudflareUID string) (*model.Video, error) {
	query := `
		SELECT id, cloudflare_uid, title, description, duration_seconds, 
			   file_size_bytes, upload_user_id, course_id, lecture_id, 
			   status, visibility, thumbnail_url, stream_url, preview_url,
			   metadata, created_at, updated_at, deleted_at
		FROM videos 
		WHERE cloudflare_uid = $1 AND deleted_at IS NULL
	`

	video := &model.Video{}
	var metadataBytes []byte

	err := vr.db.QueryRow(query, cloudflareUID).Scan(
		&video.ID, &video.CloudflareUID, &video.Title, &video.Description,
		&video.DurationSeconds, &video.FileSizeBytes, &video.UploadUserID,
		&video.CourseID, &video.LectureID, &video.Status, &video.Visibility,
		&video.ThumbnailURL, &video.StreamURL, &video.PreviewURL,
		&metadataBytes, &video.CreatedAt, &video.UpdatedAt, &video.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("video not found")
		}
		return nil, fmt.Errorf("failed to get video: %w", err)
	}

	if len(metadataBytes) > 0 {
		if err := video.Metadata.Scan(metadataBytes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return video, nil
}

// UpdateVideo updates a video record
func (vr *VideoRepository) UpdateVideo(video *model.Video) error {
	query := `
		UPDATE videos 
		SET title = $1, description = $2, duration_seconds = $3, 
			file_size_bytes = $4, status = $5, visibility = $6, 
			thumbnail_url = $7, stream_url = $8, preview_url = $9, 
			metadata = $10, updated_at = NOW()
		WHERE id = $11 AND deleted_at IS NULL
	`

	var metadataBytes []byte
	if video.Metadata != nil {
		value, err := video.Metadata.Value()
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		if value != nil {
			metadataBytes = value.([]byte)
		}
	}

	result, err := vr.db.Exec(
		query, video.Title, video.Description, video.DurationSeconds,
		video.FileSizeBytes, video.Status, video.Visibility,
		video.ThumbnailURL, video.StreamURL, video.PreviewURL,
		metadataBytes, video.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update video: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("video not found or already deleted")
	}

	return nil
}

// DeleteVideo soft deletes a video record
func (vr *VideoRepository) DeleteVideo(id uuid.UUID) error {
	query := `
		UPDATE videos 
		SET deleted_at = NOW() 
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := vr.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("video not found or already deleted")
	}

	return nil
}

// ListVideosByUser lists videos uploaded by a specific user
func (vr *VideoRepository) ListVideosByUser(userID uuid.UUID, limit, offset int) ([]*model.Video, error) {
	query := `
		SELECT id, cloudflare_uid, title, description, duration_seconds, 
			   file_size_bytes, upload_user_id, course_id, lecture_id, 
			   status, visibility, thumbnail_url, stream_url, preview_url,
			   metadata, created_at, updated_at
		FROM videos 
		WHERE upload_user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := vr.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query videos: %w", err)
	}
	defer rows.Close()

	var videos []*model.Video
	for rows.Next() {
		video := &model.Video{}
		var metadataBytes []byte

		err := rows.Scan(
			&video.ID, &video.CloudflareUID, &video.Title, &video.Description,
			&video.DurationSeconds, &video.FileSizeBytes, &video.UploadUserID,
			&video.CourseID, &video.LectureID, &video.Status, &video.Visibility,
			&video.ThumbnailURL, &video.StreamURL, &video.PreviewURL,
			&metadataBytes, &video.CreatedAt, &video.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}

		if len(metadataBytes) > 0 {
			if err := video.Metadata.Scan(metadataBytes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return videos, nil
}

// ListVideosByCourse lists videos for a specific course
func (vr *VideoRepository) ListVideosByCourse(courseID uuid.UUID, limit, offset int) ([]*model.Video, error) {
	query := `
		SELECT id, cloudflare_uid, title, description, duration_seconds, 
			   file_size_bytes, upload_user_id, course_id, lecture_id, 
			   status, visibility, thumbnail_url, stream_url, preview_url,
			   metadata, created_at, updated_at
		FROM videos 
		WHERE course_id = $1 AND deleted_at IS NULL AND status = 'ready'
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := vr.db.Query(query, courseID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query videos: %w", err)
	}
	defer rows.Close()

	var videos []*model.Video
	for rows.Next() {
		video := &model.Video{}
		var metadataBytes []byte

		err := rows.Scan(
			&video.ID, &video.CloudflareUID, &video.Title, &video.Description,
			&video.DurationSeconds, &video.FileSizeBytes, &video.UploadUserID,
			&video.CourseID, &video.LectureID, &video.Status, &video.Visibility,
			&video.ThumbnailURL, &video.StreamURL, &video.PreviewURL,
			&metadataBytes, &video.CreatedAt, &video.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}

		if len(metadataBytes) > 0 {
			if err := video.Metadata.Scan(metadataBytes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return videos, nil
}

// SearchVideos searches videos by title and description
func (vr *VideoRepository) SearchVideos(query string, limit, offset int) ([]*model.Video, error) {
	sqlQuery := `
		SELECT id, cloudflare_uid, title, description, duration_seconds, 
			   file_size_bytes, upload_user_id, course_id, lecture_id, 
			   status, visibility, thumbnail_url, stream_url, preview_url,
			   metadata, created_at, updated_at
		FROM videos 
		WHERE (title ILIKE $1 OR description ILIKE $1) 
		  AND deleted_at IS NULL AND status = 'ready' AND visibility IN ('public', 'unlisted')
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	searchTerm := "%" + strings.ToLower(query) + "%"

	rows, err := vr.db.Query(sqlQuery, searchTerm, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search videos: %w", err)
	}
	defer rows.Close()

	var videos []*model.Video
	for rows.Next() {
		video := &model.Video{}
		var metadataBytes []byte

		err := rows.Scan(
			&video.ID, &video.CloudflareUID, &video.Title, &video.Description,
			&video.DurationSeconds, &video.FileSizeBytes, &video.UploadUserID,
			&video.CourseID, &video.LectureID, &video.Status, &video.Visibility,
			&video.ThumbnailURL, &video.StreamURL, &video.PreviewURL,
			&metadataBytes, &video.CreatedAt, &video.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}

		if len(metadataBytes) > 0 {
			if err := video.Metadata.Scan(metadataBytes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return videos, nil
}

// GetVideoQualitiesByVideoID retrieves all quality variants for a video
func (vr *VideoRepository) GetVideoQualitiesByVideoID(videoID uuid.UUID) ([]*model.VideoQuality, error) {
	query := `
		SELECT id, video_id, quality_label, bitrate_kbps, width, height, 
			   fps, codec, url, file_size_bytes, created_at
		FROM video_qualities 
		WHERE video_id = $1
		ORDER BY bitrate_kbps ASC
	`

	rows, err := vr.db.Query(query, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to query video qualities: %w", err)
	}
	defer rows.Close()

	var qualities []*model.VideoQuality
	for rows.Next() {
		quality := &model.VideoQuality{}

		err := rows.Scan(
			&quality.ID, &quality.VideoID, &quality.QualityLabel,
			&quality.BitrateKbps, &quality.Width, &quality.Height,
			&quality.FPS, &quality.Codec, &quality.URL,
			&quality.FileSizeBytes, &quality.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan video quality: %w", err)
		}

		qualities = append(qualities, quality)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return qualities, nil
}

// CreateVideoQuality creates a new video quality record
func (vr *VideoRepository) CreateVideoQuality(quality *model.VideoQuality) error {
	query := `
		INSERT INTO video_qualities (
			id, video_id, quality_label, bitrate_kbps, width, height, 
			fps, codec, url, file_size_bytes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at
	`

	if quality.ID == uuid.Nil {
		quality.ID = uuid.New()
	}

	err := vr.db.QueryRow(
		query, quality.ID, quality.VideoID, quality.QualityLabel,
		quality.BitrateKbps, quality.Width, quality.Height,
		quality.FPS, quality.Codec, quality.URL, quality.FileSizeBytes,
	).Scan(&quality.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create video quality: %w", err)
	}

	return nil
}

// UpdateVideoStatus updates the status of a video
func (vr *VideoRepository) UpdateVideoStatus(id uuid.UUID, status string) error {
	query := `
		UPDATE videos 
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := vr.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update video status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("video not found or already deleted")
	}

	return nil
}

// UpdateVideoURLs updates streaming URLs for a video
func (vr *VideoRepository) UpdateVideoURLs(id uuid.UUID, thumbnailURL, streamURL, previewURL string) error {
	query := `
		UPDATE videos 
		SET thumbnail_url = $1, stream_url = $2, preview_url = $3, updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := vr.db.Exec(query, thumbnailURL, streamURL, previewURL, id)
	if err != nil {
		return fmt.Errorf("failed to update video URLs: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("video not found or already deleted")
	}

	return nil
}

// GetVideosByStatus retrieves videos with a specific status
func (vr *VideoRepository) GetVideosByStatus(status string, limit, offset int) ([]*model.Video, error) {
	query := `
		SELECT id, cloudflare_uid, title, description, duration_seconds, 
			   file_size_bytes, upload_user_id, course_id, lecture_id, 
			   status, visibility, thumbnail_url, stream_url, preview_url,
			   metadata, created_at, updated_at
		FROM videos 
		WHERE status = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := vr.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query videos by status: %w", err)
	}
	defer rows.Close()

	var videos []*model.Video
	for rows.Next() {
		video := &model.Video{}
		var metadataBytes []byte

		err := rows.Scan(
			&video.ID, &video.CloudflareUID, &video.Title, &video.Description,
			&video.DurationSeconds, &video.FileSizeBytes, &video.UploadUserID,
			&video.CourseID, &video.LectureID, &video.Status, &video.Visibility,
			&video.ThumbnailURL, &video.StreamURL, &video.PreviewURL,
			&metadataBytes, &video.CreatedAt, &video.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}

		if len(metadataBytes) > 0 {
			if err := video.Metadata.Scan(metadataBytes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return videos, nil
}