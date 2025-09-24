package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/study-platform/course-service/internal/model"
	"github.com/study-platform/course-service/internal/repository"
	"github.com/study-platform/pkg/logger"
)

type CourseService struct {
	courseRepo         *repository.CourseRepository
	lectureRepo        *repository.LectureRepository
	enrollmentRepo     *repository.EnrollmentRepository
	courseResourceRepo *repository.CourseResourceRepository
	logger             logger.Logger
}

func NewCourseService(
	courseRepo *repository.CourseRepository,
	lectureRepo *repository.LectureRepository,
	enrollmentRepo *repository.EnrollmentRepository,
	courseResourceRepo *repository.CourseResourceRepository,
	logger logger.Logger,
) *CourseService {
	return &CourseService{
		courseRepo:         courseRepo,
		lectureRepo:        lectureRepo,
		enrollmentRepo:     enrollmentRepo,
		courseResourceRepo: courseResourceRepo,
		logger:             logger,
	}
}

func (s *CourseService) CreateCourse(ctx context.Context, course *model.Course) error {
	s.logger.Infof("Creating new course: %s (instructor: %s)", course.Title, course.InstructorID.String())
	
	// Validate course data
	if err := s.validateCourse(course); err != nil {
		s.logger.Errorf("Course validation failed: %v", err)
		return err
	}
	
	// Set default values
	if course.Status == "" {
		course.Status = model.CourseStatusDraft
	}
	course.EnrollmentCount = 0
	course.Rating = 0
	course.RatingCount = 0
	course.DurationMinutes = 0

	if course.Currency == "" {
		course.Currency = "USD"
	}

	// Set default values for new fields
	if course.DifficultyLevel == "" {
		course.DifficultyLevel = "intermediate"
	}
	if course.Language == "" {
		course.Language = "en"
	}
	if course.EstimatedDurationHours == 0 {
		course.EstimatedDurationHours = 10 // default 10 hours
	}
	// Set default boolean values
	course.AutoApproveEnrollment = true
	course.AllowPreviews = true
	course.HasCertificate = false
	course.MobileAccess = true
	
	err := s.courseRepo.Create(ctx, course)
	if err != nil {
		s.logger.Errorf("Failed to create course: %v", err)
		return fmt.Errorf("failed to create course: %w", err)
	}
	
	s.logger.Infof("Course created successfully: %s", course.ID.String())
	return nil
}

func (s *CourseService) GetCourse(ctx context.Context, id uuid.UUID) (*model.Course, error) {
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to get course %s: %v", id.String(), err)
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	// Fetch lectures for the course
	lectures, err := s.lectureRepo.GetByCourseID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to get lectures for course %s: %v", id.String(), err)
		// Don't fail the whole request, just log the error
		lectures = []model.Lecture{}
	}
	course.Lectures = lectures

	// Fetch resources for the course
	resources, err := s.courseResourceRepo.GetByCourseID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to get resources for course %s: %v", id.String(), err)
		// Don't fail the whole request, just log the error
		resources = []model.CourseResource{}
	}
	course.Resources = resources

	return course, nil
}

func (s *CourseService) UpdateCourse(ctx context.Context, course *model.Course) error {
	s.logger.Infof("Updating course: %s", course.ID.String())

	// Check if course exists and is not deleted
	existingCourse, err := s.courseRepo.GetByID(ctx, course.ID)
	if err != nil {
		s.logger.Errorf("Course not found for update %s: %v", course.ID.String(), err)
		return fmt.Errorf("course not found: %w", err)
	}

	if existingCourse.DeletedAt != nil {
		return fmt.Errorf("cannot update deleted course")
	}

	// Validate course data
	if err := s.validateCourse(course); err != nil {
		s.logger.Errorf("Course validation failed: %v", err)
		return err
	}

	err = s.courseRepo.Update(ctx, course)
	if err != nil {
		s.logger.Errorf("Failed to update course %s: %v", course.ID.String(), err)
		return fmt.Errorf("failed to update course: %w", err)
	}

	s.logger.Infof("Course updated successfully: %s", course.ID.String())
	return nil
}

