package service

import (
	"context"
	"time"

	"github.com/study-platform/payment-service/internal/repository"
	"github.com/study-platform/pkg/logger"
)

// AccessValidator provides basic course access validation
type AccessValidator struct {
	courseRepo      *repository.CourseRepository
	enrollmentRepo  *repository.EnrollmentRepository
	transactionRepo *repository.TransactionRepository
	previewRepo     *repository.PreviewSessionRepository
	logger          logger.Logger
}

// CourseAccessValidation represents the result of course access validation
type CourseAccessValidation struct {
	UserID              string     `json:"user_id"`
	CourseID            string     `json:"course_id"`
	HasAccess           bool       `json:"has_access"`
	AccessLevel         string     `json:"access_level"` // full, preview, denied
	PaymentRequired     bool       `json:"payment_required"`
	PaymentVerified     bool       `json:"payment_verified"`
	TransactionID       *string    `json:"transaction_id,omitempty"`
	CoursePrice         float64    `json:"course_price,omitempty"`
	Currency            string     `json:"currency,omitempty"`
	CheckoutURL         string     `json:"checkout_url,omitempty"`
	RiskScore           float64    `json:"risk_score,omitempty"`
	AccessExpiresAt     *time.Time `json:"access_expires_at,omitempty"`
	IsPreview           bool       `json:"is_preview"`
	PreviewTimeRemaining int       `json:"preview_time_remaining"`
	Message             string     `json:"message"`
	CachedUntil         time.Time  `json:"cached_until"`
}

// AccessAuditInfo contains information for audit logging
type AccessAuditInfo struct {
	UserID       string `json:"user_id"`
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"` // course, lecture
	ClientIP     string `json:"client_ip,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
}

// NewAccessValidator creates a new access validator instance
func NewAccessValidator(
	courseRepo *repository.CourseRepository,
	enrollmentRepo *repository.EnrollmentRepository,
	transactionRepo *repository.TransactionRepository,
	previewRepo *repository.PreviewSessionRepository,
	logger logger.Logger,
) *AccessValidator {
	return &AccessValidator{
		courseRepo:      courseRepo,
		enrollmentRepo:  enrollmentRepo,
		transactionRepo: transactionRepo,
		previewRepo:     previewRepo,
		logger:          logger,
	}
}

// ValidateCourseAccess performs basic course access validation
func (v *AccessValidator) ValidateCourseAccess(ctx context.Context, userID, courseID string, auditInfo *AccessAuditInfo) (*CourseAccessValidation, error) {
	v.logger.Infof("Validating course access for user %s, course %s", userID, courseID)

	// Get course details
	course, err := v.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		v.logger.Errorf("Course not found: %s, error: %v", courseID, err)
		return nil, err
	}

	result := &CourseAccessValidation{
		UserID:          userID,
		CourseID:        courseID,
		CoursePrice:     course.Price,
		Currency:        course.Currency,
		PaymentRequired: course.IsPaid,
		CachedUntil:     time.Now().Add(5 * time.Minute),
	}

	// If course is free, grant full access
	if !course.IsPaid {
		result.HasAccess = true
		result.AccessLevel = "full"
		result.PaymentVerified = true
		result.Message = "Free course - full access granted"
		return result, nil
	}

	// For paid courses, check enrollment and payment status
	enrollment, err := v.enrollmentRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		// No enrollment - access denied
		result.HasAccess = false
		result.AccessLevel = "denied"
		result.Message = "Purchase required to access this course"
		return result, nil
	}

	// Check enrollment payment status
	if enrollment.PaymentStatus == "paid" && enrollment.PaymentVerifiedAt != nil {
		result.HasAccess = true
		result.AccessLevel = "full"
		result.PaymentVerified = true
		result.Message = "Course purchased - full access granted"
		result.TransactionID = enrollment.TransactionID
		return result, nil
	}

	// Enrollment exists but payment not verified
	result.HasAccess = false
	result.AccessLevel = "denied"
	result.PaymentVerified = false
	result.Message = "Payment verification pending"
	result.TransactionID = enrollment.TransactionID

	return result, nil
}

// ValidateLectureAccess validates user access to a specific lecture
func (v *AccessValidator) ValidateLectureAccess(ctx context.Context, userID, courseID, lectureID string, auditInfo *AccessAuditInfo) (*CourseAccessValidation, error) {
	v.logger.Infof("Validating lecture access for user %s, course %s, lecture %s", userID, courseID, lectureID)

	// For lectures, we first validate course access, then check lecture-specific permissions
	courseAccess, err := v.ValidateCourseAccess(ctx, userID, courseID, auditInfo)
	if err != nil {
		return nil, err
	}

	// If user doesn't have course access, they can't access lectures
	if !courseAccess.HasAccess {
		return courseAccess, nil
	}

	// Additional lecture-specific validation could go here
	// For now, if they have course access, they have lecture access
	return courseAccess, nil
}

// ClearUserCache clears cached access for a specific user and course
func (v *AccessValidator) ClearUserCache(ctx context.Context, userID, courseID string) error {
	// For now, just log the cache clear operation
	v.logger.Infof("Clearing access cache for user %s, course %s", userID, courseID)
	return nil
}

// UpdatePreviewSession updates a preview session duration
func (v *AccessValidator) UpdatePreviewSession(ctx context.Context, userID, lectureID string, durationSeconds int) error {
	v.logger.Infof("Updating preview session for user %s, lecture %s, duration %d", userID, lectureID, durationSeconds)

	// Get existing session
	session, err := v.previewRepo.GetByUserAndLecture(ctx, userID, lectureID)
	if err != nil {
		v.logger.Errorf("Failed to get preview session: %v", err)
		return err
	}

	// Update session duration
	session.SessionDurationSeconds += durationSeconds

	// Check if preview is exhausted
	if session.SessionDurationSeconds >= session.PreviewLimitSeconds {
		session.PreviewExhausted = true
	}

	// Update the session
	err = v.previewRepo.UpdatePreviewSession(session)
	if err != nil {
		v.logger.Errorf("Failed to update preview session: %v", err)
		return err
	}

	return nil
}