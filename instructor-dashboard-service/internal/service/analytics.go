package service

import (
	"fmt"
	"time"
	"instructor-dashboard-service/internal/model"
	"instructor-dashboard-service/internal/repository"

	"github.com/google/uuid"
)

type AnalyticsService struct {
	analyticsRepo *repository.AnalyticsRepository
}

func NewAnalyticsService(analyticsRepo *repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: analyticsRepo,
	}
}

// GetRevenueAnalytics retrieves revenue analytics with default date range if not provided
func (s *AnalyticsService) GetRevenueAnalytics(instructorID uuid.UUID, req *model.AnalyticsRequest) (*model.RevenueAnalytics, error) {
	// Set default date range if not provided
	if req.StartDate == nil {
		defaultStart := time.Now().AddDate(0, -1, 0) // Last month
		req.StartDate = &defaultStart
	}

	if req.EndDate == nil {
		defaultEnd := time.Now()
		req.EndDate = &defaultEnd
	}

	if req.Period == "" {
		req.Period = "monthly"
	}

	analytics, err := s.analyticsRepo.GetRevenueAnalytics(instructorID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue analytics: %w", err)
	}

	return analytics, nil
}

// GetEngagementAnalytics retrieves engagement analytics
func (s *AnalyticsService) GetEngagementAnalytics(instructorID uuid.UUID, req *model.AnalyticsRequest) (*model.EngagementAnalytics, error) {
	// Set default date range if not provided
	if req.StartDate == nil {
		defaultStart := time.Now().AddDate(0, -1, 0) // Last month
		req.StartDate = &defaultStart
	}

	if req.EndDate == nil {
		defaultEnd := time.Now()
		req.EndDate = &defaultEnd
	}

	if req.Period == "" {
		req.Period = "monthly"
	}

	analytics, err := s.analyticsRepo.GetEngagementAnalytics(instructorID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get engagement analytics: %w", err)
	}

	return analytics, nil
}

// GetStudentAnalytics retrieves student analytics
func (s *AnalyticsService) GetStudentAnalytics(instructorID uuid.UUID, req *model.AnalyticsRequest) (*model.StudentAnalytics, error) {
	// Set default date range if not provided
	if req.StartDate == nil {
		defaultStart := time.Now().AddDate(0, -1, 0) // Last month
		req.StartDate = &defaultStart
	}

	if req.EndDate == nil {
		defaultEnd := time.Now()
		req.EndDate = &defaultEnd
	}

	if req.Period == "" {
		req.Period = "monthly"
	}

	analytics, err := s.analyticsRepo.GetStudentAnalytics(instructorID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get student analytics: %w", err)
	}

	return analytics, nil
}

// GetCourseAnalytics retrieves analytics for a specific course
func (s *AnalyticsService) GetCourseAnalytics(instructorID, courseID uuid.UUID, req *model.AnalyticsRequest) (*model.CourseEngagementMetric, error) {
	// Set default date range if not provided
	if req.StartDate == nil {
		defaultStart := time.Now().AddDate(0, -1, 0) // Last month
		req.StartDate = &defaultStart
	}

	if req.EndDate == nil {
		defaultEnd := time.Now()
		req.EndDate = &defaultEnd
	}

	// Get engagement analytics and find the specific course
	engagementAnalytics, err := s.analyticsRepo.GetEngagementAnalytics(instructorID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get course analytics: %w", err)
	}

	// Find the specific course in the engagement breakdown
	for _, courseMetric := range engagementAnalytics.CourseEngagement {
		if courseMetric.CourseID == courseID {
			return &courseMetric, nil
		}
	}

	return nil, fmt.Errorf("course analytics not found")
}

