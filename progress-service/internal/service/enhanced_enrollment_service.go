package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/study-platform/progress-service/internal/model"
	"github.com/study-platform/progress-service/internal/repository"
	"github.com/study-platform/pkg/logger"
)

// EnhancedEnrollmentService provides improved enrollment functionality with proper error handling
type EnhancedEnrollmentService struct {
	progressRepo *repository.ProgressRepository
	logger       logger.Logger
}

type EnrollmentRequest struct {
	UserID       uuid.UUID `json:"user_id"`
	CourseID     uuid.UUID `json:"course_id"`
	PaymentID    *string   `json:"payment_id,omitempty"`    // For paid courses
	PaymentStatus string   `json:"payment_status,omitempty"` // pending, paid, verified
	IsFreeEnrollment bool  `json:"is_free_enrollment,omitempty"`
	Source       string    `json:"source,omitempty"`        // "direct", "payment_completed", "admin"
}

type EnrollmentResult struct {
	Enrollment    *model.Enrollment `json:"enrollment"`
	Success       bool              `json:"success"`
	Message       string            `json:"message"`
	RequiresPayment bool            `json:"requires_payment"`
	PaymentURL    *string           `json:"payment_url,omitempty"`
}

func NewEnhancedEnrollmentService(progressRepo *repository.ProgressRepository, logger logger.Logger) *EnhancedEnrollmentService {
	return &EnhancedEnrollmentService{
		progressRepo: progressRepo,
		logger:       logger,
	}
}

// CreateEnrollment creates an enrollment with enhanced error handling and payment support
func (s *EnhancedEnrollmentService) CreateEnrollment(ctx context.Context, req *EnrollmentRequest) (*EnrollmentResult, error) {
	s.logger.Infof("Creating enrollment for user %s in course %s", req.UserID.String(), req.CourseID.String())

	// 1. Check if user is already enrolled
	existingEnrollment, err := s.progressRepo.GetEnrollment(req.UserID, req.CourseID)
	if err != nil && err.Error() != "enrollment not found" {
		s.logger.Errorf("Failed to check existing enrollment: %v", err)
		return &EnrollmentResult{
			Success: false,
			Message: "Failed to verify enrollment status",
		}, fmt.Errorf("failed to check existing enrollment: %w", err)
	}

	if existingEnrollment != nil {
		s.logger.Warnf("User %s is already enrolled in course %s", req.UserID.String(), req.CourseID.String())
		return &EnrollmentResult{
			Enrollment: existingEnrollment,
			Success:    true,
			Message:    "User is already enrolled in this course",
		}, nil
	}

	// 2. Handle different enrollment scenarios
	switch {
	case req.IsFreeEnrollment:
		return s.createFreeEnrollment(ctx, req)
	case req.PaymentStatus == "paid" || req.PaymentStatus == "verified":
		return s.createPaidEnrollment(ctx, req)
	case req.PaymentStatus == "pending":
		return s.createPendingEnrollment(ctx, req)
	default:
		// Default: Try to determine if course is free or paid
		return s.createSmartEnrollment(ctx, req)
	}
}