// UpdateCourseWithS3Resources provides enhanced course update with S3 file handling
func (s *CourseService) UpdateCourseWithS3Resources(ctx context.Context, course *model.Course, oldThumbnailURL string) error {
	s.logger.Infof("Updating course with S3 resource management: %s", course.ID.String())

	// Check if course exists and is not deleted
	existingCourse, err := s.courseRepo.GetByID(ctx, course.ID)
	if err != nil {
		s.logger.Errorf("Course not found for update %s: %v", course.ID.String(), err)
		return fmt.Errorf("course not found: %w", err)
	}

	if existingCourse.DeletedAt != nil {
		return fmt.Errorf("cannot update deleted course")
	}

	// Validate course data
	if err := s.validateCourse(course); err != nil {
		s.logger.Errorf("Course validation failed: %v", err)
		return err
	}

	// TODO: Implement S3 file cleanup logic here if thumbnail URL changed
	// This would involve:
	// 1. Comparing old and new thumbnail URLs
	// 2. If different, marking old file for cleanup in S3
	// 3. This should be done asynchronously to avoid blocking the update

	if oldThumbnailURL != "" && oldThumbnailURL != course.ThumbnailURL {
		s.logger.Infof("Thumbnail URL changed for course %s - old: %s, new: %s",
			course.ID.String(), oldThumbnailURL, course.ThumbnailURL)
		// Here you would call bucket service to clean up old thumbnail
		// For now, we just log it
	}

	err = s.courseRepo.Update(ctx, course)
	if err != nil {
		s.logger.Errorf("Failed to update course %s: %v", course.ID.String(), err)
		return fmt.Errorf("failed to update course: %w", err)
	}

	s.logger.Infof("Course updated successfully with S3 resource management: %s", course.ID.String())
	return nil
}

func (s *CourseService) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	s.logger.Infof("Soft deleting course: %s", id.String())

	// Check if course exists and is not already deleted
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to get course %s: %v", id.String(), err)
		return fmt.Errorf("course not found: %w", err)
	}

	if course.DeletedAt != nil {
		return fmt.Errorf("course is already deleted")
	}

	// Use cascade soft delete to preserve billing history
	// This will soft delete the course and all related lectures and enrollments
	err = s.courseRepo.SoftDeleteCourseWithCascade(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to soft delete course %s: %v", id.String(), err)
		return fmt.Errorf("failed to delete course: %w", err)
	}

	s.logger.Infof("Course soft deleted successfully: %s", id.String())
	return nil
}

func (s *CourseService) ListCourses(ctx context.Context, filters model.CourseFilters) (*model.CourseSearchResult, error) {
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 20
	}
	
	result, err := s.courseRepo.List(ctx, filters)
	if err != nil {
		s.logger.Errorf("Failed to list courses: %v", err)
		return nil, fmt.Errorf("failed to list courses: %w", err)
	}
	
	return result, nil
}

func (s *CourseService) SearchCourses(ctx context.Context, filters model.CourseFilters) (*model.CourseSearchResult, error) {
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 20
	}
	
	result, err := s.courseRepo.Search(ctx, filters)
	if err != nil {
		s.logger.Errorf("Failed to search courses: %v", err)
		return nil, fmt.Errorf("failed to search courses: %w", err)
	}
	
	return result, nil
}

func (s *CourseService) CreateLecture(ctx context.Context, lecture *model.Lecture) error {
	s.logger.Infof("Creating new lecture: %s (course: %s)", lecture.Title, lecture.CourseID.String())
	
	// Validate lecture data
	if err := s.validateLecture(lecture); err != nil {
		s.logger.Errorf("Lecture validation failed: %v", err)
		return err
	}
	
	// Verify course exists
	_, err := s.courseRepo.GetByID(ctx, lecture.CourseID)
	if err != nil {
		s.logger.Errorf("Course not found for lecture %s: %v", lecture.CourseID.String(), err)
		return fmt.Errorf("course not found: %w", err)
	}
	
	// Set default values
	lecture.Status = model.LectureStatusDraft
	
	err = s.lectureRepo.Create(ctx, lecture)
	if err != nil {
		s.logger.Errorf("Failed to create lecture: %v", err)
		return fmt.Errorf("failed to create lecture: %w", err)
	}
	
	// Update course duration
	err = s.lectureRepo.UpdateCourseDuration(ctx, lecture.CourseID)
	if err != nil {
		s.logger.Errorf("Failed to update course duration %s: %v", lecture.CourseID.String(), err)
	}
	
	s.logger.Infof("Lecture created successfully: %s", lecture.ID.String())
	return nil
}

