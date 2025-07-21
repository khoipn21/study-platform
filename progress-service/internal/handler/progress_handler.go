package handler

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/study-platform/progress-service/internal/service"
	pb "github.com/study-platform/progress-service/proto"
	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProgressHandler struct {
	pb.UnimplementedProgressServiceServer
	progressService *service.ProgressService
	logger          logger.Logger
}

func NewProgressHandler(progressService *service.ProgressService, logger logger.Logger) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
		logger:          logger,
	}
}

// Progress tracking methods
func (h *ProgressHandler) UpdateProgress(ctx context.Context, req *pb.UpdateProgressRequest) (*pb.UpdateProgressResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" || req.LectureId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id, course_id, and lecture_id are required")
	}
	
	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id")
	}
	
	lectureID, err := uuid.Parse(req.LectureId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid lecture_id")
	}
	
	// Validate progress percentage
	if req.ProgressPercentage < 0 || req.ProgressPercentage > 100 {
		return nil, status.Error(codes.InvalidArgument, "progress_percentage must be between 0 and 100")
	}
	
	// Update progress
	progress, err := h.progressService.UpdateProgress(userID, courseID, lectureID, req.ProgressPercentage, req.WatchTimeSeconds, req.IsCompleted)
	if err != nil {
		h.logger.Errorf("Failed to update progress: %v", err)
		return nil, status.Error(codes.Internal, "failed to update progress")
	}
	
	// Convert to proto
	progressProto := &pb.UserProgress{
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
		progressProto.CompletedAt = timestamppb.New(*progress.CompletedAt)
	}
	
	return &pb.UpdateProgressResponse{
		Progress: progressProto,
		Message:  "Progress updated successfully",
	}, nil
}

func (h *ProgressHandler) GetProgress(ctx context.Context, req *pb.GetProgressRequest) (*pb.GetProgressResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" || req.LectureId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id, course_id, and lecture_id are required")
	}
	
	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id")
	}
	
	lectureID, err := uuid.Parse(req.LectureId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid lecture_id")
	}
	
	// Get progress
	progress, err := h.progressService.GetProgress(userID, courseID, lectureID)
	if err != nil {
		if err.Error() == "progress not found" {
			return nil, status.Error(codes.NotFound, "progress not found")
		}
		h.logger.Errorf("Failed to get progress: %v", err)
		return nil, status.Error(codes.Internal, "failed to get progress")
	}
	
	// Convert to proto
	progressProto := &pb.UserProgress{
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
		progressProto.CompletedAt = timestamppb.New(*progress.CompletedAt)
	}
	
	return &pb.GetProgressResponse{
		Progress: progressProto,
	}, nil
}

func (h *ProgressHandler) GetUserProgress(ctx context.Context, req *pb.GetUserProgressRequest) (*pb.GetUserProgressResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and course_id are required")
	}
	
	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id")
	}
	
	// Get user progress
	progressList, overallProgress, err := h.progressService.GetUserProgress(userID, courseID)
	if err != nil {
		h.logger.Errorf("Failed to get user progress: %v", err)
		return nil, status.Error(codes.Internal, "failed to get user progress")
	}
	
	// Convert to proto
	var progressProtos []*pb.UserProgress
	for _, progress := range progressList {
		progressProto := &pb.UserProgress{
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
			progressProto.CompletedAt = timestamppb.New(*progress.CompletedAt)
		}
		
		progressProtos = append(progressProtos, progressProto)
	}
	
	return &pb.GetUserProgressResponse{
		Progress:                  progressProtos,
		OverallProgressPercentage: overallProgress,
	}, nil
}

