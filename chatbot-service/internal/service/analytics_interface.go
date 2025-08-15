package service

import (
	"context"

	"chatbot-service/internal/model"

	"github.com/google/uuid"
)

// AnalyticsServiceInterface defines the interface for analytics services
type AnalyticsServiceInterface interface {
	GetOverallAnalytics(ctx context.Context, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error)
	GetUserAnalytics(ctx context.Context, userID uuid.UUID, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error)
	GetCourseAnalytics(ctx context.Context, courseID uuid.UUID, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error)
	GetTimeBasedAnalytics(ctx context.Context, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error)
	GetRealTimeMetrics(ctx context.Context) (*model.RealTimeMetrics, error)
	GetSessionMetrics(ctx context.Context, userID uuid.UUID, limit int) ([]*model.SessionMetrics, error)
	GetUsageStats(ctx context.Context, userID uuid.UUID, period string) (*model.UsageStats, error)
	GenerateAnalyticsReport(ctx context.Context, req *model.AnalyticsRequest) (*model.AnalyticsResponse, error)
	CalculateResponseQuality(ctx context.Context, req *model.AnalyticsRequest) (*model.ResponseQuality, error)
}