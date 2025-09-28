package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/study-platform/course-service/internal/model"
	"github.com/study-platform/course-service/internal/repository"
	"github.com/study-platform/pkg/logger"
)

type LectureResourceService struct {
	lectureRepo         *repository.LectureRepository
	lectureResourceRepo *repository.LectureResourceRepository
	enrollmentRepo      *repository.EnrollmentRepository
	logger              logger.Logger
}

func NewLectureResourceService(
	lectureRepo *repository.LectureRepository,
	lectureResourceRepo *repository.LectureResourceRepository,
	enrollmentRepo *repository.EnrollmentRepository,
	logger logger.Logger,
) *LectureResourceService {
	return &LectureResourceService{
		lectureRepo:         lectureRepo,
		lectureResourceRepo: lectureResourceRepo,
		enrollmentRepo:      enrollmentRepo,
		logger:              logger,
	}
}

// CreateResource creates a new resource for a lecture
func (s *LectureResourceService) CreateResource(ctx context.Context, resource *model.LectureResource) error {
	// Validate lecture exists
	_, err := s.lectureRepo.GetByID(ctx, resource.LectureID)
	if err != nil {
		return fmt.Errorf("lecture not found: %w", err)
	}

	// Create the resource
	err = s.lectureResourceRepo.Create(ctx, resource)
	if err != nil {
		s.logger.Errorf("Failed to create lecture resource: %v", err)
		return fmt.Errorf("failed to create resource: %w", err)
	}

	s.logger.Infof("Created resource %s for lecture %s", resource.ID, resource.LectureID)
	return nil
}

// GetResourcesByLecture gets all resources for a specific lecture
func (s *LectureResourceService) GetResourcesByLecture(ctx context.Context, lectureID uuid.UUID) ([]model.LectureResource, error) {
	// Validate lecture exists
	_, err := s.lectureRepo.GetByID(ctx, lectureID)
	if err != nil {
		return nil, fmt.Errorf("lecture not found: %w", err)
	}

	resources, err := s.lectureResourceRepo.GetByLectureID(ctx, lectureID)
	if err != nil {
		s.logger.Errorf("Failed to get resources for lecture %s: %v", lectureID, err)
		return nil, fmt.Errorf("failed to get resources: %w", err)
	}

	return resources, nil
}

// GetResourcesByCourse gets all resources for all lectures in a course
func (s *LectureResourceService) GetResourcesByCourse(ctx context.Context, courseID uuid.UUID) (map[uuid.UUID][]model.LectureResource, error) {
	resources, err := s.lectureResourceRepo.GetByCourseID(ctx, courseID)
	if err != nil {
		s.logger.Errorf("Failed to get resources for course %s: %v", courseID, err)
		return nil, fmt.Errorf("failed to get course resources: %w", err)
	}

	return resources, nil
}

