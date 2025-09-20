package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/study-platform/course-service/internal/service"
	pb "github.com/study-platform/course-service/proto"
	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccessValidationHandler struct {
	pb.UnimplementedAccessValidationServiceServer
	courseService *service.CourseService
	logger        logger.Logger
}

func NewAccessValidationHandler(courseService *service.CourseService, logger logger.Logger) *AccessValidationHandler {
	return &AccessValidationHandler{
		courseService: courseService,
		logger:        logger,
	}
}

// ValidateCourseAccess validates if a user has access to a specific course
func (h *AccessValidationHandler) ValidateCourseAccess(ctx context.Context, req *pb.ValidateCourseAccessRequest) (*pb.ValidateCourseAccessResponse, error) {
	h.logger.Info("ValidateCourseAccess request received")

	// Parse user ID
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	// Parse course ID
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid course ID: %v", err)
	}

	// Check if course exists
	course, err := h.courseService.GetCourse(ctx, courseID)
	if err != nil {
		h.logger.Errorf("Course not found: %v", err)
		return &pb.ValidateCourseAccessResponse{
			HasAccess:    false,
			AccessLevel:  "none",
			Message:      "Course not found",
			ErrorCode:    "COURSE_NOT_FOUND",
		}, nil
	}

	// Check if course is published
	if course.Status != "published" {
		// Only instructors and admins can access unpublished courses
		if req.UserRole != "instructor" && req.UserRole != "admin" {
			return &pb.ValidateCourseAccessResponse{
				HasAccess:    false,
				AccessLevel:  "none",
				Message:      "Course is not published",
				ErrorCode:    "COURSE_NOT_PUBLISHED",
			}, nil
		}
	}

	// Check if user is the course instructor
	if course.InstructorID == userID {
		return &pb.ValidateCourseAccessResponse{
			HasAccess:    true,
			AccessLevel:  "owner",
			Message:      "Course owner access granted",
			ExpiresAt:    0, // No expiration for owners
		}, nil
	}

	// Check if user is admin
	if req.UserRole == "admin" {
		return &pb.ValidateCourseAccessResponse{
			HasAccess:    true,
			AccessLevel:  "admin",
			Message:      "Admin access granted",
			ExpiresAt:    0, // No expiration for admins
		}, nil
	}

	// For regular users, check enrollment
	enrollment, err := h.courseService.GetEnrollment(ctx, userID, courseID)
	if err != nil {
		// No enrollment found
		if course.Price > 0 {
			return &pb.ValidateCourseAccessResponse{
				HasAccess:    false,
				AccessLevel:  "none",
				Message:      "Payment required to access this course",
				ErrorCode:    "PAYMENT_REQUIRED",
				CoursePrice:  course.Price,
				CourseCurrency: course.Currency,
			}, nil
		}
		// Free course but not enrolled
		return &pb.ValidateCourseAccessResponse{
			HasAccess:    false,
			AccessLevel:  "none",
			Message:      "Enrollment required to access this course",
			ErrorCode:    "ENROLLMENT_REQUIRED",
		}, nil
	}

	// Check if enrollment is active
	if enrollment.PaidAt == nil && course.Price > 0 {
		return &pb.ValidateCourseAccessResponse{
			HasAccess:    false,
			AccessLevel:  "none",
			Message:      "Payment required - enrollment not paid",
			ErrorCode:    "PAYMENT_REQUIRED",
			CoursePrice:  course.Price,
			CourseCurrency: course.Currency,
		}, nil
	}

	// Determine access level
	accessLevel := "student"
	expiresAt := int64(0) // No expiration by default

	// Note: UpdateLastAccessed method would need to be implemented in course service
	// For now, we skip this functionality

	return &pb.ValidateCourseAccessResponse{
		HasAccess:    true,
		AccessLevel:  accessLevel,
		Message:      "Access granted",
		ExpiresAt:    expiresAt,
		EnrollmentId: enrollment.ID.String(),
	}, nil
}

