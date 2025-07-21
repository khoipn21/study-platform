package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/study-platform/course-service/internal/model"
	"github.com/study-platform/course-service/internal/repository"
	"github.com/study-platform/pkg/logger"
)

type CourseService struct {
	courseRepo     *repository.CourseRepository
	lectureRepo    *repository.LectureRepository
	enrollmentRepo *repository.EnrollmentRepository
	logger         *logger.Logger
}

func NewCourseService(
	courseRepo *repository.CourseRepository,
	lectureRepo *repository.LectureRepository,
	enrollmentRepo *repository.EnrollmentRepository,
	logger *logger.Logger,
) *CourseService {
	return &CourseService{
		courseRepo:     courseRepo,
		lectureRepo:    lectureRepo,
		enrollmentRepo: enrollmentRepo,
		logger:         logger,
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
	course.Status = model.CourseStatusDraft
	course.EnrollmentCount = 0
	course.Rating = 0
	course.RatingCount = 0
	course.DurationMinutes = 0
	
	if course.Currency == "" {
		course.Currency = "USD"
	}
	
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
	
	return course, nil
}

func (s *CourseService) UpdateCourse(ctx context.Context, course *model.Course) error {
	s.logger.Infof("Updating course: %s", course.ID.String())
	
	// Validate course data
	if err := s.validateCourse(course); err != nil {
		s.logger.Errorf("Course validation failed: %v", err)
		return err
	}
	
	err := s.courseRepo.Update(ctx, course)
	if err != nil {
		s.logger.Errorf("Failed to update course %s: %v", course.ID.String(), err)
		return fmt.Errorf("failed to update course: %w", err)
	}
	
	s.logger.Infof("Course updated successfully: %s", course.ID.String())
	return nil
}

func (s *CourseService) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	s.logger.Infof("Deleting course: %s", id.String())
	
	// Check if course has enrollments
	enrollments, err := s.enrollmentRepo.GetCourseEnrollments(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to check course enrollments %s: %v", id.String(), err)
		return fmt.Errorf("failed to check course enrollments: %w", err)
	}
	
	if len(enrollments) > 0 {
		return fmt.Errorf("cannot delete course with active enrollments")
	}
	
	err = s.courseRepo.Delete(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to delete course %s: %v", id.String(), err)
		return fmt.Errorf("failed to delete course: %w", err)
	}
	
	s.logger.Infof("Course deleted successfully: %s", id.String())
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
	
	s.logger.Infof("User enrolled successfully: %s", enrollment.ID.String())
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