func (h *ProgressHandler) GetCourseProgress(ctx context.Context, req *pb.GetCourseProgressRequest) (*pb.GetCourseProgressResponse, error) {
	// Validate request
	if req.CourseId == "" {
		return nil, status.Error(codes.InvalidArgument, "course_id is required")
	}
	
	// Parse UUID
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id")
	}
	
	// Default pagination
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	
	// Get course progress
	progressList, totalCount, err := h.progressService.GetCourseProgress(courseID, int(page), int(pageSize))
	if err != nil {
		h.logger.Errorf("Failed to get course progress: %v", err)
		return nil, status.Error(codes.Internal, "failed to get course progress")
	}
	
	// Convert to proto
	var progressProtos []*pb.UserProgress
	for _, progress := range progressList {
		progressProto := &pb.UserProgress{
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
			progressProto.CompletedAt = timestamppb.New(*progress.CompletedAt)
		}
		
		progressProtos = append(progressProtos, progressProto)
	}
	
	return &pb.GetCourseProgressResponse{
		Progress:   progressProtos,
		TotalCount: int32(totalCount),
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// Enrollment methods
func (h *ProgressHandler) CreateEnrollment(ctx context.Context, req *pb.CreateEnrollmentRequest) (*pb.CreateEnrollmentResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and course_id are required")
	}
	
	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id")
	}
	
	// Create enrollment
	enrollment, err := h.progressService.CreateEnrollment(userID, courseID)
	if err != nil {
		h.logger.Errorf("Failed to create enrollment: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create enrollment: %v", err))
	}
	
	// Convert to proto
	enrollmentProto := &pb.Enrollment{
		Id:                    enrollment.ID.String(),
		UserId:                enrollment.UserID.String(),
		CourseId:              enrollment.CourseID.String(),
		Status:                enrollment.Status,
		ProgressPercentage:    enrollment.ProgressPercentage,
		CompletedLectures:     enrollment.CompletedLectures,
		TotalLectures:         enrollment.TotalLectures,
		TotalWatchTimeSeconds: enrollment.TotalWatchTimeSeconds,
		EnrolledAt:            timestamppb.New(enrollment.EnrolledAt),
		CreatedAt:             timestamppb.New(enrollment.CreatedAt),
		UpdatedAt:             timestamppb.New(enrollment.UpdatedAt),
	}
	
	if enrollment.CompletedAt != nil {
		enrollmentProto.CompletedAt = timestamppb.New(*enrollment.CompletedAt)
	}
	
	if enrollment.LastAccessed != nil {
		enrollmentProto.LastAccessed = timestamppb.New(*enrollment.LastAccessed)
	}
	
	return &pb.CreateEnrollmentResponse{
		Enrollment: enrollmentProto,
		Message:    "Enrollment created successfully",
	}, nil
}

func (h *ProgressHandler) GetEnrollment(ctx context.Context, req *pb.GetEnrollmentRequest) (*pb.GetEnrollmentResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and course_id are required")
	}
	
	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id")
	}
	
	// Get enrollment
	enrollment, err := h.progressService.GetEnrollment(userID, courseID)
	if err != nil {
		if err.Error() == "enrollment not found" {
			return nil, status.Error(codes.NotFound, "enrollment not found")
		}
		h.logger.Errorf("Failed to get enrollment: %v", err)
		return nil, status.Error(codes.Internal, "failed to get enrollment")
	}
	
	// Convert to proto
	enrollmentProto := &pb.Enrollment{
		Id:                    enrollment.ID.String(),
		UserId:                enrollment.UserID.String(),
		CourseId:              enrollment.CourseID.String(),
		Status:                enrollment.Status,
		ProgressPercentage:    enrollment.ProgressPercentage,
		CompletedLectures:     enrollment.CompletedLectures,
		TotalLectures:         enrollment.TotalLectures,
		TotalWatchTimeSeconds: enrollment.TotalWatchTimeSeconds,
		EnrolledAt:            timestamppb.New(enrollment.EnrolledAt),
		CreatedAt:             timestamppb.New(enrollment.CreatedAt),
		UpdatedAt:             timestamppb.New(enrollment.UpdatedAt),
	}
	
	if enrollment.CompletedAt != nil {
		enrollmentProto.CompletedAt = timestamppb.New(*enrollment.CompletedAt)
	}
	
	if enrollment.LastAccessed != nil {
		enrollmentProto.LastAccessed = timestamppb.New(*enrollment.LastAccessed)
	}
	
	return &pb.GetEnrollmentResponse{
		Enrollment: enrollmentProto,
	}, nil
}