func (s *CourseService) GetLecture(ctx context.Context, id uuid.UUID) (*model.Lecture, error) {
	lecture, err := s.lectureRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to get lecture %s: %v", id.String(), err)
		return nil, fmt.Errorf("failed to get lecture: %w", err)
	}
	
	return lecture, nil
}

func (s *CourseService) UpdateLecture(ctx context.Context, lecture *model.Lecture) error {
	s.logger.Infof("Updating lecture: %s", lecture.ID.String())
	
	// Validate lecture data
	if err := s.validateLecture(lecture); err != nil {
		s.logger.Errorf("Lecture validation failed: %v", err)
		return err
	}
	
	err := s.lectureRepo.Update(ctx, lecture)
	if err != nil {
		s.logger.Errorf("Failed to update lecture %s: %v", lecture.ID.String(), err)
		return fmt.Errorf("failed to update lecture: %w", err)
	}
	
	// Update course duration
	err = s.lectureRepo.UpdateCourseDuration(ctx, lecture.CourseID)
	if err != nil {
		s.logger.Errorf("Failed to update course duration %s: %v", lecture.CourseID.String(), err)
	}
	
	s.logger.Infof("Lecture updated successfully: %s", lecture.ID.String())
	return nil
}

func (s *CourseService) DeleteLecture(ctx context.Context, id uuid.UUID) error {
	s.logger.Infof("Deleting lecture: %s", id.String())
	
	// Get lecture to get course ID for duration update
	lecture, err := s.lectureRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to get lecture for deletion %s: %v", id.String(), err)
		return fmt.Errorf("failed to get lecture: %w", err)
	}
	
	err = s.lectureRepo.Delete(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to delete lecture %s: %v", id.String(), err)
		return fmt.Errorf("failed to delete lecture: %w", err)
	}
	
	// Update course duration
	err = s.lectureRepo.UpdateCourseDuration(ctx, lecture.CourseID)
	if err != nil {
		s.logger.Errorf("Failed to update course duration %s: %v", lecture.CourseID.String(), err)
	}
	
	s.logger.Infof("Lecture deleted successfully: %s", id.String())
	return nil
}

func (s *CourseService) ListLectures(ctx context.Context, filters model.LectureFilters) (*model.LectureSearchResult, error) {
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 50
	}
	
	result, err := s.lectureRepo.List(ctx, filters)
	if err != nil {
		s.logger.Errorf("Failed to list lectures: %v", err)
		return nil, fmt.Errorf("failed to list lectures: %w", err)
	}
	
	return result, nil
}

