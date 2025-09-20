package handler

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/study-platform/progress-service/internal/model"
	"github.com/study-platform/progress-service/internal/service"
	pb "github.com/study-platform/progress-service/proto"
	"github.com/study-platform/pkg/logger"
)

type EnhancedProgressHandlerFixed struct {
	pb.UnimplementedProgressServiceServer
	progressService            *service.ProgressService
	enhancedEnrollmentService  *service.EnhancedEnrollmentService
	logger                     logger.Logger
}

func NewEnhancedProgressHandlerFixed(
	progressService *service.ProgressService,
	enhancedEnrollmentService *service.EnhancedEnrollmentService,
	logger logger.Logger,
) *EnhancedProgressHandlerFixed {
	return &EnhancedProgressHandlerFixed{
		progressService:           progressService,
		enhancedEnrollmentService: enhancedEnrollmentService,
		logger:                    logger,
	}
}

// CreateEnrollment creates a new enrollment using the basic protobuf definition
func (h *EnhancedProgressHandlerFixed) CreateEnrollment(ctx context.Context, req *pb.CreateEnrollmentRequest) (*pb.CreateEnrollmentResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and course_id are required")
	}

	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id format")
	}

	// Create enrollment request with basic info
	enrollmentReq := &service.EnrollmentRequest{
		UserID:           userID,
		CourseID:         courseID,
		IsFreeEnrollment: true, // Default to free enrollment for now
		PaymentStatus:    "completed",
		Source:           "api_request",
	}

	// Create enrollment using enhanced service
	result, err := h.enhancedEnrollmentService.CreateEnrollment(ctx, enrollmentReq)
	if err != nil {
		h.logger.Errorf("Failed to create enrollment: %v", err)
		return nil, status.Error(codes.Internal, "failed to create enrollment")
	}

	// Handle enrollment result
	if !result.Success {
		return &pb.CreateEnrollmentResponse{
			Enrollment: nil,
			Message:    result.Message,
		}, nil
	}

	// Convert enrollment to proto
	enrollmentProto := h.convertEnrollmentToProto(result.Enrollment)

	return &pb.CreateEnrollmentResponse{
		Enrollment: enrollmentProto,
		Message:    "Enrollment created successfully",
	}, nil
}

// GetEnrollment retrieves an enrollment
func (h *EnhancedProgressHandlerFixed) GetEnrollment(ctx context.Context, req *pb.GetEnrollmentRequest) (*pb.GetEnrollmentResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and course_id are required")
	}

	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id format")
	}

	// Get enrollment
	enrollment, err := h.progressService.GetEnrollment(userID, courseID)
	if err != nil {
		h.logger.Errorf("Failed to get enrollment: %v", err)
		return nil, status.Error(codes.NotFound, "enrollment not found")
	}

	return &pb.GetEnrollmentResponse{
		Enrollment: h.convertEnrollmentToProto(enrollment),
	}, nil
}

// ListEnrollments lists enrollments for a user
func (h *EnhancedProgressHandlerFixed) ListEnrollments(ctx context.Context, req *pb.ListEnrollmentsRequest) (*pb.ListEnrollmentsResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Parse UUID
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	// Set defaults for pagination
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	// List enrollments
	enrollments, totalCount, err := h.progressService.ListEnrollments(userID, req.Status, page, pageSize)
	if err != nil {
		h.logger.Errorf("Failed to list enrollments: %v", err)
		return nil, status.Error(codes.Internal, "failed to list enrollments")
	}

	// Convert to proto
	protoEnrollments := make([]*pb.Enrollment, len(enrollments))
	for i, enrollment := range enrollments {
		protoEnrollments[i] = h.convertEnrollmentToProto(enrollment)
	}

	return &pb.ListEnrollmentsResponse{
		Enrollments: protoEnrollments,
		TotalCount:  int32(totalCount),
		Page:        int32(page),
		PageSize:    int32(pageSize),
	}, nil
}

// UpdateEnrollmentStatus updates enrollment status
func (h *EnhancedProgressHandlerFixed) UpdateEnrollmentStatus(ctx context.Context, req *pb.UpdateEnrollmentStatusRequest) (*pb.UpdateEnrollmentStatusResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" || req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id, course_id and status are required")
	}

	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id format")
	}

	// Update enrollment status
	enrollment, err := h.progressService.UpdateEnrollmentStatus(userID, courseID, req.Status)
	if err != nil {
		h.logger.Errorf("Failed to update enrollment status: %v", err)
		return nil, status.Error(codes.Internal, "failed to update enrollment status")
	}

	return &pb.UpdateEnrollmentStatusResponse{
		Enrollment: h.convertEnrollmentToProto(enrollment),
		Message:    "Enrollment status updated successfully",
	}, nil
}

// UpdateProgress updates user progress for a lecture
func (h *EnhancedProgressHandlerFixed) UpdateProgress(ctx context.Context, req *pb.UpdateProgressRequest) (*pb.UpdateProgressResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" || req.LectureId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id, course_id and lecture_id are required")
	}

	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id format")
	}

	lectureID, err := uuid.Parse(req.LectureId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid lecture_id format")
	}

	// Update progress
	progress, err := h.progressService.UpdateProgress(userID, courseID, lectureID, req.ProgressPercentage, req.WatchTimeSeconds, req.IsCompleted)
	if err != nil {
		h.logger.Errorf("Failed to update progress: %v", err)
		return nil, status.Error(codes.Internal, "failed to update progress")
	}

	return &pb.UpdateProgressResponse{
		Progress: h.convertProgressToProto(progress),
		Message:  "Progress updated successfully",
	}, nil
}