func (h *ProgressHandler) ListEnrollments(ctx context.Context, req *pb.ListEnrollmentsRequest) (*pb.ListEnrollmentsResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	
	// Parse UUID
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	
	// Default pagination
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	
	// List enrollments
	enrollments, totalCount, err := h.progressService.ListEnrollments(userID, req.Status, int(page), int(pageSize))
	if err != nil {
		h.logger.Errorf("Failed to list enrollments: %v", err)
		return nil, status.Error(codes.Internal, "failed to list enrollments")
	}
	
	// Convert to proto
	var enrollmentProtos []*pb.Enrollment
	for _, enrollment := range enrollments {
		enrollmentProto := &pb.Enrollment{
			Id:                    enrollment.ID.String(),
			UserId:                enrollment.UserID.String(),
			CourseId:              enrollment.CourseID.String(),
			Status:                enrollment.Status,
			ProgressPercentage:    enrollment.ProgressPercentage,
			CompletedLectures:     enrollment.CompletedLectures,
			TotalLectures:         enrollment.TotalLectures,
			TotalWatchTimeSeconds: enrollment.TotalWatchTimeSeconds,
			EnrolledAt:            timestamppb.New(enrollment.EnrolledAt),
			CreatedAt:             timestamppb.New(enrollment.CreatedAt),
			UpdatedAt:             timestamppb.New(enrollment.UpdatedAt),
		}
		
		if enrollment.CompletedAt != nil {
			enrollmentProto.CompletedAt = timestamppb.New(*enrollment.CompletedAt)
		}
		
		if enrollment.LastAccessed != nil {
			enrollmentProto.LastAccessed = timestamppb.New(*enrollment.LastAccessed)
		}
		
		enrollmentProtos = append(enrollmentProtos, enrollmentProto)
	}
	
	return &pb.ListEnrollmentsResponse{
		Enrollments: enrollmentProtos,
		TotalCount:  int32(totalCount),
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

func (h *ProgressHandler) UpdateEnrollmentStatus(ctx context.Context, req *pb.UpdateEnrollmentStatusRequest) (*pb.UpdateEnrollmentStatusResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" || req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id, course_id, and status are required")
	}
	
	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id")
	}
	
	// Update enrollment status
	enrollment, err := h.progressService.UpdateEnrollmentStatus(userID, courseID, req.Status)
	if err != nil {
		h.logger.Errorf("Failed to update enrollment status: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update enrollment status: %v", err))
	}
	
	// Convert to proto
	enrollmentProto := &pb.Enrollment{
		Id:                    enrollment.ID.String(),
		UserId:                enrollment.UserID.String(),
		CourseId:              enrollment.CourseID.String(),
		Status:                enrollment.Status,
		ProgressPercentage:    enrollment.ProgressPercentage,
		CompletedLectures:     enrollment.CompletedLectures,
		TotalLectures:         enrollment.TotalLectures,
		TotalWatchTimeSeconds: enrollment.TotalWatchTimeSeconds,
		EnrolledAt:            timestamppb.New(enrollment.EnrolledAt),
		CreatedAt:             timestamppb.New(enrollment.CreatedAt),
		UpdatedAt:             timestamppb.New(enrollment.UpdatedAt),
	}
	
	if enrollment.CompletedAt != nil {
		enrollmentProto.CompletedAt = timestamppb.New(*enrollment.CompletedAt)
	}
	
	if enrollment.LastAccessed != nil {
		enrollmentProto.LastAccessed = timestamppb.New(*enrollment.LastAccessed)
	}
	
	return &pb.UpdateEnrollmentStatusResponse{
		Enrollment: enrollmentProto,
		Message:    "Enrollment status updated successfully",
	}, nil
}

// Course completion methods
func (h *ProgressHandler) MarkLectureComplete(ctx context.Context, req *pb.MarkLectureCompleteRequest) (*pb.MarkLectureCompleteResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" || req.LectureId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id, course_id, and lecture_id are required")
	}
	
	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id")
	}
	
	lectureID, err := uuid.Parse(req.LectureId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid lecture_id")
	}
	
	// Mark lecture complete
	progress, courseCompleted, err := h.progressService.MarkLectureComplete(userID, courseID, lectureID, req.WatchTimeSeconds)
	if err != nil {
		h.logger.Errorf("Failed to mark lecture complete: %v", err)
		return nil, status.Error(codes.Internal, "failed to mark lecture complete")
	}
	
	// Convert to proto
	progressProto := &pb.UserProgress{
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
		progressProto.CompletedAt = timestamppb.New(*progress.CompletedAt)
	}
	
	message := "Lecture marked as complete"
	if courseCompleted {
		message = "Lecture marked as complete and course completed"
	}
	
	return &pb.MarkLectureCompleteResponse{
		Progress:        progressProto,
		CourseCompleted: courseCompleted,
		Message:         message,
	}, nil
}