func (s *CourseService) EnrollInCourse(ctx context.Context, enrollment *model.Enrollment) error {
	s.logger.Infof("Enrolling user %s in course %s", enrollment.UserID.String(), enrollment.CourseID.String())

	// Verify course exists and is published
	course, err := s.courseRepo.GetByID(ctx, enrollment.CourseID)
	if err != nil {
		s.logger.Errorf("Course not found for enrollment %s: %v", enrollment.CourseID.String(), err)
		return fmt.Errorf("course not found: %w", err)
	}

	if course.Status != model.CourseStatusPublished {
		return fmt.Errorf("course is not available for enrollment")
	}

	// Check if user is already enrolled
	existingEnrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, enrollment.UserID, enrollment.CourseID)
	if err == nil && existingEnrollment != nil {
		return fmt.Errorf("user is already enrolled in this course")
	}

	// Set payment requirements based on course type
	if course.IsPaid && course.Price > 0 {
		// CRITICAL SECURITY FIX: For paid courses, DO NOT create enrollment without payment verification
		s.logger.Warnf("PAYMENT REQUIRED: User %s attempting to enroll in paid course %s (Price: %.2f %s) - ENROLLMENT DENIED",
			enrollment.UserID.String(), enrollment.CourseID.String(), course.Price, course.Currency)

		// Return error immediately - no enrollment should be created for paid courses
		// Enrollment will only be created after successful payment verification through the payment service
		return fmt.Errorf("payment required: course costs %.2f %s. Please purchase this course through the payment system to access content", course.Price, course.Currency)
	} else {
		// Free course - allow direct enrollment
		enrollment.PaymentRequired = false
		enrollment.PaymentStatus = "not_required"
		enrollment.Status = model.EnrollmentStatusActive

		err = s.enrollmentRepo.Create(ctx, enrollment)
		if err != nil {
			s.logger.Errorf("Failed to create enrollment: %v", err)
			return fmt.Errorf("failed to enroll in course: %w", err)
		}

		// Update course enrollment count
		err = s.courseRepo.UpdateEnrollmentCount(ctx, enrollment.CourseID)
		if err != nil {
			s.logger.Errorf("Failed to update enrollment count %s: %v", enrollment.CourseID.String(), err)
		}

		s.logger.Infof("User enrolled successfully in free course: %s", enrollment.ID.String())
	}

	return nil
}

func (s *CourseService) GetEnrollment(ctx context.Context, userID, courseID uuid.UUID) (*model.Enrollment, error) {
	enrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		s.logger.Errorf("Failed to get enrollment (user: %s, course: %s): %v", userID.String(), courseID.String(), err)
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}
	
	return enrollment, nil
}

func (s *CourseService) ListEnrollments(ctx context.Context, filters model.EnrollmentFilters) (*model.EnrollmentSearchResult, error) {
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 20
	}

	result, err := s.enrollmentRepo.List(ctx, filters)
	if err != nil {
		s.logger.Errorf("Failed to list enrollments: %v", err)
		return nil, fmt.Errorf("failed to list enrollments: %w", err)
	}

	return result, nil
}

// CreatePaidEnrollment creates an enrollment after successful payment verification
func (s *CourseService) CreatePaidEnrollment(ctx context.Context, userID, courseID uuid.UUID, orderID string, paidAmount float64, currency string) error {
	s.logger.Infof("Creating paid enrollment - User: %s, Course: %s, Order: %s", userID.String(), courseID.String(), orderID)

	// Verify course exists and is paid
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		s.logger.Errorf("Course not found for paid enrollment %s: %v", courseID.String(), err)
		return fmt.Errorf("course not found: %w", err)
	}

	if !course.IsPaid || course.Price <= 0 {
		return fmt.Errorf("course is not a paid course")
	}

	// Check if user is already enrolled
	existingEnrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err == nil && existingEnrollment != nil {
		s.logger.Warnf("User %s is already enrolled in course %s", userID.String(), courseID.String())
		return fmt.Errorf("user is already enrolled in this course")
	}

	// Create enrollment with payment information
	now := time.Now()
	enrollment := &model.Enrollment{
		UserID:                userID,
		CourseID:              courseID,
		Status:                model.EnrollmentStatusActive,
		PaymentRequired:       true,
		PaymentStatus:         "completed",
		LemonSqueezyOrderID:   &orderID,
		PaymentAmount:         &paidAmount,
		PaymentCurrency:       &currency,
		PaidAt:                &now,
	}

	err = s.enrollmentRepo.Create(ctx, enrollment)
	if err != nil {
		s.logger.Errorf("Failed to create paid enrollment: %v", err)
		return fmt.Errorf("failed to create enrollment: %w", err)
	}

	// Update course enrollment count
	err = s.courseRepo.UpdateEnrollmentCount(ctx, courseID)
	if err != nil {
		s.logger.Errorf("Failed to update enrollment count %s: %v", courseID.String(), err)
	}

	s.logger.Infof("Paid enrollment created successfully: %s", enrollment.ID.String())
	return nil
}

