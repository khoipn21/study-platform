package model

import (
	"time"

	"github.com/google/uuid"
)

type LectureStatus string

const (
	LectureStatusDraft     LectureStatus = "draft"
	LectureStatusPublished LectureStatus = "published"
)

type Lecture struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	CourseID        uuid.UUID     `json:"course_id" db:"course_id"`
	Title           string        `json:"title" db:"title"`
	Description     string        `json:"description" db:"description"`
	Type            string        `json:"type" db:"type"`
	OrderNumber     int32         `json:"order_number" db:"order_number"`
	DurationMinutes int32         `json:"duration_minutes" db:"duration_minutes"`
	VideoURL        string        `json:"video_url" db:"video_url"`
	VideoID         string        `json:"video_id" db:"video_id"`
	Status          LectureStatus `json:"status" db:"status"`
	IsFree          bool          `json:"is_free" db:"is_free"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time    `json:"deleted_at,omitempty" db:"deleted_at"`

	// Related data (not stored in DB, populated on fetch)
	Resources       []LectureResource `json:"resources,omitempty"`
}

type LectureFilters struct {
	CourseID string
	Status   string
	Page     int32
	PageSize int32
}

type LectureSearchResult struct {
	Lectures   []Lecture `json:"lectures"`
	TotalCount int32     `json:"total_count"`
	Page       int32     `json:"page"`
	PageSize   int32     `json:"page_size"`
}

// LectureResource represents a file resource attached to a lecture
type LectureResource struct {
	ID           uuid.UUID `json:"id" db:"id"`
	LectureID    uuid.UUID `json:"lecture_id" db:"lecture_id"`
	FileID       uuid.UUID `json:"file_id" db:"file_id"`
	ResourceType string    `json:"resource_type" db:"resource_type"`
	DisplayOrder int32     `json:"display_order" db:"display_order"`
	IsRequired   bool      `json:"is_required" db:"is_required"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`

	// File information (populated from files table join)
	Filename     string    `json:"filename" db:"filename"`
	OriginalName string    `json:"original_name" db:"original_filename"`
	FileType     string    `json:"file_type" db:"content_type"`
	FileSize     int64     `json:"file_size" db:"size_bytes"`
	DownloadURL  string    `json:"download_url" db:"download_url"`
	IsPublic     bool      `json:"is_public" db:"is_public"`
	UploadedAt   time.Time `json:"uploaded_at" db:"file_created_at"`
	BucketName   string    `json:"bucket_name" db:"bucket_name"`
	ObjectKey    string    `json:"object_key" db:"object_key"`
}