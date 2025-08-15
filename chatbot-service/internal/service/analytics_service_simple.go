package service

import (
	"context"
	"time"

	"chatbot-service/internal/model"
	"chatbot-service/internal/repository"

	"github.com/google/uuid"
)

type SimpleAnalyticsService struct {
	analyticsRepo *repository.SimpleAnalyticsRepository
	chatRepo      *repository.SimpleChatRepository
}

func NewSimpleAnalyticsService(analyticsRepo *repository.SimpleAnalyticsRepository, chatRepo *repository.SimpleChatRepository) *SimpleAnalyticsService {
	return &SimpleAnalyticsService{
		analyticsRepo: analyticsRepo,
		chatRepo:      chatRepo,
	}
}

// GetOverallAnalytics retrieves overall system analytics
func (s *SimpleAnalyticsService) GetOverallAnalytics(ctx context.Context, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error) {
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
func (s *SimpleAnalyticsService) GetUserAnalytics(ctx context.Context, userID uuid.UUID, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error) {
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

// GetTimeBasedAnalytics retrieves analytics data over time
func (s *SimpleAnalyticsService) GetTimeBasedAnalytics(ctx context.Context, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error) {
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
func (s *SimpleAnalyticsService) GetRealTimeMetrics(ctx context.Context) (*model.RealTimeMetrics, error) {
	return s.analyticsRepo.GetRealTimeMetrics(ctx)
}

// GetSessionMetrics retrieves detailed session metrics for a user
func (s *SimpleAnalyticsService) GetSessionMetrics(ctx context.Context, userID uuid.UUID, limit int) ([]*model.SessionMetrics, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.analyticsRepo.GetSessionMetrics(ctx, userID, limit)
}

// GetUsageStats retrieves usage statistics for billing/quota tracking
func (s *SimpleAnalyticsService) GetUsageStats(ctx context.Context, userID uuid.UUID, period string) (*model.UsageStats, error) {
	return s.analyticsRepo.GetUsageStats(ctx, userID, period)
}

// Course analytics placeholder
func (s *SimpleAnalyticsService) GetCourseAnalytics(ctx context.Context, courseID uuid.UUID, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error) {
	// Placeholder - not supported with simple schema
	analytics := &model.CourseAnalytics{
		CourseID:            courseID,
		TotalSessions:       0,
		TotalMessages:       0,
		UniqueUsers:         0,
		AverageResponseTime: 0,
		TopQuestions:        []string{},
	}

	startDate, endDate := s.getDateRange(req)
	return &model.AnalyticsResponse{
		Period:      s.getPeriodString(req),
		StartDate:   startDate,
		EndDate:     endDate,
		GeneratedAt: time.Now(),
		Data:        analytics,
	}, nil
}

// GenerateAnalyticsReport generates a comprehensive analytics report
func (s *SimpleAnalyticsService) GenerateAnalyticsReport(ctx context.Context, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error) {
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

	// Add user-specific data if requested
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

	return &model.AnalyticsResponse{
		Period:      period,
		StartDate:   startDate,
		EndDate:     endDate,
		GeneratedAt: time.Now(),
		Data:        report,
	}, nil
}

// CalculateResponseQuality calculates response quality metrics
func (s *SimpleAnalyticsService) CalculateResponseQuality(ctx context.Context, req *model.AnalyticsRequest) (*model.ResponseQuality, error) {
	// Placeholder implementation
	quality := &model.ResponseQuality{
		AverageLength:       150.0, // characters
		AverageTokens:       45.0,  // tokens
		AverageResponseTime: 1200.0, // milliseconds
		SuccessRate:         0.95,   // 95%
		RetryRate:          0.05,   // 5%
	}

	return quality, nil
}

// Helper methods
func (s *SimpleAnalyticsService) getDateRange(req *model.AnalyticsRequest) (time.Time, time.Time) {
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

func (s *SimpleAnalyticsService) getPeriodString(req *model.AnalyticsRequest) string {
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