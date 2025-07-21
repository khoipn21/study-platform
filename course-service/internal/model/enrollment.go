package model

import (
	"time"

	"github.com/google/uuid"
)

type EnrollmentStatus string

const (
	EnrollmentStatusEnrolled  EnrollmentStatus = "enrolled"
	EnrollmentStatusCompleted EnrollmentStatus = "completed"
	EnrollmentStatusCancelled EnrollmentStatus = "cancelled"
)

type Enrollment struct {
	ID                 uuid.UUID        `json:"id" db:"id"`
	UserID             uuid.UUID        `json:"user_id" db:"user_id"`
	CourseID           uuid.UUID        `json:"course_id" db:"course_id"`
	Status             EnrollmentStatus `json:"status" db:"status"`
	ProgressPercentage float64          `json:"progress_percentage" db:"progress_percentage"`
	EnrolledAt         time.Time        `json:"enrolled_at" db:"enrolled_at"`
	CompletedAt        *time.Time       `json:"completed_at" db:"completed_at"`
	LastAccessed       *time.Time       `json:"last_accessed" db:"last_accessed"`
}

type EnrollmentFilters struct {
	UserID   string
	CourseID string
	Status   string
	Page     int32
	PageSize int32
}

type EnrollmentSearchResult struct {
	Enrollments []Enrollment `json:"enrollments"`
	TotalCount  int32        `json:"total_count"`
	Page        int32        `json:"page"`
	PageSize    int32        `json:"page_size"`
}