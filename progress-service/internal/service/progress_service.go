package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/study-platform/progress-service/internal/model"
	"github.com/study-platform/progress-service/internal/repository"
	"github.com/study-platform/pkg/logger"
)

type ProgressService struct {
	progressRepo *repository.ProgressRepository
	logger       logger.Logger
}

func NewProgressService(progressRepo *repository.ProgressRepository, logger logger.Logger) *ProgressService {
	return &ProgressService{
		progressRepo: progressRepo,
		logger:       logger,
	}
}

// Progress tracking methods
func (s *ProgressService) UpdateProgress(userID, courseID, lectureID uuid.UUID, progressPercentage float64, watchTimeSeconds int32, isCompleted bool) (*model.UserProgress, error) {
	// Check if progress already exists
	existingProgress, err := s.progressRepo.GetProgress(userID, courseID, lectureID)
	now := time.Now()
	
	if err != nil && err.Error() != "progress not found" {
		return nil, fmt.Errorf("failed to check existing progress: %w", err)
	}
	
	if existingProgress != nil {
		// Update existing progress
		existingProgress.ProgressPercentage = progressPercentage
		existingProgress.WatchTimeSeconds = watchTimeSeconds
		existingProgress.IsCompleted = isCompleted
		existingProgress.LastAccessed = now
		existingProgress.UpdatedAt = now
		
		if isCompleted && existingProgress.CompletedAt == nil {
			existingProgress.CompletedAt = &now
		}
		
		err = s.progressRepo.UpdateProgress(existingProgress)
		if err != nil {
			return nil, fmt.Errorf("failed to update progress: %w", err)
		}
		
		// Update enrollment progress
		s.updateEnrollmentProgress(userID, courseID)
		
		return existingProgress, nil
	}
	
	// Create new progress
	progress := &model.UserProgress{
		ID:                 uuid.New(),
		UserID:             userID,
		CourseID:           courseID,
		LectureID:          lectureID,
		ProgressPercentage: progressPercentage,
		WatchTimeSeconds:   watchTimeSeconds,
		IsCompleted:        isCompleted,
		LastAccessed:       now,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	
	if isCompleted {
		progress.CompletedAt = &now
	}
	
	err = s.progressRepo.CreateProgress(progress)
	if err != nil {
		return nil, fmt.Errorf("failed to create progress: %w", err)
	}
	
	// Update enrollment progress
	s.updateEnrollmentProgress(userID, courseID)
	
	return progress, nil
}

func (s *ProgressService) GetProgress(userID, courseID, lectureID uuid.UUID) (*model.UserProgress, error) {
	progress, err := s.progressRepo.GetProgress(userID, courseID, lectureID)
	if err != nil {
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}
	return progress, nil
}

func (s *ProgressService) GetUserProgress(userID, courseID uuid.UUID) ([]*model.UserProgress, float64, error) {
	progressList, err := s.progressRepo.GetUserProgress(userID, courseID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user progress: %w", err)
	}
	
	// Calculate overall progress percentage
	if len(progressList) == 0 {
		return progressList, 0, nil
	}
	
	totalProgress := 0.0
	for _, p := range progressList {
		totalProgress += p.ProgressPercentage
	}
	overallProgress := totalProgress / float64(len(progressList))
	
	return progressList, overallProgress, nil
}

func (s *ProgressService) GetCourseProgress(courseID uuid.UUID, page, pageSize int) ([]*model.UserProgress, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	
	progressList, totalCount, err := s.progressRepo.GetCourseProgress(courseID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get course progress: %w", err)
	}
	
	return progressList, totalCount, nil
}

// Enrollment methods
func (s *ProgressService) CreateEnrollment(userID, courseID uuid.UUID) (*model.Enrollment, error) {
	// Check if enrollment already exists
	existingEnrollment, err := s.progressRepo.GetEnrollment(userID, courseID)
	if err != nil && err.Error() != "enrollment not found" {
		return nil, fmt.Errorf("failed to check existing enrollment: %w", err)
	}
	
	if existingEnrollment != nil {
		return nil, fmt.Errorf("user is already enrolled in this course")
	}
	
	now := time.Now()
	enrollment := &model.Enrollment{
		ID:                    uuid.New(),
		UserID:                userID,
		CourseID:              courseID,
		Status:                "enrolled",
		ProgressPercentage:    0.0,
		CompletedLectures:     0,
		TotalLectures:         0, // This should be set based on actual course data
		TotalWatchTimeSeconds: 0,
		EnrolledAt:            now,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	
	err = s.progressRepo.CreateEnrollment(enrollment)
	if err != nil {
		return nil, fmt.Errorf("failed to create enrollment: %w", err)
	}
	
	return enrollment, nil
}

func (s *ProgressService) GetEnrollment(userID, courseID uuid.UUID) (*model.Enrollment, error) {
	enrollment, err := s.progressRepo.GetEnrollment(userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}
	return enrollment, nil
}

func (s *ProgressService) ListEnrollments(userID uuid.UUID, status string, page, pageSize int) ([]*model.Enrollment, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	
	enrollments, totalCount, err := s.progressRepo.ListEnrollments(userID, status, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list enrollments: %w", err)
	}
	
	return enrollments, totalCount, nil
}

func (s *ProgressService) UpdateEnrollmentStatus(userID, courseID uuid.UUID, status string) (*model.Enrollment, error) {
	// Validate status
	validStatuses := map[string]bool{
		"enrolled":  true,
		"completed": true,
		"cancelled": true,
		"suspended": true,
	}
	
	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid status: %s", status)
	}
	
	err := s.progressRepo.UpdateEnrollmentStatus(userID, courseID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to update enrollment status: %w", err)
	}
	
	// Return updated enrollment
	enrollment, err := s.progressRepo.GetEnrollment(userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated enrollment: %w", err)
	}
	
	return enrollment, nil
}

// Course completion methods
func (s *ProgressService) MarkLectureComplete(userID, courseID, lectureID uuid.UUID, watchTimeSeconds int32) (*model.UserProgress, bool, error) {
	// Mark lecture as completed
	progress, err := s.UpdateProgress(userID, courseID, lectureID, 100.0, watchTimeSeconds, true)
	if err != nil {
		return nil, false, fmt.Errorf("failed to mark lecture complete: %w", err)
	}
	
	// Check if course is completed
	courseCompleted := s.checkCourseCompletion(userID, courseID)
	
	return progress, courseCompleted, nil
}

func (s *ProgressService) GetLectureProgress(userID, courseID uuid.UUID) ([]*model.LectureProgress, error) {
	lectureProgress, err := s.progressRepo.GetLectureProgress(userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lecture progress: %w", err)
	}
	return lectureProgress, nil
}

func (s *ProgressService) GetCourseCompletion(userID, courseID uuid.UUID) (*model.CourseCompletion, error) {
	completion, err := s.progressRepo.GetCourseCompletion(userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course completion: %w", err)
	}
	return completion, nil
}

// Analytics methods
func (s *ProgressService) GetUserAnalytics(userID uuid.UUID) (*model.UserAnalytics, error) {
	analytics, err := s.progressRepo.GetUserAnalytics(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user analytics: %w", err)
	}
	return analytics, nil
}

func (s *ProgressService) GetCourseAnalytics(courseID uuid.UUID) (*model.CourseAnalytics, error) {
	analytics, err := s.progressRepo.GetCourseAnalytics(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course analytics: %w", err)
	}
	return analytics, nil
}

func (s *ProgressService) GetInstructorAnalytics(instructorID uuid.UUID) (*model.InstructorAnalytics, error) {
	analytics, err := s.progressRepo.GetInstructorAnalytics(instructorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get instructor analytics: %w", err)
	}
	return analytics, nil
}

// Helper methods
func (s *ProgressService) updateEnrollmentProgress(userID, courseID uuid.UUID) {
	// Get all progress for this user and course
	progressList, err := s.progressRepo.GetUserProgress(userID, courseID)
	if err != nil {
		s.logger.Errorf("Failed to get user progress for enrollment update: %v", err)
		return
	}
	
	if len(progressList) == 0 {
		return
	}
	
	// Calculate overall progress
	totalProgress := 0.0
	completedLectures := int32(0)
	totalWatchTime := int32(0)
	
	for _, p := range progressList {
		totalProgress += p.ProgressPercentage
		totalWatchTime += p.WatchTimeSeconds
		if p.IsCompleted {
			completedLectures++
		}
	}
	
	overallProgress := totalProgress / float64(len(progressList))
	
	// Update enrollment
	err = s.progressRepo.UpdateEnrollmentProgress(userID, courseID, overallProgress, completedLectures, totalWatchTime)
	if err != nil {
		s.logger.Errorf("Failed to update enrollment progress: %v", err)
	}
}

func (s *ProgressService) checkCourseCompletion(userID, courseID uuid.UUID) bool {
	completion, err := s.progressRepo.GetCourseCompletion(userID, courseID)
	if err != nil {
		s.logger.Errorf("Failed to check course completion: %v", err)
		return false
	}
	
	return completion.CompletionPercentage >= 100.0
}