package service

import (
	"context"
	"time"

	"chatbot-service/internal/model"
	"chatbot-service/internal/repository"

	"github.com/google/uuid"
)

type AnalyticsService struct {
	analyticsRepo *repository.AnalyticsRepository
	chatRepo      *repository.ChatRepository
}

func NewAnalyticsService(analyticsRepo *repository.AnalyticsRepository, chatRepo *repository.ChatRepository) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: analyticsRepo,
		chatRepo:      chatRepo,
	}
}

// GetOverallAnalytics retrieves overall system analytics
func (s *AnalyticsService) GetOverallAnalytics(ctx context.Context, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error) {
	startDate, endDate := s.getDateRange(req)

	analytics, err := s.analyticsRepo.GetOverallAnalytics(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &model.AnalyticsResponse{
		Period:      s.getPeriodString(req),
		StartDate:   startDate,
		EndDate:     endDate,
		GeneratedAt: time.Now(),
		Data:        analytics,
	}, nil
}

// GetUserAnalytics retrieves analytics for a specific user
func (s *AnalyticsService) GetUserAnalytics(ctx context.Context, userID uuid.UUID, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error) {
	startDate, endDate := s.getDateRange(req)

	analytics, err := s.analyticsRepo.GetUserAnalytics(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &model.AnalyticsResponse{
		Period:      s.getPeriodString(req),
		StartDate:   startDate,
		EndDate:     endDate,
		GeneratedAt: time.Now(),
		Data:        analytics,
	}, nil
}

// GetCourseAnalytics retrieves analytics for a specific course
func (s *AnalyticsService) GetCourseAnalytics(ctx context.Context, courseID uuid.UUID, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error) {
	startDate, endDate := s.getDateRange(req)

	analytics, err := s.analyticsRepo.GetCourseAnalytics(ctx, courseID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &model.AnalyticsResponse{
		Period:      s.getPeriodString(req),
		StartDate:   startDate,
		EndDate:     endDate,
		GeneratedAt: time.Now(),
		Data:        analytics,
	}, nil
}

// GetTimeBasedAnalytics retrieves analytics data over time
func (s *AnalyticsService) GetTimeBasedAnalytics(ctx context.Context, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error) {
	startDate, endDate := s.getDateRange(req)
	period := req.Period
	if period == "" {
		period = "daily"
	}

	analytics, err := s.analyticsRepo.GetTimeBasedAnalytics(ctx, period, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &model.AnalyticsResponse{
		Period:      period,
		StartDate:   startDate,
		EndDate:     endDate,
		GeneratedAt: time.Now(),
		Data:        analytics,
	}, nil
}

// GetRealTimeMetrics retrieves current system metrics
func (s *AnalyticsService) GetRealTimeMetrics(ctx context.Context) (*model.RealTimeMetrics, error) {
	return s.analyticsRepo.GetRealTimeMetrics(ctx)
}

// GetSessionMetrics retrieves detailed session metrics for a user
func (s *AnalyticsService) GetSessionMetrics(ctx context.Context, userID uuid.UUID, limit int) ([]*model.SessionMetrics, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.analyticsRepo.GetSessionMetrics(ctx, userID, limit)
}

// GetUsageStats retrieves usage statistics for billing/quota tracking
func (s *AnalyticsService) GetUsageStats(ctx context.Context, userID uuid.UUID, period string) (*model.UsageStats, error) {
	return s.analyticsRepo.GetUsageStats(ctx, userID, period)
}

// GenerateAnalyticsReport generates a comprehensive analytics report
func (s *AnalyticsService) GenerateAnalyticsReport(ctx context.Context, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error) {
	startDate, endDate := s.getDateRange(req)

	// Get overall analytics
	overall, err := s.analyticsRepo.GetOverallAnalytics(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get time-based analytics
	period := req.Period
	if period == "" {
		period = "daily"
	}
	timeBased, err := s.analyticsRepo.GetTimeBasedAnalytics(ctx, period, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get real-time metrics
	realTime, err := s.analyticsRepo.GetRealTimeMetrics(ctx)
	if err != nil {
		return nil, err
	}

	// Combine all data into a comprehensive report
	report := map[string]interface{}{
		"overall_analytics":   overall,
		"time_based_data":     timeBased,
		"real_time_metrics":   realTime,
		"report_generated_at": time.Now(),
		"period_analyzed":     period,
	}

	// Add user-specific or course-specific data if requested
	if req.UserID != nil {
		userAnalytics, err := s.analyticsRepo.GetUserAnalytics(ctx, *req.UserID, startDate, endDate)
		if err == nil {
			report["user_analytics"] = userAnalytics
		}

		usageStats, err := s.analyticsRepo.GetUsageStats(ctx, *req.UserID, period)
		if err == nil {
			report["usage_stats"] = usageStats
		}
	}

	if req.CourseID != nil {
		courseAnalytics, err := s.analyticsRepo.GetCourseAnalytics(ctx, *req.CourseID, startDate, endDate)
		if err == nil {
			report["course_analytics"] = courseAnalytics
		}
	}

	return &model.AnalyticsResponse{
		Period:      period,
		StartDate:   startDate,
		EndDate:     endDate,
		GeneratedAt: time.Now(),
		Data:        report,
	}, nil
}

// CalculateResponseQuality calculates response quality metrics
func (s *AnalyticsService) CalculateResponseQuality(ctx context.Context, req *model.AnalyticsRequest) (*model.ResponseQuality, error) {
	// This would be implemented with more sophisticated analysis
	// For now, we'll return basic metrics from the database

	// Placeholder implementation - this would need more complex queries
	// to calculate actual response quality metrics
	quality := &model.ResponseQuality{
		AverageLength:       150.0, // characters
		AverageTokens:       45.0,  // tokens
		AverageResponseTime: 1200.0, // milliseconds
		SuccessRate:         0.95,   // 95%
		RetryRate:          0.05,   // 5%
	}

	return quality, nil
}

// TrackUserEngagement tracks user engagement patterns
func (s *AnalyticsService) TrackUserEngagement(ctx context.Context, userID uuid.UUID) error {
	// This could be used to track user engagement events
	// For now, it's a placeholder for future implementation
	
	// Could track:
	// - Session start/end times
	// - Message frequency
	// - Response satisfaction
	// - Feature usage
	
	return nil
}

// Helper methods

func (s *AnalyticsService) getDateRange(req *model.AnalyticsRequest) (time.Time, time.Time) {
	now := time.Now()
	var startDate, endDate time.Time

	if req.StartDate != nil && req.EndDate != nil {
		startDate = *req.StartDate
		endDate = *req.EndDate
	} else {
		// Default to last 30 days
		endDate = now
		startDate = now.AddDate(0, 0, -30)
	}

	return startDate, endDate
}

func (s *AnalyticsService) getPeriodString(req *model.AnalyticsRequest) string {
	if req.Period != "" {
		return req.Period
	}
	
	startDate, endDate := s.getDateRange(req)
	days := int(endDate.Sub(startDate).Hours() / 24)
	
	switch {
	case days <= 7:
		return "daily"
	case days <= 60:
		return "weekly"
	default:
		return "monthly"
	}
}