func (h *ProgressHandler) GetLectureProgress(ctx context.Context, req *pb.GetLectureProgressRequest) (*pb.GetLectureProgressResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and course_id are required")
	}
	
	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id")
	}
	
	// Get lecture progress
	lectureProgress, err := h.progressService.GetLectureProgress(userID, courseID)
	if err != nil {
		h.logger.Errorf("Failed to get lecture progress: %v", err)
		return nil, status.Error(codes.Internal, "failed to get lecture progress")
	}
	
	// Convert to proto
	var lectureProgressProtos []*pb.LectureProgress
	for _, lp := range lectureProgress {
		lectureProgressProto := &pb.LectureProgress{
			LectureId:          lp.LectureID.String(),
			Title:              lp.Title,
			OrderNumber:        lp.OrderNumber,
			ProgressPercentage: lp.ProgressPercentage,
			WatchTimeSeconds:   lp.WatchTimeSeconds,
			IsCompleted:        lp.IsCompleted,
		}
		
		if lp.LastAccessed != nil {
			lectureProgressProto.LastAccessed = timestamppb.New(*lp.LastAccessed)
		}
		
		if lp.CompletedAt != nil {
			lectureProgressProto.CompletedAt = timestamppb.New(*lp.CompletedAt)
		}
		
		lectureProgressProtos = append(lectureProgressProtos, lectureProgressProto)
	}
	
	return &pb.GetLectureProgressResponse{
		LectureProgress: lectureProgressProtos,
	}, nil
}

func (h *ProgressHandler) GetCourseCompletion(ctx context.Context, req *pb.GetCourseCompletionRequest) (*pb.GetCourseCompletionResponse, error) {
	// Validate request
	if req.UserId == "" || req.CourseId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and course_id are required")
	}
	
	// Parse UUIDs
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id")
	}
	
	// Get course completion
	completion, err := h.progressService.GetCourseCompletion(userID, courseID)
	if err != nil {
		h.logger.Errorf("Failed to get course completion: %v", err)
		return nil, status.Error(codes.Internal, "failed to get course completion")
	}
	
	// Convert lecture progress to proto
	var lectureProgressProtos []*pb.LectureProgress
	for _, lp := range completion.LectureProgress {
		lectureProgressProto := &pb.LectureProgress{
			LectureId:          lp.LectureID.String(),
			Title:              lp.Title,
			OrderNumber:        lp.OrderNumber,
			ProgressPercentage: lp.ProgressPercentage,
			WatchTimeSeconds:   lp.WatchTimeSeconds,
			IsCompleted:        lp.IsCompleted,
		}
		
		if lp.LastAccessed != nil {
			lectureProgressProto.LastAccessed = timestamppb.New(*lp.LastAccessed)
		}
		
		if lp.CompletedAt != nil {
			lectureProgressProto.CompletedAt = timestamppb.New(*lp.CompletedAt)
		}
		
		lectureProgressProtos = append(lectureProgressProtos, lectureProgressProto)
	}
	
	// Convert to proto
	completionProto := &pb.CourseCompletion{
		CourseId:              completion.CourseID.String(),
		CourseTitle:           completion.CourseTitle,
		UserId:                completion.UserID.String(),
		CompletionPercentage:  completion.CompletionPercentage,
		CompletedLectures:     completion.CompletedLectures,
		TotalLectures:         completion.TotalLectures,
		TotalWatchTimeSeconds: completion.TotalWatchTimeSeconds,
		LectureProgress:       lectureProgressProtos,
	}
	
	if completion.StartedAt != nil {
		completionProto.StartedAt = timestamppb.New(*completion.StartedAt)
	}
	
	if completion.CompletedAt != nil {
		completionProto.CompletedAt = timestamppb.New(*completion.CompletedAt)
	}
	
	if completion.LastAccessed != nil {
		completionProto.LastAccessed = timestamppb.New(*completion.LastAccessed)
	}
	
	return &pb.GetCourseCompletionResponse{
		Completion: completionProto,
	}, nil
}