// CompleteEnrollmentPayment activates an enrollment after successful payment (legacy method, kept for compatibility)
func (s *CourseService) CompleteEnrollmentPayment(ctx context.Context, userID, courseID uuid.UUID, orderID string, paidAmount float64, currency string) error {
	s.logger.Infof("Completing payment for enrollment - User: %s, Course: %s, Order: %s", userID.String(), courseID.String(), orderID)

	// Try to find existing pending enrollment first
	enrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		// No existing enrollment found, create new paid enrollment
		return s.CreatePaidEnrollment(ctx, userID, courseID, orderID, paidAmount, currency)
	}

	if enrollment.Status != model.EnrollmentStatusPending {
		// If enrollment exists but not pending, check if payment was already completed
		if enrollment.Status == model.EnrollmentStatusActive && enrollment.PaymentStatus == "completed" {
			s.logger.Infof("Enrollment already active and paid for User: %s, Course: %s", userID.String(), courseID.String())
			return nil
		}
		return fmt.Errorf("enrollment is not in pending status (current: %s)", enrollment.Status)
	}

	if !enrollment.PaymentRequired {
		return fmt.Errorf("payment not required for this enrollment")
	}

	// Update enrollment with payment information
	now := time.Now()
	enrollment.Status = model.EnrollmentStatusActive
	enrollment.PaymentStatus = "completed"
	enrollment.LemonSqueezyOrderID = &orderID
	enrollment.PaymentAmount = &paidAmount
	enrollment.PaymentCurrency = &currency
	enrollment.PaidAt = &now

	err = s.enrollmentRepo.Update(ctx, enrollment)
	if err != nil {
		s.logger.Errorf("Failed to update enrollment after payment: %v", err)
		return fmt.Errorf("failed to activate enrollment: %w", err)
	}

	// Update course enrollment count
	err = s.courseRepo.UpdateEnrollmentCount(ctx, enrollment.CourseID)
	if err != nil {
		s.logger.Errorf("Failed to update enrollment count %s: %v", enrollment.CourseID.String(), err)
	}

	s.logger.Infof("Enrollment activated after successful payment: %s", enrollment.ID.String())
	return nil
}

func (s *CourseService) validateCourse(course *model.Course) error {
	if course.Title == "" {
		return fmt.Errorf("course title is required")
	}
	if course.Description == "" {
		return fmt.Errorf("course description is required")
	}
	if course.InstructorID == uuid.Nil {
		return fmt.Errorf("instructor ID is required")
	}
	if course.Category == "" {
		return fmt.Errorf("course category is required")
	}
	if course.Level == "" {
		return fmt.Errorf("course level is required")
	}
	if course.Price < 0 {
		return fmt.Errorf("course price cannot be negative")
	}
	
	// Validate level
	switch course.Level {
	case model.CourseLevelBeginner, model.CourseLevelIntermediate, model.CourseLevelAdvanced:
		// Valid level
	default:
		return fmt.Errorf("invalid course level: %s", course.Level)
	}
	
	return nil
}

func (s *CourseService) validateLecture(lecture *model.Lecture) error {
	if lecture.Title == "" {
		return fmt.Errorf("lecture title is required")
	}
	if lecture.CourseID == uuid.Nil {
		return fmt.Errorf("course ID is required")
	}
	if lecture.OrderNumber <= 0 {
		return fmt.Errorf("lecture order number must be positive")
	}
	if lecture.DurationMinutes < 0 {
		return fmt.Errorf("lecture duration cannot be negative")
	}

	return nil
}