// GetVideoAnalytics retrieves video analytics for an instructor
func (s *AnalyticsService) GetVideoAnalytics(instructorID uuid.UUID, req *model.AnalyticsRequest) (*model.VideoAnalytics, error) {
	// Set default date range if not provided
	if req.StartDate == nil {
		defaultStart := time.Now().AddDate(0, -1, 0) // Last month
		req.StartDate = &defaultStart
	}

	if req.EndDate == nil {
		defaultEnd := time.Now()
		req.EndDate = &defaultEnd
	}

	if req.Period == "" {
		req.Period = "monthly"
	}

	// This would be implemented to get comprehensive video analytics
	// For now, return a basic structure with mock data
	analytics := &model.VideoAnalytics{
		Period:            req.Period,
		StartDate:         *req.StartDate,
		EndDate:           *req.EndDate,
		TotalVideos:       0,
		TotalViews:        0,
		TotalWatchTime:    0,
		AvgWatchTime:      0,
		AvgCompletionRate: 0,
		TopVideos:         []model.VideoPerformance{},
		UnderperformingVideos: []model.VideoPerformance{},
		EngagementHeatmap: []model.EngagementHeatmap{},
		ViewingPatterns:   model.ViewingPatterns{},
	}

	return analytics, nil
}

// GetVideoEngagement retrieves engagement data for a specific video
func (s *AnalyticsService) GetVideoEngagement(instructorID, videoID uuid.UUID, req *model.AnalyticsRequest) (*model.VideoEngagementMetric, error) {
	// Set default date range if not provided
	if req.StartDate == nil {
		defaultStart := time.Now().AddDate(0, -1, 0) // Last month
		req.StartDate = &defaultStart
	}

	if req.EndDate == nil {
		defaultEnd := time.Now()
		req.EndDate = &defaultEnd
	}

	// This would be implemented to get specific video engagement data
	// For now, return a basic structure
	engagement := &model.VideoEngagementMetric{
		VideoID:         videoID,
		VideoTitle:      "Sample Video",
		TotalViews:      0,
		UniqueViewers:   0,
		AvgWatchTime:    0,
		CompletionRate:  0,
		ReplayRate:      0,
		DropOffPoints:   []model.DropOffPoint{},
		EngagementScore: 0,
		AIQuestions:     0,
		BookmarksCreated: 0,
	}

	return engagement, nil
}

// GetAnalyticsSummary provides a high-level analytics summary
func (s *AnalyticsService) GetAnalyticsSummary(instructorID uuid.UUID, period string) (*model.AnalyticsSummary, error) {
	// Set date range based on period
	var startDate, endDate time.Time
	now := time.Now()

	switch period {
	case "daily":
		startDate = now.AddDate(0, 0, -1)
		endDate = now
	case "weekly":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "monthly":
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	case "yearly":
		startDate = now.AddDate(-1, 0, 0)
		endDate = now
	default:
		startDate = now.AddDate(0, -1, 0) // Default to monthly
		endDate = now
	}

	req := &model.AnalyticsRequest{
		Period:    period,
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	// Get all analytics types
	revenueAnalytics, err := s.GetRevenueAnalytics(instructorID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue analytics: %w", err)
	}

	engagementAnalytics, err := s.GetEngagementAnalytics(instructorID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get engagement analytics: %w", err)
	}

	studentAnalytics, err := s.GetStudentAnalytics(instructorID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get student analytics: %w", err)
	}

	// Combine into summary
	summary := &model.AnalyticsSummary{
		Period:              period,
		StartDate:           startDate,
		EndDate:             endDate,
		TotalRevenue:        revenueAnalytics.TotalRevenue,
		RevenueGrowth:       revenueAnalytics.RevenueGrowth,
		TotalStudents:       studentAnalytics.TotalStudents,
		ActiveStudents:      engagementAnalytics.ActiveStudents,
		EngagementRate:      engagementAnalytics.EngagementRate,
		CompletionRate:      engagementAnalytics.CompletionRate,
		StudentSatisfaction: studentAnalytics.StudentSatisfaction,
		ConversionRate:      revenueAnalytics.ConversionRate,
	}

	return summary, nil
}