// ValidateLectureAccess validates if a user has access to a specific lecture
func (h *AccessValidationHandler) ValidateLectureAccess(ctx context.Context, req *pb.ValidateLectureAccessRequest) (*pb.ValidateLectureAccessResponse, error) {
	h.logger.Info("ValidateLectureAccess request received")

	// Parse user ID - validation only
	_, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	// Parse lecture ID
	lectureID, err := uuid.Parse(req.LectureId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid lecture ID: %v", err)
	}

	// Get lecture details
	lecture, err := h.courseService.GetLecture(ctx, lectureID)
	if err != nil {
		h.logger.Errorf("Lecture not found: %v", err)
		return &pb.ValidateLectureAccessResponse{
			HasAccess:    false,
			AccessLevel:  "none",
			Message:      "Lecture not found",
			ErrorCode:    "LECTURE_NOT_FOUND",
		}, nil
	}

	// First validate course access
	courseAccessReq := &pb.ValidateCourseAccessRequest{
		UserId:   req.UserId,
		CourseId: lecture.CourseID.String(),
		UserRole: req.UserRole,
	}

	courseAccessResp, err := h.ValidateCourseAccess(ctx, courseAccessReq)
	if err != nil {
		return nil, err
	}

	if !courseAccessResp.HasAccess {
		return &pb.ValidateLectureAccessResponse{
			HasAccess:    false,
			AccessLevel:  courseAccessResp.AccessLevel,
			Message:      courseAccessResp.Message,
			ErrorCode:    courseAccessResp.ErrorCode,
			CoursePrice:  courseAccessResp.CoursePrice,
			CourseCurrency: courseAccessResp.CourseCurrency,
		}, nil
	}

	// Check if lecture is published
	if lecture.Status != "published" && courseAccessResp.AccessLevel != "owner" && courseAccessResp.AccessLevel != "admin" {
		return &pb.ValidateLectureAccessResponse{
			HasAccess:    false,
			AccessLevel:  "none",
			Message:      "Lecture is not published",
			ErrorCode:    "LECTURE_NOT_PUBLISHED",
		}, nil
	}

	// Check sequential access if required
	if lecture.OrderNumber > 1 && req.EnforceSequentialAccess {
		// Note: CheckPreviousLecturesCompleted method would need to be implemented
		// For now, we skip sequential access validation
		h.logger.Infof("Sequential access check requested for lecture %d but not implemented", lecture.OrderNumber)
	}

	// Note: Time-based access restrictions (AvailableFrom, AvailableUntil)
	// are not implemented in the current Lecture model
	// These would need to be added as fields if required

	return &pb.ValidateLectureAccessResponse{
		HasAccess:     true,
		AccessLevel:   courseAccessResp.AccessLevel,
		Message:       "Lecture access granted",
		CourseId:      lecture.CourseID.String(),
		LectureOrder:  lecture.OrderNumber,
		Duration:      lecture.DurationMinutes,
	}, nil
}

// ValidateResourceAccess validates access to course resources (files, materials, etc.)
func (h *AccessValidationHandler) ValidateResourceAccess(ctx context.Context, req *pb.ValidateResourceAccessRequest) (*pb.ValidateResourceAccessResponse, error) {
	h.logger.Info("ValidateResourceAccess request received")

	// Parse user ID - validation only
	_, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	// Parse resource ID
	resourceID, err := uuid.Parse(req.ResourceId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid resource ID: %v", err)
	}

	// Note: Resource management is not yet implemented in the course service
	// This would require a Resource model and related methods
	h.logger.Infof("Resource access validation requested for resource %s but not implemented", resourceID.String())

	return &pb.ValidateResourceAccessResponse{
		HasAccess:      false,
		AccessLevel:    "none",
		Message:        "Resource validation not implemented",
		ErrorCode:      "NOT_IMPLEMENTED",
	}, nil
}

// BatchValidateAccess validates access to multiple courses/lectures/resources at once
func (h *AccessValidationHandler) BatchValidateAccess(ctx context.Context, req *pb.BatchValidateAccessRequest) (*pb.BatchValidateAccessResponse, error) {
	h.logger.Info("BatchValidateAccess request received")

	responses := make([]*pb.BatchAccessValidationResult, 0)

	for _, item := range req.Items {
		var result *pb.BatchAccessValidationResult

		switch item.ResourceType {
		case "course":
			courseReq := &pb.ValidateCourseAccessRequest{
				UserId:   req.UserId,
				CourseId: item.ResourceId,
				UserRole: req.UserRole,
			}
			courseResp, err := h.ValidateCourseAccess(ctx, courseReq)
			if err != nil {
				result = &pb.BatchAccessValidationResult{
					ResourceId:   item.ResourceId,
					ResourceType: "course",
					HasAccess:    false,
					Message:      "Validation error",
					ErrorCode:    "VALIDATION_ERROR",
				}
			} else {
				result = &pb.BatchAccessValidationResult{
					ResourceId:   item.ResourceId,
					ResourceType: "course",
					HasAccess:    courseResp.HasAccess,
					AccessLevel:  courseResp.AccessLevel,
					Message:      courseResp.Message,
					ErrorCode:    courseResp.ErrorCode,
				}
			}

		case "lecture":
			lectureReq := &pb.ValidateLectureAccessRequest{
				UserId:    req.UserId,
				LectureId: item.ResourceId,
				UserRole:  req.UserRole,
			}
			lectureResp, err := h.ValidateLectureAccess(ctx, lectureReq)
			if err != nil {
				result = &pb.BatchAccessValidationResult{
					ResourceId:   item.ResourceId,
					ResourceType: "lecture",
					HasAccess:    false,
					Message:      "Validation error",
					ErrorCode:    "VALIDATION_ERROR",
				}
			} else {
				result = &pb.BatchAccessValidationResult{
					ResourceId:   item.ResourceId,
					ResourceType: "lecture",
					HasAccess:    lectureResp.HasAccess,
					AccessLevel:  lectureResp.AccessLevel,
					Message:      lectureResp.Message,
					ErrorCode:    lectureResp.ErrorCode,
				}
			}

		default:
			result = &pb.BatchAccessValidationResult{
				ResourceId:   item.ResourceId,
				ResourceType: item.ResourceType,
				HasAccess:    false,
				Message:      "Unsupported resource type",
				ErrorCode:    "UNSUPPORTED_RESOURCE_TYPE",
			}
		}

		responses = append(responses, result)
	}

	return &pb.BatchValidateAccessResponse{
		Results: responses,
	}, nil
}