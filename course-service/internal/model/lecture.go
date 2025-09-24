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