// ValidateCourseAccess checks if a user has access to a course
func (s *CourseService) ValidateCourseAccess(ctx context.Context, userID, courseID uuid.UUID) (*model.CourseAccessResult, error) {
	// Get course details
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		s.logger.Errorf("Failed to get course %s: %v", courseID.String(), err)
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	// If course is free, grant full access
	if !course.IsPaid {
		return &model.CourseAccessResult{
			HasAccess:   true,
			AccessLevel: model.AccessLevelFull,
			CourseType:  model.CourseTypeFree,
			Message:     "Free course - full access granted",
		}, nil
	}

	// For paid courses, check enrollment
	enrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		// No enrollment found, check if preview is allowed
		return &model.CourseAccessResult{
			HasAccess:   false,
			AccessLevel: model.AccessLevelPreview,
			CourseType:  model.CourseTypePaid,
			CoursePrice: course.Price,
			Currency:    course.Currency,
			Message:     "Purchase required to access this content",
		}, nil
	}

	// Check enrollment status
	if enrollment.Status == model.EnrollmentStatusActive {
		return &model.CourseAccessResult{
			HasAccess:   true,
			AccessLevel: model.AccessLevelFull,
			CourseType:  model.CourseTypePaid,
			Message:     "Course purchased - full access granted",
		}, nil
	}

	// Enrollment exists but not active
	return &model.CourseAccessResult{
		HasAccess:   false,
		AccessLevel: model.AccessLevelDenied,
		CourseType:  model.CourseTypePaid,
		CoursePrice: course.Price,
		Currency:    course.Currency,
		Message:     "Enrollment inactive - purchase required",
	}, nil
}

// ValidateLectureAccess checks if a user has access to a specific lecture
func (s *CourseService) ValidateLectureAccess(ctx context.Context, userID, lectureID uuid.UUID) (*model.LectureAccessResult, error) {
	// Get lecture details
	lecture, err := s.lectureRepo.GetByID(ctx, lectureID)
	if err != nil {
		s.logger.Errorf("Failed to get lecture %s: %v", lectureID.String(), err)
		return nil, fmt.Errorf("failed to get lecture: %w", err)
	}

	// Get course access first
	courseAccess, err := s.ValidateCourseAccess(ctx, userID, lecture.CourseID)
	if err != nil {
		return nil, err
	}

	// If lecture is marked as free, allow access
	if lecture.IsFree {
		return &model.LectureAccessResult{
			HasAccess:     true,
			AccessLevel:   model.AccessLevelFull,
			LectureType:   model.LectureTypeFree,
			CourseAccess:  courseAccess,
			Message:       "Free lecture - full access granted",
		}, nil
	}

	// If user has full course access, grant lecture access
	if courseAccess.HasAccess && courseAccess.AccessLevel == model.AccessLevelFull {
		return &model.LectureAccessResult{
			HasAccess:     true,
			AccessLevel:   model.AccessLevelFull,
			LectureType:   model.LectureTypePaid,
			CourseAccess:  courseAccess,
			Message:       "Full course access - lecture available",
		}, nil
	}

	// For preview access, check if this is a preview lecture (first lecture or marked as preview)
	if s.isPreviewLecture(ctx, lecture) {
		return &model.LectureAccessResult{
			HasAccess:     true,
			AccessLevel:   model.AccessLevelPreview,
			LectureType:   model.LectureTypePreview,
			CourseAccess:  courseAccess,
			PreviewTimeLimit: 600, // 10 minutes preview
			Message:       "Preview access - limited time available",
		}, nil
	}

	// No access
	return &model.LectureAccessResult{
		HasAccess:     false,
		AccessLevel:   model.AccessLevelDenied,
		LectureType:   model.LectureTypePaid,
		CourseAccess:  courseAccess,
		Message:       "Purchase required to access this lecture",
	}, nil
}

// isPreviewLecture determines if a lecture should be available for preview
func (s *CourseService) isPreviewLecture(ctx context.Context, lecture *model.Lecture) bool {
	// Check if it's the first lecture in the course
	lectures, err := s.lectureRepo.GetByCourseID(ctx, lecture.CourseID)
	if err != nil {
		s.logger.Errorf("Failed to get course lectures: %v", err)
		return false
	}

	// If this is the first lecture (lowest order number), allow preview
	if len(lectures) > 0 {
		firstLecture := lectures[0]
		for _, l := range lectures {
			if l.OrderNumber < firstLecture.OrderNumber {
				firstLecture = l
			}
		}
		if lecture.ID == firstLecture.ID {
			return true
		}
	}

	return false
}