// GetResource gets a specific resource by ID
func (s *LectureResourceService) GetResource(ctx context.Context, resourceID uuid.UUID) (*model.LectureResource, error) {
	resource, err := s.lectureResourceRepo.GetByID(ctx, resourceID)
	if err != nil {
		s.logger.Errorf("Failed to get resource %s: %v", resourceID, err)
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	return resource, nil
}

// UpdateResource updates an existing lecture resource
func (s *LectureResourceService) UpdateResource(ctx context.Context, resource *model.LectureResource) error {
	// Check if resource exists and get current data
	existing, err := s.lectureResourceRepo.GetByID(ctx, resource.ID)
	if err != nil {
		return fmt.Errorf("resource not found: %w", err)
	}

	// Validate lecture exists if lecture_id is being changed
	if existing.LectureID != resource.LectureID {
		_, err := s.lectureRepo.GetByID(ctx, resource.LectureID)
		if err != nil {
			return fmt.Errorf("target lecture not found: %w", err)
		}
	}

	// Update the resource
	err = s.lectureResourceRepo.Update(ctx, resource)
	if err != nil {
		s.logger.Errorf("Failed to update resource %s: %v", resource.ID, err)
		return fmt.Errorf("failed to update resource: %w", err)
	}

	s.logger.Infof("Updated resource %s", resource.ID)
	return nil
}

// DeleteResource deletes a lecture resource
func (s *LectureResourceService) DeleteResource(ctx context.Context, resourceID uuid.UUID) error {
	// Check if resource exists
	_, err := s.lectureResourceRepo.GetByID(ctx, resourceID)
	if err != nil {
		return fmt.Errorf("resource not found: %w", err)
	}

	err = s.lectureResourceRepo.Delete(ctx, resourceID)
	if err != nil {
		s.logger.Errorf("Failed to delete resource %s: %v", resourceID, err)
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	s.logger.Infof("Deleted resource %s", resourceID)
	return nil
}

// BulkCreateResources creates multiple resources for a lecture
func (s *LectureResourceService) BulkCreateResources(ctx context.Context, lectureID uuid.UUID, resources []model.LectureResource) error {
	// Validate lecture exists
	_, err := s.lectureRepo.GetByID(ctx, lectureID)
	if err != nil {
		return fmt.Errorf("lecture not found: %w", err)
	}

	// Set lecture ID for all resources
	for i := range resources {
		resources[i].LectureID = lectureID
	}

	err = s.lectureResourceRepo.BulkCreate(ctx, lectureID, resources)
	if err != nil {
		s.logger.Errorf("Failed to bulk create resources for lecture %s: %v", lectureID, err)
		return fmt.Errorf("failed to bulk create resources: %w", err)
	}

	s.logger.Infof("Created %d resources for lecture %s", len(resources), lectureID)
	return nil
}

// ReorderResources updates the display order of multiple resources
func (s *LectureResourceService) ReorderResources(ctx context.Context, lectureID uuid.UUID, resourceOrders []struct {
	ResourceID   uuid.UUID `json:"resource_id"`
	DisplayOrder int32     `json:"display_order"`
}) error {
	// Validate lecture exists
	_, err := s.lectureRepo.GetByID(ctx, lectureID)
	if err != nil {
		return fmt.Errorf("lecture not found: %w", err)
	}

	// Validate all resources belong to the lecture
	for _, order := range resourceOrders {
		resource, err := s.lectureResourceRepo.GetByID(ctx, order.ResourceID)
		if err != nil {
			return fmt.Errorf("resource %s not found: %w", order.ResourceID, err)
		}
		if resource.LectureID != lectureID {
			return fmt.Errorf("resource %s does not belong to lecture %s", order.ResourceID, lectureID)
		}
	}

	// Prepare updates for repository
	updates := make([]struct {
		ID           uuid.UUID
		DisplayOrder int32
	}, len(resourceOrders))

	for i, order := range resourceOrders {
		updates[i] = struct {
			ID           uuid.UUID
			DisplayOrder int32
		}{
			ID:           order.ResourceID,
			DisplayOrder: order.DisplayOrder,
		}
	}

	err = s.lectureResourceRepo.UpdateDisplayOrder(ctx, updates)
	if err != nil {
		s.logger.Errorf("Failed to reorder resources for lecture %s: %v", lectureID, err)
		return fmt.Errorf("failed to reorder resources: %w", err)
	}

	s.logger.Infof("Reordered %d resources for lecture %s", len(resourceOrders), lectureID)
	return nil
}

// DeleteResourcesByLecture deletes all resources for a lecture
func (s *LectureResourceService) DeleteResourcesByLecture(ctx context.Context, lectureID uuid.UUID) error {
	err := s.lectureResourceRepo.DeleteByLectureID(ctx, lectureID)
	if err != nil {
		s.logger.Errorf("Failed to delete resources for lecture %s: %v", lectureID, err)
		return fmt.Errorf("failed to delete lecture resources: %w", err)
	}

	s.logger.Infof("Deleted all resources for lecture %s", lectureID)
	return nil
}

// GetLectureWithResources gets a lecture with its resources populated
func (s *LectureResourceService) GetLectureWithResources(ctx context.Context, lectureID uuid.UUID) (*model.Lecture, error) {
	lecture, err := s.lectureRepo.GetByID(ctx, lectureID)
	if err != nil {
		return nil, fmt.Errorf("lecture not found: %w", err)
	}

	resources, err := s.lectureResourceRepo.GetByLectureID(ctx, lectureID)
	if err != nil {
		s.logger.Errorf("Failed to get resources for lecture %s: %v", lectureID, err)
		// Don't fail the entire request if resources can't be loaded
		resources = []model.LectureResource{}
	}

	lecture.Resources = resources
	return lecture, nil
}

// PopulateLecturesWithResources populates a slice of lectures with their resources
func (s *LectureResourceService) PopulateLecturesWithResources(ctx context.Context, lectures []*model.Lecture) error {
	if len(lectures) == 0 {
		s.logger.Infof("DEBUG: No lectures to populate with resources")
		return nil
	}

	// Get course ID from the first lecture (assuming all lectures are from the same course)
	courseID := lectures[0].CourseID
	s.logger.Infof("DEBUG: Populating %d lectures with resources for course %s", len(lectures), courseID)

	// Get all resources for the course grouped by lecture
	resourcesMap, err := s.lectureResourceRepo.GetByCourseID(ctx, courseID)
	if err != nil {
		s.logger.Errorf("Failed to get resources for course %s: %v", courseID, err)
		// Don't fail the entire request if resources can't be loaded
		return nil
	}

	s.logger.Infof("DEBUG: Found %d lecture resource groups for course %s", len(resourcesMap), courseID)

	// Populate each lecture with its resources
	for i := range lectures {
		if resources, exists := resourcesMap[lectures[i].ID]; exists {
			lectures[i].Resources = resources
			s.logger.Infof("DEBUG: Populated lecture %s with %d resources", lectures[i].ID, len(resources))
		} else {
			lectures[i].Resources = []model.LectureResource{}
			s.logger.Infof("DEBUG: No resources found for lecture %s", lectures[i].ID)
		}
	}

	return nil
}

// CheckUserCourseAccess checks if a user has access to a course by lecture ID
func (s *LectureResourceService) CheckUserCourseAccess(ctx context.Context, userID uuid.UUID, lectureID uuid.UUID) (bool, error) {
	// Get the lecture to find the course ID
	lecture, err := s.lectureRepo.GetByID(ctx, lectureID)
	if err != nil {
		s.logger.Errorf("Failed to get lecture %s: %v", lectureID, err)
		return false, fmt.Errorf("failed to get lecture: %w", err)
	}

	// Check if user is enrolled in the course
	enrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, lecture.CourseID)
	if err != nil {
		// If enrollment not found, user doesn't have access
		s.logger.Infof("User %s not enrolled in course %s", userID, lecture.CourseID)
		return false, nil
	}

	// Check if enrollment is active
	if enrollment.Status != "active" {
		s.logger.Infof("User %s has inactive enrollment in course %s: %s", userID, lecture.CourseID, enrollment.Status)
		return false, nil
	}

	s.logger.Infof("User %s has access to course %s via lecture %s", userID, lecture.CourseID, lectureID)
	return true, nil
}