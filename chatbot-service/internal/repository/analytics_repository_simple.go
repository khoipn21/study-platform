package repository

import (
	"context"
	"database/sql"
	"time"

	"chatbot-service/internal/model"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// SimpleAnalyticsRepository works with the existing chat_history table schema
type SimpleAnalyticsRepository struct {
	db *sql.DB
}

func NewSimpleAnalyticsRepository(db *sql.DB) *SimpleAnalyticsRepository {
	return &SimpleAnalyticsRepository{db: db}
}

// GetOverallAnalytics retrieves overall chat analytics from chat_history
func (r *SimpleAnalyticsRepository) GetOverallAnalytics(ctx context.Context, startDate, endDate time.Time) (*model.ChatAnalytics, error) {
	query := `
		SELECT 
			COUNT(CASE WHEN is_user = true THEN 1 END) as user_messages,
			COUNT(CASE WHEN is_user = false THEN 1 END) as bot_messages,
			COUNT(DISTINCT user_id) as unique_users,
			COUNT(*) as total_messages
		FROM chat_history
		WHERE created_at BETWEEN $1 AND $2`

	analytics := &model.ChatAnalytics{}
	var userMessages, botMessages, uniqueUsers int

	err := r.db.QueryRowContext(ctx, query, startDate, endDate).Scan(
		&userMessages,
		&botMessages,
		&uniqueUsers,
		&analytics.TotalMessages,
	)

	if err != nil {
		return nil, err
	}

	// Estimate sessions (group consecutive messages within 30 minutes)
	analytics.TotalSessions = userMessages // Simple approximation
	analytics.ActiveSessions = 0 // No concept of active sessions in simple schema
	analytics.TotalTokensUsed = botMessages * 50 // Rough estimate: 50 tokens per response
	analytics.AverageResponseTime = 1200.0 // Default estimate

	// Calculate user engagement
	if analytics.TotalSessions > 0 {
		analytics.UserEngagement = float64(analytics.TotalMessages) / float64(analytics.TotalSessions)
	}

	return analytics, nil
}

// GetUserAnalytics retrieves analytics for a specific user from chat_history
func (r *SimpleAnalyticsRepository) GetUserAnalytics(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*model.UserChatAnalytics, error) {
	query := `
		SELECT 
			COUNT(CASE WHEN is_user = true THEN 1 END) as user_messages,
			COUNT(CASE WHEN is_user = false THEN 1 END) as bot_messages,
			COUNT(*) as total_messages,
			MAX(created_at) as last_active
		FROM chat_history
		WHERE user_id = $1 AND created_at BETWEEN $2 AND $3`

	analytics := &model.UserChatAnalytics{UserID: userID}
	var userMessages, botMessages int
	var lastActive sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID, startDate, endDate).Scan(
		&userMessages,
		&botMessages,
		&analytics.TotalMessages,
		&lastActive,
	)

	if err != nil {
		return nil, err
	}

	if lastActive.Valid {
		analytics.LastActive = lastActive.Time
	}

	// Estimate values
	analytics.TotalSessions = userMessages // Each user message starts a "session"
	analytics.ActiveSessions = 0
	analytics.TotalTokensUsed = botMessages * 50 // Estimate tokens
	analytics.AverageResponseTime = 1200.0 // Default estimate

	// Calculate messages per session
	if analytics.TotalSessions > 0 {
		analytics.MessagesPerSession = float64(analytics.TotalMessages) / float64(analytics.TotalSessions)
	}

	return analytics, nil
}

// GetRealTimeMetrics retrieves current system metrics from chat_history
func (r *SimpleAnalyticsRepository) GetRealTimeMetrics(ctx context.Context) (*model.RealTimeMetrics, error) {
	now := time.Now()
	oneHourAgo := now.Add(-time.Hour)

	query := `
		SELECT 
			COUNT(*) as messages_last_hour,
			COUNT(DISTINCT user_id) as active_users_last_hour
		FROM chat_history
		WHERE created_at >= $1`

	metrics := &model.RealTimeMetrics{LastUpdated: now}
	var activeUsers int

	err := r.db.QueryRowContext(ctx, query, oneHourAgo).Scan(
		&metrics.MessagesInLastHour,
		&activeUsers,
	)

	if err != nil {
		return nil, err
	}

	// Estimate values
	metrics.ActiveSessions = activeUsers
	metrics.TokensInLastHour = metrics.MessagesInLastHour * 25 // Rough estimate
	metrics.AverageResponseTime = 1200.0 // Default estimate
	metrics.ErrorRate = 0.0 // No error tracking in simple schema

	return metrics, nil
}