// Analytics methods
func (h *ProgressHandler) GetUserAnalytics(ctx context.Context, req *pb.GetUserAnalyticsRequest) (*pb.GetUserAnalyticsResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	
	// Parse UUID
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	
	// Get user analytics
	analytics, err := h.progressService.GetUserAnalytics(userID)
	if err != nil {
		h.logger.Errorf("Failed to get user analytics: %v", err)
		return nil, status.Error(codes.Internal, "failed to get user analytics")
	}
	
	// Convert to proto
	analyticsProto := &pb.UserAnalytics{
		UserId:                    analytics.UserID.String(),
		TotalCoursesEnrolled:      analytics.TotalCoursesEnrolled,
		TotalCoursesCompleted:     analytics.TotalCoursesCompleted,
		TotalLecturesCompleted:    analytics.TotalLecturesCompleted,
		TotalWatchTimeSeconds:     analytics.TotalWatchTimeSeconds,
		AverageProgressPercentage: analytics.AverageProgressPercentage,
		CoursesInProgress:         analytics.CoursesInProgress,
		MostActiveDay:             analytics.MostActiveDay,
		StreakDays:                analytics.StreakDays,
	}
	
	if analytics.LastActivity != nil {
		analyticsProto.LastActivity = timestamppb.New(*analytics.LastActivity)
	}
	
	return &pb.GetUserAnalyticsResponse{
		Analytics: analyticsProto,
	}, nil
}

func (h *ProgressHandler) GetCourseAnalytics(ctx context.Context, req *pb.GetCourseAnalyticsRequest) (*pb.GetCourseAnalyticsResponse, error) {
	// Validate request
	if req.CourseId == "" {
		return nil, status.Error(codes.InvalidArgument, "course_id is required")
	}
	
	// Parse UUID
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course_id")
	}
	
	// Get course analytics
	analytics, err := h.progressService.GetCourseAnalytics(courseID)
	if err != nil {
		h.logger.Errorf("Failed to get course analytics: %v", err)
		return nil, status.Error(codes.Internal, "failed to get course analytics")
	}
	
	// Convert to proto
	analyticsProto := &pb.CourseAnalytics{
		CourseId:                  analytics.CourseID.String(),
		TotalEnrollments:          analytics.TotalEnrollments,
		TotalCompletions:          analytics.TotalCompletions,
		CompletionRate:            analytics.CompletionRate,
		AverageProgressPercentage: analytics.AverageProgressPercentage,
		ActiveStudents:            analytics.ActiveStudents,
		TotalWatchTimeSeconds:     analytics.TotalWatchTimeSeconds,
		AverageWatchTimePerStudent: analytics.AverageWatchTimePerStudent,
		MostPopularLecture:        analytics.MostPopularLecture,
		DropoutRate:               analytics.DropoutRate,
	}
	
	return &pb.GetCourseAnalyticsResponse{
		Analytics: analyticsProto,
	}, nil
}

func (h *ProgressHandler) GetInstructorAnalytics(ctx context.Context, req *pb.GetInstructorAnalyticsRequest) (*pb.GetInstructorAnalyticsResponse, error) {
	// Validate request
	if req.InstructorId == "" {
		return nil, status.Error(codes.InvalidArgument, "instructor_id is required")
	}
	
	// Parse UUID
	instructorID, err := uuid.Parse(req.InstructorId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid instructor_id")
	}
	
	// Get instructor analytics
	analytics, err := h.progressService.GetInstructorAnalytics(instructorID)
	if err != nil {
		h.logger.Errorf("Failed to get instructor analytics: %v", err)
		return nil, status.Error(codes.Internal, "failed to get instructor analytics")
	}
	
	// Convert to proto
	analyticsProto := &pb.InstructorAnalytics{
		InstructorId:          analytics.InstructorID.String(),
		TotalCourses:          analytics.TotalCourses,
		TotalStudents:         analytics.TotalStudents,
		TotalCompletions:      analytics.TotalCompletions,
		AverageCompletionRate: analytics.AverageCompletionRate,
		TotalWatchTimeSeconds: analytics.TotalWatchTimeSeconds,
		AverageCourseRating:   analytics.AverageCourseRating,
		ActiveCourses:         analytics.ActiveCourses,
		BestPerformingCourse:  analytics.BestPerformingCourse,
		TotalRevenueCents:     analytics.TotalRevenueCents,
	}
	
	return &pb.GetInstructorAnalyticsResponse{
		Analytics: analyticsProto,
	}, nil
}