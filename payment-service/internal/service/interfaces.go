package service

import (
	"context"
	"time"

	"github.com/study-platform/payment-service/internal/model"
)

// ProgressServiceClient interface for communicating with progress service
type ProgressServiceClient interface {
	UpdateEnrollment(ctx context.Context, enrollment *model.Enrollment) error
	GetUserProgress(ctx context.Context, userID, courseID string) (*UserProgress, error)
}

// CourseServiceClient interface for communicating with course service
type CourseServiceClient interface {
	GetCourse(ctx context.Context, courseID string) (*model.Course, error)
	ValidateCourseAccess(ctx context.Context, userID, courseID string) (*CourseAccessValidation, error)
}

// Removed duplicate type declarations - they exist in access_validator.go

// CheckoutRiskAssessment contains data for fraud detection
type CheckoutRiskAssessment struct {
	UserID   string
	CourseID string
	ClientIP string
	Amount   float64
	Currency string
}

// UserProgress represents user's progress in a course
type UserProgress struct {
	UserID             string    `json:"user_id"`
	CourseID           string    `json:"course_id"`
	ProgressPercentage float64   `json:"progress_percentage"`
	LastAccessed       time.Time `json:"last_accessed"`
}

// WebhookRetryInfo contains information for webhook retry
type WebhookRetryInfo struct {
	EventID     string
	Payload     interface{}
	Signature   string
	RetryCount  int
	NextRetryAt time.Time
}