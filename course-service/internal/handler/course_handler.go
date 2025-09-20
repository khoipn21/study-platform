package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/study-platform/course-service/internal/model"
	"github.com/study-platform/course-service/internal/service"
	pb "github.com/study-platform/course-service/proto"
	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CourseHandler struct {
	pb.UnimplementedCourseServiceServer
	courseService *service.CourseService
	logger        logger.Logger
}

func NewCourseHandler(courseService *service.CourseService, logger logger.Logger) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
		logger:        logger,
	}
}

func (h *CourseHandler) CreateCourse(ctx context.Context, req *pb.CreateCourseRequest) (*pb.CreateCourseResponse, error) {
	h.logger.Info("Request received")
	
	instructorID, err := uuid.Parse(req.InstructorId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid instructor ID: %v", err)
	}
	
	course := &model.Course{
		Title:        req.Title,
		Description:  req.Description,
		InstructorID: instructorID,
		Category:     req.Category,
		Level:        model.CourseLevel(req.Level),
		Price:        req.Price,
		Currency:     req.Currency,
		ThumbnailURL: req.ThumbnailUrl,
		Status:       model.CourseStatus(req.Status),
		Tags:         req.Tags,
	}
	
	err = h.courseService.CreateCourse(ctx, course)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create course: %v", err)
	}
	
	return &pb.CreateCourseResponse{
		Course:  h.courseToProto(course),
		Message: "Course created successfully",
	}, nil
}

func (h *CourseHandler) GetCourse(ctx context.Context, req *pb.GetCourseRequest) (*pb.GetCourseResponse, error) {
	h.logger.Info("Request received")
	
	courseID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid course ID: %v", err)
	}
	
	course, err := h.courseService.GetCourse(ctx, courseID)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.NotFound, "course not found: %v", err)
	}
	
	return &pb.GetCourseResponse{
		Course: h.courseToProto(course),
	}, nil
}

func (h *CourseHandler) UpdateCourse(ctx context.Context, req *pb.UpdateCourseRequest) (*pb.UpdateCourseResponse, error) {
	h.logger.Info("Request received")
	
	courseID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid course ID: %v", err)
	}
	
	course := &model.Course{
		ID:           courseID,
		Title:        req.Title,
		Description:  req.Description,
		Category:     req.Category,
		Level:        model.CourseLevel(req.Level),
		Price:        req.Price,
		Currency:     req.Currency,
		ThumbnailURL: req.ThumbnailUrl,
		Status:       model.CourseStatus(req.Status),
		Tags:         req.Tags,
	}
	
	err = h.courseService.UpdateCourse(ctx, course)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update course: %v", err)
	}
	
	// Get updated course
	updatedCourse, err := h.courseService.GetCourse(ctx, courseID)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get updated course: %v", err)
	}
	
	return &pb.UpdateCourseResponse{
		Course:  h.courseToProto(updatedCourse),
		Message: "Course updated successfully",
	}, nil
}

func (h *CourseHandler) DeleteCourse(ctx context.Context, req *pb.DeleteCourseRequest) (*pb.DeleteCourseResponse, error) {
	h.logger.Info("Request received")
	
	courseID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid course ID: %v", err)
	}
	
	err = h.courseService.DeleteCourse(ctx, courseID)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete course: %v", err)
	}
	
	return &pb.DeleteCourseResponse{
		Message: "Course deleted successfully",
	}, nil
}

func (h *CourseHandler) ListCourses(ctx context.Context, req *pb.ListCoursesRequest) (*pb.ListCoursesResponse, error) {
	h.logger.Info("Request received")
	
	filters := model.CourseFilters{
		Category:     req.Category,
		Level:        req.Level,
		Status:       req.Status,
		InstructorID: req.InstructorId,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}
	
	result, err := h.courseService.ListCourses(ctx, filters)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list courses: %v", err)
	}
	
	courses := make([]*pb.Course, len(result.Courses))
	for i, course := range result.Courses {
		courses[i] = h.courseToProto(&course)
	}
	
	return &pb.ListCoursesResponse{
		Courses:    courses,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
	}, nil
}