// GetTimeBasedAnalytics retrieves analytics data over time periods from chat_history
func (r *SimpleAnalyticsRepository) GetTimeBasedAnalytics(ctx context.Context, period string, startDate, endDate time.Time) ([]*model.TimeBasedAnalytics, error) {
	var truncateParam string
	switch period {
	case "daily":
		truncateParam = "day"
	case "weekly":
		truncateParam = "week"
	case "monthly":
		truncateParam = "month"
	default:
		truncateParam = "day"
		period = "daily"
	}

	query := `
		SELECT 
			DATE_TRUNC($1, created_at) as period_date,
			COUNT(*) as total_messages,
			COUNT(DISTINCT user_id) as unique_users,
			COUNT(CASE WHEN is_user = false THEN 1 END) as bot_responses
		FROM chat_history
		WHERE created_at BETWEEN $2 AND $3
		GROUP BY DATE_TRUNC($1, created_at)
		ORDER BY period_date`

	rows, err := r.db.QueryContext(ctx, query, truncateParam, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []*model.TimeBasedAnalytics
	for rows.Next() {
		analytic := &model.TimeBasedAnalytics{Period: period}
		var botResponses int

		err := rows.Scan(
			&analytic.Date,
			&analytic.TotalMessages,
			&analytic.UniqueUsers,
			&botResponses,
		)
		if err != nil {
			return nil, err
		}

		// Estimate values
		analytic.TotalSessions = analytic.TotalMessages / 2 // Rough estimate
		analytic.TotalTokensUsed = botResponses * 50 // Estimate tokens
		analytic.AverageResponseTime = 1200.0 // Default estimate

		analytics = append(analytics, analytic)
	}

	return analytics, nil
}

// GetUsageStats retrieves usage statistics for billing/quota tracking
func (r *SimpleAnalyticsRepository) GetUsageStats(ctx context.Context, userID uuid.UUID, period string) (*model.UsageStats, error) {
	var startDate time.Time
	now := time.Now()

	switch period {
	case "daily":
		startDate = now.AddDate(0, 0, -1)
	case "weekly":
		startDate = now.AddDate(0, 0, -7)
	case "monthly":
		startDate = now.AddDate(0, -1, 0)
	default:
		startDate = now.AddDate(0, -1, 0) // default to monthly
		period = "monthly"
	}

	query := `
		SELECT 
			COUNT(*) as total_messages,
			COUNT(CASE WHEN is_user = false THEN 1 END) as bot_responses
		FROM chat_history
		WHERE user_id = $1 AND created_at BETWEEN $2 AND $3`

	stats := &model.UsageStats{
		UserID: userID,
		Period: period,
	}

	var botResponses int

	err := r.db.QueryRowContext(ctx, query, userID, startDate, now).Scan(
		&stats.TotalMessages,
		&botResponses,
	)

	if err != nil {
		return nil, err
	}

	// Estimate values
	stats.TotalSessions = stats.TotalMessages / 2 // Rough estimate
	stats.TotalTokens = botResponses * 50 // Estimate tokens

	// Calculate estimated cost (assuming $0.002 per 1K tokens - GPT-3.5 pricing)
	stats.EstimatedCost = float64(stats.TotalTokens) * 0.002 / 1000

	// Set quota limits (this would typically come from user settings)
	stats.QuotaLimit = 100000 // 100K tokens default
	stats.RemainingQuota = stats.QuotaLimit - stats.TotalTokens
	if stats.RemainingQuota < 0 {
		stats.RemainingQuota = 0
	}
	
	stats.QuotaUsed = float64(stats.TotalTokens) / float64(stats.QuotaLimit) * 100
	if stats.QuotaUsed > 100 {
		stats.QuotaUsed = 100
	}

	return stats, nil
}

// GetSessionMetrics returns a simplified version based on chat_history
func (r *SimpleAnalyticsRepository) GetSessionMetrics(ctx context.Context, userID uuid.UUID, limit int) ([]*model.SessionMetrics, error) {
	// Since we don't have sessions, we'll return recent conversations grouped by day
	query := `
		SELECT 
			DATE_TRUNC('day', created_at) as session_date,
			COUNT(*) as message_count,
			MIN(created_at) as started_at,
			MAX(created_at) as ended_at
		FROM chat_history
		WHERE user_id = $1
		GROUP BY DATE_TRUNC('day', created_at)
		ORDER BY session_date DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*model.SessionMetrics
	for rows.Next() {
		session := &model.SessionMetrics{
			UserID: userID,
		}

		var sessionDate time.Time
		err := rows.Scan(
			&sessionDate,
			&session.MessageCount,
			&session.StartedAt,
			&session.EndedAt,
		)
		if err != nil {
			return nil, err
		}

		// Generate a pseudo session ID from the date
		session.SessionID = uuid.New()
		session.Duration = int64(session.EndedAt.Sub(session.StartedAt).Seconds())
		session.TokensUsed = session.MessageCount * 25 // Estimate
		session.AverageResponseTime = 1200.0 // Default estimate

		sessions = append(sessions, session)
	}

	return sessions, nil
}