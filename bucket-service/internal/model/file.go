package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Filename         string    `gorm:"size:255;not null" json:"filename"`
	OriginalFilename string    `gorm:"size:255;not null" json:"original_filename"`
	ContentType      string    `gorm:"size:100;not null" json:"content_type"`
	SizeBytes        int64     `gorm:"not null" json:"size_bytes"`
	BucketName       string    `gorm:"size:100;not null" json:"bucket_name"`
	ObjectKey        string    `gorm:"size:500;not null" json:"object_key"`
	UploadUserID     uuid.UUID `gorm:"type:uuid;not null" json:"upload_user_id"`
	IsPublic         bool      `gorm:"default:false" json:"is_public"`
	Metadata         *string   `gorm:"type:jsonb" json:"metadata,omitempty"`
	Checksum         *string   `gorm:"size:64" json:"checksum,omitempty"`
	ThumbnailURL     *string   `gorm:"size:500" json:"thumbnail_url,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type FilePermission struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FileID         uuid.UUID `gorm:"type:uuid;not null;index:idx_file_permissions_file_user" json:"file_id"`
	UserID         *uuid.UUID `gorm:"type:uuid;index:idx_file_permissions_file_user" json:"user_id,omitempty"`
	PermissionType string    `gorm:"size:20;not null" json:"permission_type"` // 'read', 'write', 'delete'
	GrantedBy      uuid.UUID `gorm:"type:uuid;not null" json:"granted_by"`
	CreatedAt      time.Time `json:"created_at"`
	
	File File `gorm:"foreignKey:FileID;constraint:OnDelete:CASCADE" json:"-"`
}

type UploadSession struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UploadID     string    `gorm:"size:255;not null" json:"upload_id"` // S3 multipart upload ID
	Filename     string    `gorm:"size:255;not null" json:"filename"`
	ContentType  string    `gorm:"size:100;not null" json:"content_type"`
	TotalSize    int64     `gorm:"not null" json:"total_size"`
	UploadedSize int64     `gorm:"default:0" json:"uploaded_size"`
	BucketName   string    `gorm:"size:100;not null" json:"bucket_name"`
	ObjectKey    string    `gorm:"size:500;not null" json:"object_key"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index:idx_upload_sessions_user_status" json:"user_id"`
	Status       string    `gorm:"size:20;default:'active';index:idx_upload_sessions_user_status" json:"status"` // 'active', 'completed', 'aborted'
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
}

// FileType represents supported file types
type FileType string

const (
	FileTypeDocument FileType = "document"
	FileTypeImage    FileType = "image"
	FileTypeArchive  FileType = "archive"
	FileTypeOther    FileType = "other"
)

// GetFileType determines the file type based on content type
func GetFileType(contentType string) FileType {
	documentTypes := map[string]bool{
		"application/pdf":  true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/vnd.ms-powerpoint": true,
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
		"application/vnd.ms-excel": true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
		"text/plain": true,
		"application/rtf": true,
	}

	imageTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
		"image/svg+xml": true,
		"image/bmp": true,
	}

	archiveTypes := map[string]bool{
		"application/zip": true,
		"application/x-rar-compressed": true,
		"application/x-7z-compressed": true,
		"application/x-tar": true,
		"application/gzip": true,
	}

	if documentTypes[contentType] {
		return FileTypeDocument
	}
	if imageTypes[contentType] {
		return FileTypeImage
	}
	if archiveTypes[contentType] {
		return FileTypeArchive
	}
	return FileTypeOther
}

// GetBucketForFileType returns the appropriate bucket name for a file type
func GetBucketForFileType(fileType FileType, baseBucketName string) string {
	switch fileType {
	case FileTypeDocument:
		return baseBucketName + "-documents"
	case FileTypeImage:
		return baseBucketName + "-images"
	case FileTypeArchive:
		return baseBucketName + "-course-materials"
	default:
		return baseBucketName + "-files"
	}
}