package model

import (
	"time"

	"github.com/google/uuid"
)

type CourseStatus string
type CourseLevel string

const (
	CourseStatusDraft     CourseStatus = "draft"
	CourseStatusPublished CourseStatus = "published"
	CourseStatusArchived  CourseStatus = "archived"
)

const (
	CourseLevelBeginner     CourseLevel = "beginner"
	CourseLevelIntermediate CourseLevel = "intermediate"
	CourseLevelAdvanced     CourseLevel = "advanced"
)

type Course struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	Title           string        `json:"title" db:"title"`
	Description     string        `json:"description" db:"description"`
	InstructorID    uuid.UUID     `json:"instructor_id" db:"instructor_id"`
	InstructorName  string        `json:"instructor_name" db:"instructor_name"`
	Category        string        `json:"category" db:"category"`
	Level           CourseLevel   `json:"level" db:"level"`
	Price           float64       `json:"price" db:"price"`
	Currency        string        `json:"currency" db:"currency"`
	ThumbnailURL    string        `json:"thumbnail_url" db:"thumbnail_url"`
	Status          CourseStatus  `json:"status" db:"status"`
	DurationMinutes int32         `json:"duration_minutes" db:"duration_minutes"`
	EnrollmentCount int32         `json:"enrollment_count" db:"enrollment_count"`
	Rating          float64       `json:"rating" db:"rating"`
	RatingCount     int32         `json:"rating_count" db:"rating_count"`
	Tags            []string      `json:"tags" db:"tags"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
}

type CourseFilters struct {
	Category     string
	Level        string
	Status       string
	InstructorID string
	MinPrice     float64
	MaxPrice     float64
	MinRating    float64
	Query        string
	Page         int32
	PageSize     int32
}

type CourseSearchResult struct {
	Courses    []Course `json:"courses"`
	TotalCount int32    `json:"total_count"`
	Page       int32    `json:"page"`
	PageSize   int32    `json:"page_size"`
}