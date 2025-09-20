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
	ID                     uuid.UUID     `json:"id" db:"id"`
	Title                  string        `json:"title" db:"title"`
	Description            string        `json:"description" db:"description"`
	InstructorID           uuid.UUID     `json:"instructor_id" db:"instructor_id"`
	InstructorName         string        `json:"instructor_name" db:"instructor_name"`
	Category               string        `json:"category" db:"category"`
	Level                  CourseLevel   `json:"level" db:"level"`
	Price                  float64       `json:"price" db:"price"`
	Currency               string        `json:"currency" db:"currency"`
	ThumbnailURL           string        `json:"thumbnail_url" db:"thumbnail_url"`
	Status                 CourseStatus  `json:"status" db:"status"`
	DurationMinutes        int32         `json:"duration_minutes" db:"duration_minutes"`
	EnrollmentCount        int32         `json:"enrollment_count" db:"enrollment_count"`
	Rating                 float64       `json:"rating" db:"rating"`
	RatingCount            int32         `json:"rating_count" db:"rating_count"`
	Tags                   []string      `json:"tags" db:"tags"`
	IsPaid                 bool          `json:"is_paid" db:"is_paid"`
	LemonSqueezyProductID  *string       `json:"lemon_squeezy_product_id,omitempty" db:"lemon_squeezy_product_id"`
	LemonSqueezyVariantID  *string       `json:"lemon_squeezy_variant_id,omitempty" db:"lemon_squeezy_variant_id"`
	CreatedAt              time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time     `json:"updated_at" db:"updated_at"`
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

// Lemon Squeezy specific models
type LemonSqueezyProduct struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	LemonSqueezyProductID string    `json:"lemon_squeezy_product_id" db:"lemon_squeezy_product_id"`
	Name                 string    `json:"name" db:"name"`
	Description          string    `json:"description" db:"description"`
	Status               string    `json:"status" db:"status"`
	StoreID              string    `json:"store_id" db:"store_id"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

type LemonSqueezyVariant struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	LemonSqueezyVariantID string    `json:"lemon_squeezy_variant_id" db:"lemon_squeezy_variant_id"`
	LemonSqueezyProductID string    `json:"lemon_squeezy_product_id" db:"lemon_squeezy_product_id"`
	Name                  string    `json:"name" db:"name"`
	Price                 float64   `json:"price" db:"price"`
	Currency              string    `json:"currency" db:"currency"`
	Status                string    `json:"status" db:"status"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

// Request models for creating/updating courses with Lemon Squeezy integration
type CreateCourseRequest struct {
	Title                 string   `json:"title" binding:"required"`
	Description           string   `json:"description"`
	Category              string   `json:"category" binding:"required"`
	Level                 string   `json:"level" binding:"required"`
	Price                 float64  `json:"price"`
	Currency              string   `json:"currency"`
	ThumbnailURL          string   `json:"thumbnail_url"`
	Tags                  []string `json:"tags"`
	IsPaid                bool     `json:"is_paid"`
	LemonSqueezyProductID *string  `json:"lemon_squeezy_product_id,omitempty"`
	LemonSqueezyVariantID *string  `json:"lemon_squeezy_variant_id,omitempty"`
}

type UpdateCourseRequest struct {
	Title                 *string  `json:"title,omitempty"`
	Description           *string  `json:"description,omitempty"`
	Category              *string  `json:"category,omitempty"`
	Level                 *string  `json:"level,omitempty"`
	Price                 *float64 `json:"price,omitempty"`
	Currency              *string  `json:"currency,omitempty"`
	ThumbnailURL          *string  `json:"thumbnail_url,omitempty"`
	Status                *string  `json:"status,omitempty"`
	Tags                  []string `json:"tags,omitempty"`
	IsPaid                *bool    `json:"is_paid,omitempty"`
	LemonSqueezyProductID *string  `json:"lemon_squeezy_product_id,omitempty"`
	LemonSqueezyVariantID *string  `json:"lemon_squeezy_variant_id,omitempty"`
}

// Access Control Types
type AccessLevel string
type CourseType string
type LectureType string

const (
	AccessLevelFull    AccessLevel = "full"
	AccessLevelPreview AccessLevel = "preview"
	AccessLevelDenied  AccessLevel = "denied"
)

const (
	CourseTypeFree CourseType = "free"
	CourseTypePaid CourseType = "paid"
)

const (
	LectureTypeFree    LectureType = "free"
	LectureTypePaid    LectureType = "paid"
	LectureTypePreview LectureType = "preview"
)

// CourseAccessResult contains the result of course access validation
type CourseAccessResult struct {
	HasAccess   bool        `json:"has_access"`
	AccessLevel AccessLevel `json:"access_level"`
	CourseType  CourseType  `json:"course_type"`
	CoursePrice float64     `json:"course_price,omitempty"`
	Currency    string      `json:"currency,omitempty"`
	Message     string      `json:"message"`
}

// LectureAccessResult contains the result of lecture access validation
type LectureAccessResult struct {
	HasAccess        bool                `json:"has_access"`
	AccessLevel      AccessLevel         `json:"access_level"`
	LectureType      LectureType         `json:"lecture_type"`
	CourseAccess     *CourseAccessResult `json:"course_access"`
	PreviewTimeLimit int                 `json:"preview_time_limit,omitempty"` // in seconds
	Message          string              `json:"message"`
}