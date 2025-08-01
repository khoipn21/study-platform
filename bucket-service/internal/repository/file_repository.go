package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"bucket-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

// File operations
func (r *FileRepository) CreateFile(ctx context.Context, file *model.File) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *FileRepository) GetFileByID(ctx context.Context, id uuid.UUID) (*model.File, error) {
	var file model.File
	err := r.db.WithContext(ctx).First(&file, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) GetFileByObjectKey(ctx context.Context, bucketName, objectKey string) (*model.File, error) {
	var file model.File
	err := r.db.WithContext(ctx).First(&file, "bucket_name = ? AND object_key = ?", bucketName, objectKey).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) UpdateFile(ctx context.Context, file *model.File) error {
	file.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(file).Error
}

func (r *FileRepository) DeleteFile(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.File{}, "id = ?", id).Error
}

func (r *FileRepository) SoftDeleteFile(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.File{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

type ListFilesOptions struct {
	UserID     *uuid.UUID
	BucketName string
	ContentType string
	IsPublic   *bool
	Page       int
	Limit      int
	SortBy     string
	SortOrder  string
	Search     string
}

type ListFilesResult struct {
	Files      []*model.File
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}

func (r *FileRepository) ListFiles(ctx context.Context, opts *ListFilesOptions) (*ListFilesResult, error) {
	query := r.db.WithContext(ctx).Model(&model.File{})

	// Apply filters
	if opts.UserID != nil {
		query = query.Where("upload_user_id = ?", *opts.UserID)
	}
	if opts.BucketName != "" {
		query = query.Where("bucket_name = ?", opts.BucketName)
	}
	if opts.ContentType != "" {
		query = query.Where("content_type LIKE ?", opts.ContentType+"%")
	}
	if opts.IsPublic != nil {
		query = query.Where("is_public = ?", *opts.IsPublic)
	}
	if opts.Search != "" {
		searchTerm := "%" + strings.ToLower(opts.Search) + "%"
		query = query.Where("LOWER(filename) LIKE ? OR LOWER(original_filename) LIKE ?", 
			searchTerm, searchTerm)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	sortBy := "created_at"
	if opts.SortBy != "" {
		sortBy = opts.SortBy
	}
	sortOrder := "DESC"
	if opts.SortOrder != "" {
		sortOrder = strings.ToUpper(opts.SortOrder)
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	// Apply pagination
	if opts.Limit <= 0 {
		opts.Limit = 20
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}
	offset := (opts.Page - 1) * opts.Limit
	query = query.Offset(offset).Limit(opts.Limit)

	// Execute query
	var files []*model.File
	if err := query.Find(&files).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(opts.Limit) - 1) / int64(opts.Limit))

	return &ListFilesResult{
		Files:      files,
		Total:      total,
		Page:       opts.Page,
		Limit:      opts.Limit,
		TotalPages: totalPages,
	}, nil
}

// File permissions
func (r *FileRepository) CreateFilePermission(ctx context.Context, permission *model.FilePermission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

func (r *FileRepository) GetFilePermissions(ctx context.Context, fileID uuid.UUID) ([]*model.FilePermission, error) {
	var permissions []*model.FilePermission
	err := r.db.WithContext(ctx).Where("file_id = ?", fileID).Find(&permissions).Error
	return permissions, err
}

func (r *FileRepository) HasFilePermission(ctx context.Context, fileID, userID uuid.UUID, permissionType string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.FilePermission{}).
		Where("file_id = ? AND user_id = ? AND permission_type = ?", fileID, userID, permissionType).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *FileRepository) DeleteFilePermission(ctx context.Context, fileID, userID uuid.UUID, permissionType string) error {
	return r.db.WithContext(ctx).
		Where("file_id = ? AND user_id = ? AND permission_type = ?", fileID, userID, permissionType).
		Delete(&model.FilePermission{}).Error
}

// Upload sessions
func (r *FileRepository) CreateUploadSession(ctx context.Context, session *model.UploadSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *FileRepository) GetUploadSession(ctx context.Context, sessionID uuid.UUID) (*model.UploadSession, error) {
	var session model.UploadSession
	err := r.db.WithContext(ctx).First(&session, "id = ?", sessionID).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *FileRepository) GetUploadSessionByUploadID(ctx context.Context, uploadID string) (*model.UploadSession, error) {
	var session model.UploadSession
	err := r.db.WithContext(ctx).First(&session, "upload_id = ?", uploadID).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *FileRepository) UpdateUploadSession(ctx context.Context, session *model.UploadSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *FileRepository) DeleteUploadSession(ctx context.Context, sessionID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.UploadSession{}, "id = ?", sessionID).Error
}

func (r *FileRepository) GetExpiredUploadSessions(ctx context.Context) ([]*model.UploadSession, error) {
	var sessions []*model.UploadSession
	err := r.db.WithContext(ctx).
		Where("expires_at < ? AND status = 'active'", time.Now()).
		Find(&sessions).Error
	return sessions, err
}

func (r *FileRepository) GetUserUploadSessions(ctx context.Context, userID uuid.UUID, status string) ([]*model.UploadSession, error) {
	var sessions []*model.UploadSession
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&sessions).Error
	return sessions, err
}

// Statistics
func (r *FileRepository) GetUserFileStats(ctx context.Context, userID uuid.UUID) (*FileStats, error) {
	var stats FileStats
	
	// Total files
	err := r.db.WithContext(ctx).Model(&model.File{}).
		Where("upload_user_id = ?", userID).
		Count(&stats.TotalFiles).Error
	if err != nil {
		return nil, err
	}

	// Total size
	err = r.db.WithContext(ctx).Model(&model.File{}).
		Where("upload_user_id = ?", userID).
		Select("COALESCE(SUM(size_bytes), 0)").
		Scan(&stats.TotalSize).Error
	if err != nil {
		return nil, err
	}

	// Files by type
	var typeStats []struct {
		ContentType string
		Count       int64
		TotalSize   int64
	}
	err = r.db.WithContext(ctx).Model(&model.File{}).
		Where("upload_user_id = ?", userID).
		Select("content_type, COUNT(*) as count, COALESCE(SUM(size_bytes), 0) as total_size").
		Group("content_type").
		Scan(&typeStats).Error
	if err != nil {
		return nil, err
	}

	stats.FilesByType = make(map[string]FileTypeStats)
	for _, stat := range typeStats {
		stats.FilesByType[stat.ContentType] = FileTypeStats{
			Count:     stat.Count,
			TotalSize: stat.TotalSize,
		}
	}

	return &stats, nil
}

type FileStats struct {
	TotalFiles  int64                     `json:"total_files"`
	TotalSize   int64                     `json:"total_size"`
	FilesByType map[string]FileTypeStats  `json:"files_by_type"`
}

type FileTypeStats struct {
	Count     int64 `json:"count"`
	TotalSize int64 `json:"total_size"`
}