func (h *CourseHandler) SearchCourses(ctx context.Context, req *pb.SearchCoursesRequest) (*pb.SearchCoursesResponse, error) {
	h.logger.Info("Request received")
	
	filters := model.CourseFilters{
		Query:     req.Query,
		Category:  req.Category,
		Level:     req.Level,
		MinPrice:  req.MinPrice,
		MaxPrice:  req.MaxPrice,
		MinRating: req.MinRating,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}
	
	result, err := h.courseService.SearchCourses(ctx, filters)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to search courses: %v", err)
	}
	
	courses := make([]*pb.Course, len(result.Courses))
	for i, course := range result.Courses {
		courses[i] = h.courseToProto(&course)
	}
	
	return &pb.SearchCoursesResponse{
		Courses:    courses,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
	}, nil
}

func (h *CourseHandler) CreateLecture(ctx context.Context, req *pb.CreateLectureRequest) (*pb.CreateLectureResponse, error) {
	h.logger.Info("Request received")
	
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid course ID: %v", err)
	}
	
	lecture := &model.Lecture{
		CourseID:        courseID,
		Title:           req.Title,
		Description:     req.Description,
		OrderNumber:     req.OrderNumber,
		DurationMinutes: req.DurationMinutes,
		VideoURL:        req.VideoUrl,
		VideoID:         req.VideoId,
		IsFree:          req.IsFree,
	}
	
	err = h.courseService.CreateLecture(ctx, lecture)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create lecture: %v", err)
	}
	
	return &pb.CreateLectureResponse{
		Lecture: h.lectureToProto(lecture),
		Message: "Lecture created successfully",
	}, nil
}

func (h *CourseHandler) GetLecture(ctx context.Context, req *pb.GetLectureRequest) (*pb.GetLectureResponse, error) {
	h.logger.Info("Request received")
	
	lectureID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid lecture ID: %v", err)
	}
	
	lecture, err := h.courseService.GetLecture(ctx, lectureID)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.NotFound, "lecture not found: %v", err)
	}
	
	return &pb.GetLectureResponse{
		Lecture: h.lectureToProto(lecture),
	}, nil
}

func (h *CourseHandler) UpdateLecture(ctx context.Context, req *pb.UpdateLectureRequest) (*pb.UpdateLectureResponse, error) {
	h.logger.Info("Request received")
	
	lectureID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid lecture ID: %v", err)
	}
	
	lecture := &model.Lecture{
		ID:              lectureID,
		Title:           req.Title,
		Description:     req.Description,
		OrderNumber:     req.OrderNumber,
		DurationMinutes: req.DurationMinutes,
		VideoURL:        req.VideoUrl,
		VideoID:         req.VideoId,
		Status:          model.LectureStatus(req.Status),
		IsFree:          req.IsFree,
	}
	
	err = h.courseService.UpdateLecture(ctx, lecture)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update lecture: %v", err)
	}
	
	// Get updated lecture
	updatedLecture, err := h.courseService.GetLecture(ctx, lectureID)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get updated lecture: %v", err)
	}
	
	return &pb.UpdateLectureResponse{
		Lecture: h.lectureToProto(updatedLecture),
		Message: "Lecture updated successfully",
	}, nil
}

func (h *CourseHandler) DeleteLecture(ctx context.Context, req *pb.DeleteLectureRequest) (*pb.DeleteLectureResponse, error) {
	h.logger.Info("Request received")
	
	lectureID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid lecture ID: %v", err)
	}
	
	err = h.courseService.DeleteLecture(ctx, lectureID)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete lecture: %v", err)
	}
	
	return &pb.DeleteLectureResponse{
		Message: "Lecture deleted successfully",
	}, nil
}

func (h *CourseHandler) ListLectures(ctx context.Context, req *pb.ListLecturesRequest) (*pb.ListLecturesResponse, error) {
	h.logger.Info("Request received")
	
	filters := model.LectureFilters{
		CourseID: req.CourseId,
		Status:   req.Status,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	
	result, err := h.courseService.ListLectures(ctx, filters)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list lectures: %v", err)
	}
	
	lectures := make([]*pb.Lecture, len(result.Lectures))
	for i, lecture := range result.Lectures {
		lectures[i] = h.lectureToProto(&lecture)
	}
	
	return &pb.ListLecturesResponse{
		Lectures:   lectures,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
	}, nil
}