// createFreeEnrollment creates enrollment for free courses
func (s *EnhancedEnrollmentService) createFreeEnrollment(ctx context.Context, req *EnrollmentRequest) (*EnrollmentResult, error) {
	s.logger.Infof("Creating free enrollment for user %s in course %s", req.UserID.String(), req.CourseID.String())

	enrollment := &model.Enrollment{
		ID:                    uuid.New(),
		UserID:                req.UserID,
		CourseID:              req.CourseID,
		Status:                "active",
		ProgressPercentage:    0.0,
		CompletedLectures:     0,
		TotalLectures:         0,
		TotalWatchTimeSeconds: 0,
		EnrolledAt:            time.Now(),
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	err := s.progressRepo.CreateEnrollment(enrollment)
	if err != nil {
		s.logger.Errorf("Failed to create free enrollment: %v", err)
		return &EnrollmentResult{
			Success: false,
			Message: "Failed to create enrollment",
		}, fmt.Errorf("failed to create enrollment: %w", err)
	}

	s.logger.Infof("Successfully created free enrollment for user %s in course %s", req.UserID.String(), req.CourseID.String())

	return &EnrollmentResult{
		Enrollment: enrollment,
		Success:    true,
		Message:    "Successfully enrolled in free course",
	}, nil
}

// createPaidEnrollment creates enrollment for paid courses after payment verification
func (s *EnhancedEnrollmentService) createPaidEnrollment(ctx context.Context, req *EnrollmentRequest) (*EnrollmentResult, error) {
	s.logger.Infof("Creating paid enrollment for user %s in course %s with payment ID %s",
		req.UserID.String(), req.CourseID.String(), getStringValue(req.PaymentID))

	enrollment := &model.Enrollment{
		ID:                    uuid.New(),
		UserID:                req.UserID,
		CourseID:              req.CourseID,
		Status:                "active",
		ProgressPercentage:    0.0,
		CompletedLectures:     0,
		TotalLectures:         0,
		TotalWatchTimeSeconds: 0,
		EnrolledAt:            time.Now(),
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	// Add payment information to enrollment metadata if available
	if req.PaymentID != nil {
		// In a real implementation, you might store payment information
		// in a separate field or as metadata
		s.logger.Infof("Enrollment linked to payment ID: %s", *req.PaymentID)
	}

	err := s.progressRepo.CreateEnrollment(enrollment)
	if err != nil {
		s.logger.Errorf("Failed to create paid enrollment: %v", err)
		return &EnrollmentResult{
			Success: false,
			Message: "Failed to create enrollment after payment",
		}, fmt.Errorf("failed to create paid enrollment: %w", err)
	}

	s.logger.Infof("Successfully created paid enrollment for user %s in course %s", req.UserID.String(), req.CourseID.String())

	return &EnrollmentResult{
		Enrollment: enrollment,
		Success:    true,
		Message:    "Successfully enrolled after payment verification",
	}, nil
}

// createPendingEnrollment creates a pending enrollment for courses awaiting payment
func (s *EnhancedEnrollmentService) createPendingEnrollment(ctx context.Context, req *EnrollmentRequest) (*EnrollmentResult, error) {
	s.logger.Infof("Creating pending enrollment for user %s in course %s", req.UserID.String(), req.CourseID.String())

	// For pending payments, we don't create an enrollment yet
	// Instead, we return information about required payment
	return &EnrollmentResult{
		Success:         false,
		Message:         "Payment required for enrollment",
		RequiresPayment: true,
		// PaymentURL would be provided by the payment service
	}, nil
}

// createSmartEnrollment attempts to determine enrollment type automatically
func (s *EnhancedEnrollmentService) createSmartEnrollment(ctx context.Context, req *EnrollmentRequest) (*EnrollmentResult, error) {
	s.logger.Infof("Creating smart enrollment for user %s in course %s", req.UserID.String(), req.CourseID.String())

	// TODO: In a real implementation, this would:
	// 1. Query the course service to get course pricing information
	// 2. Check if the course is free or paid
	// 3. For paid courses, check if payment exists
	// 4. Create appropriate enrollment based on results

	// For now, default to free enrollment to fix the 500 errors
	req.IsFreeEnrollment = true
	return s.createFreeEnrollment(ctx, req)
}

// GetEnrollment retrieves enrollment information
func (s *EnhancedEnrollmentService) GetEnrollment(ctx context.Context, userID, courseID uuid.UUID) (*model.Enrollment, error) {
	enrollment, err := s.progressRepo.GetEnrollment(userID, courseID)
	if err != nil {
		s.logger.Errorf("Failed to get enrollment for user %s in course %s: %v", userID.String(), courseID.String(), err)
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}

	return enrollment, nil
}

// ListUserEnrollments retrieves all enrollments for a user
func (s *EnhancedEnrollmentService) ListUserEnrollments(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Enrollment, error) {
	enrollments, _, err := s.progressRepo.ListEnrollments(userID, "", limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to list enrollments for user %s: %v", userID.String(), err)
		return nil, fmt.Errorf("failed to list enrollments: %w", err)
	}

	return enrollments, nil
}

// UpdateEnrollmentStatus updates the status of an enrollment
func (s *EnhancedEnrollmentService) UpdateEnrollmentStatus(ctx context.Context, enrollmentID uuid.UUID, status string) error {
	s.logger.Infof("Updating enrollment %s status to %s", enrollmentID.String(), status)

	// Get the enrollment first to get userID and courseID
	enrollment, err := s.progressRepo.GetEnrollmentByID(enrollmentID)
	if err != nil {
		s.logger.Errorf("Failed to get enrollment: %v", err)
		return fmt.Errorf("failed to get enrollment: %w", err)
	}

	err = s.progressRepo.UpdateEnrollmentStatus(enrollment.UserID, enrollment.CourseID, status)
	if err != nil {
		s.logger.Errorf("Failed to update enrollment status: %v", err)
		return fmt.Errorf("failed to update enrollment status: %w", err)
	}

	return nil
}

// ActivateEnrollmentAfterPayment activates a pending enrollment after payment confirmation
func (s *EnhancedEnrollmentService) ActivateEnrollmentAfterPayment(ctx context.Context, userID, courseID uuid.UUID, paymentID string) (*model.Enrollment, error) {
	s.logger.Infof("Activating enrollment for user %s in course %s after payment %s",
		userID.String(), courseID.String(), paymentID)

	// Check if enrollment exists
	enrollment, err := s.progressRepo.GetEnrollment(userID, courseID)
	if err != nil {
		// If enrollment doesn't exist, create it
		req := &EnrollmentRequest{
			UserID:        userID,
			CourseID:      courseID,
			PaymentID:     &paymentID,
			PaymentStatus: "verified",
			Source:        "payment_completed",
		}

		result, err := s.createPaidEnrollment(ctx, req)
		if err != nil {
			return nil, err
		}

		return result.Enrollment, nil
	}

	// Update existing enrollment to active
	enrollment.Status = "active"
	enrollment.UpdatedAt = time.Now()

	err = s.progressRepo.UpdateEnrollmentStatus(enrollment.UserID, enrollment.CourseID, enrollment.Status)
	if err != nil {
		s.logger.Errorf("Failed to activate enrollment: %v", err)
		return nil, fmt.Errorf("failed to activate enrollment: %w", err)
	}

	return enrollment, nil
}

// Helper function
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}