// GetProgress gets user progress for a lecture
func (h *EnhancedProgressHandlerFixed) GetProgress(ctx context.Context, req *pb.GetProgressRequest) (*pb.GetProgressResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" || req.LectureId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id, course_id and lecture_id are required")
	}

	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id format")
	}

	lectureID, err := uuid.Parse(req.LectureId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid lecture_id format")
	}

	// Get progress
	progress, err := h.progressService.GetProgress(userID, courseID, lectureID)
	if err != nil {
		h.logger.Errorf("Failed to get progress: %v", err)
		return nil, status.Error(codes.NotFound, "progress not found")
	}

	return &pb.GetProgressResponse{
		Progress: h.convertProgressToProto(progress),
	}, nil
}

// GetUserProgress gets all progress for a user in a course
func (h *EnhancedProgressHandlerFixed) GetUserProgress(ctx context.Context, req *pb.GetUserProgressRequest) (*pb.GetUserProgressResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and course_id are required")
	}

	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id format")
	}

	// Get user progress
	progressList, overallProgress, err := h.progressService.GetUserProgress(userID, courseID)
	if err != nil {
		h.logger.Errorf("Failed to get user progress: %v", err)
		return nil, status.Error(codes.Internal, "failed to get user progress")
	}

	// Convert to proto
	protoProgress := make([]*pb.UserProgress, len(progressList))
	for i, progress := range progressList {
		protoProgress[i] = h.convertProgressToProto(progress)
	}

	return &pb.GetUserProgressResponse{
		Progress:                 protoProgress,
		OverallProgressPercentage: overallProgress,
	}, nil
}

// Helper methods to convert models to protobuf
func (h *EnhancedProgressHandlerFixed) convertEnrollmentToProto(enrollment *model.Enrollment) *pb.Enrollment {
	proto := &pb.Enrollment{
		Id:                     enrollment.ID.String(),
		UserId:                 enrollment.UserID.String(),
		CourseId:               enrollment.CourseID.String(),
		Status:                 enrollment.Status,
		ProgressPercentage:     enrollment.ProgressPercentage,
		CompletedLectures:      enrollment.CompletedLectures,
		TotalLectures:          enrollment.TotalLectures,
		TotalWatchTimeSeconds:  enrollment.TotalWatchTimeSeconds,
		EnrolledAt:             timestamppb.New(enrollment.EnrolledAt),
		LastAccessed:           timestamppb.New(*enrollment.LastAccessed),
		CreatedAt:              timestamppb.New(enrollment.CreatedAt),
		UpdatedAt:              timestamppb.New(enrollment.UpdatedAt),
	}

	if enrollment.CompletedAt != nil {
		proto.CompletedAt = timestamppb.New(*enrollment.CompletedAt)
	}

	return proto
}

func (h *EnhancedProgressHandlerFixed) convertProgressToProto(progress *model.UserProgress) *pb.UserProgress {
	proto := &pb.UserProgress{
		Id:                 progress.ID.String(),
		UserId:             progress.UserID.String(),
		CourseId:           progress.CourseID.String(),
		LectureId:          progress.LectureID.String(),
		ProgressPercentage: progress.ProgressPercentage,
		WatchTimeSeconds:   progress.WatchTimeSeconds,
		IsCompleted:        progress.IsCompleted,
		LastAccessed:       timestamppb.New(progress.LastAccessed),
		CreatedAt:          timestamppb.New(progress.CreatedAt),
		UpdatedAt:          timestamppb.New(progress.UpdatedAt),
	}

	if progress.CompletedAt != nil {
		proto.CompletedAt = timestamppb.New(*progress.CompletedAt)
	}

	return proto
}

// Placeholder methods for other protobuf service methods
func (h *EnhancedProgressHandlerFixed) GetCourseProgress(ctx context.Context, req *pb.GetCourseProgressRequest) (*pb.GetCourseProgressResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

func (h *EnhancedProgressHandlerFixed) MarkLectureComplete(ctx context.Context, req *pb.MarkLectureCompleteRequest) (*pb.MarkLectureCompleteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

func (h *EnhancedProgressHandlerFixed) GetLectureProgress(ctx context.Context, req *pb.GetLectureProgressRequest) (*pb.GetLectureProgressResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

func (h *EnhancedProgressHandlerFixed) GetCourseCompletion(ctx context.Context, req *pb.GetCourseCompletionRequest) (*pb.GetCourseCompletionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

func (h *EnhancedProgressHandlerFixed) GetUserAnalytics(ctx context.Context, req *pb.GetUserAnalyticsRequest) (*pb.GetUserAnalyticsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

func (h *EnhancedProgressHandlerFixed) GetCourseAnalytics(ctx context.Context, req *pb.GetCourseAnalyticsRequest) (*pb.GetCourseAnalyticsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

func (h *EnhancedProgressHandlerFixed) GetInstructorAnalytics(ctx context.Context, req *pb.GetInstructorAnalyticsRequest) (*pb.GetInstructorAnalyticsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}