func (h *CourseHandler) EnrollInCourse(ctx context.Context, req *pb.EnrollInCourseRequest) (*pb.EnrollInCourseResponse, error) {
	h.logger.Info("Request received")
	
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}
	
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid course ID: %v", err)
	}
	
	enrollment := &model.Enrollment{
		UserID:   userID,
		CourseID: courseID,
	}
	
	err = h.courseService.EnrollInCourse(ctx, enrollment)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to enroll in course: %v", err)
	}
	
	return &pb.EnrollInCourseResponse{
		Enrollment: h.enrollmentToProto(enrollment),
		Message:    "Successfully enrolled in course",
	}, nil
}

func (h *CourseHandler) GetEnrollment(ctx context.Context, req *pb.GetEnrollmentRequest) (*pb.GetEnrollmentResponse, error) {
	h.logger.Info("Request received")
	
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}
	
	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid course ID: %v", err)
	}
	
	enrollment, err := h.courseService.GetEnrollment(ctx, userID, courseID)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.NotFound, "enrollment not found: %v", err)
	}
	
	return &pb.GetEnrollmentResponse{
		Enrollment: h.enrollmentToProto(enrollment),
	}, nil
}

func (h *CourseHandler) ListEnrollments(ctx context.Context, req *pb.ListEnrollmentsRequest) (*pb.ListEnrollmentsResponse, error) {
	h.logger.Info("Request received")
	
	filters := model.EnrollmentFilters{
		UserID:   req.UserId,
		Status:   req.Status,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	
	result, err := h.courseService.ListEnrollments(ctx, filters)
	if err != nil {
		h.logger.Errorf("Handler error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list enrollments: %v", err)
	}
	
	enrollments := make([]*pb.Enrollment, len(result.Enrollments))
	for i, enrollment := range result.Enrollments {
		enrollments[i] = h.enrollmentToProto(&enrollment)
	}
	
	return &pb.ListEnrollmentsResponse{
		Enrollments: enrollments,
		TotalCount:  result.TotalCount,
		Page:        result.Page,
		PageSize:    result.PageSize,
	}, nil
}

// Helper methods to convert between models and protobuf messages
func (h *CourseHandler) courseToProto(course *model.Course) *pb.Course {
	return &pb.Course{
		Id:              course.ID.String(),
		Title:           course.Title,
		Description:     course.Description,
		InstructorId:    course.InstructorID.String(),
		InstructorName:  course.InstructorName,
		Category:        course.Category,
		Level:           string(course.Level),
		Price:           course.Price,
		Currency:        course.Currency,
		ThumbnailUrl:    course.ThumbnailURL,
		Status:          string(course.Status),
		DurationMinutes: course.DurationMinutes,
		EnrollmentCount: course.EnrollmentCount,
		Rating:          course.Rating,
		RatingCount:     course.RatingCount,
		Tags:            course.Tags,
		CreatedAt:       timestamppb.New(course.CreatedAt),
		UpdatedAt:       timestamppb.New(course.UpdatedAt),
	}
}

func (h *CourseHandler) lectureToProto(lecture *model.Lecture) *pb.Lecture {
	return &pb.Lecture{
		Id:              lecture.ID.String(),
		CourseId:        lecture.CourseID.String(),
		Title:           lecture.Title,
		Description:     lecture.Description,
		OrderNumber:     lecture.OrderNumber,
		DurationMinutes: lecture.DurationMinutes,
		VideoUrl:        lecture.VideoURL,
		VideoId:         lecture.VideoID,
		Status:          string(lecture.Status),
		IsFree:          lecture.IsFree,
		CreatedAt:       timestamppb.New(lecture.CreatedAt),
		UpdatedAt:       timestamppb.New(lecture.UpdatedAt),
	}
}

func (h *CourseHandler) enrollmentToProto(enrollment *model.Enrollment) *pb.Enrollment {
	pbEnrollment := &pb.Enrollment{
		Id:                 enrollment.ID.String(),
		UserId:             enrollment.UserID.String(),
		CourseId:           enrollment.CourseID.String(),
		Status:             string(enrollment.Status),
		ProgressPercentage: enrollment.ProgressPercentage,
		EnrolledAt:         timestamppb.New(enrollment.EnrolledAt),
	}
	
	if enrollment.CompletedAt != nil {
		pbEnrollment.CompletedAt = timestamppb.New(*enrollment.CompletedAt)
	}
	
	if enrollment.LastAccessed != nil {
		pbEnrollment.LastAccessed = timestamppb.New(*enrollment.LastAccessed)
	}
	
	return pbEnrollment
}