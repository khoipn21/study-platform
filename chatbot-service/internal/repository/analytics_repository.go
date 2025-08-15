package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"chatbot-service/internal/model"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

// GetOverallAnalytics retrieves overall chat analytics
func (r *AnalyticsRepository) GetOverallAnalytics(ctx context.Context, startDate, endDate time.Time) (*model.ChatAnalytics, error) {
	query := `
		SELECT 
			COUNT(DISTINCT cs.id) as total_sessions,
			COUNT(cm.id) as total_messages,
			COALESCE(SUM(cm.tokens_used), 0) as total_tokens,
			COALESCE(AVG(cm.response_time), 0) as avg_response_time,
			COUNT(DISTINCT CASE WHEN cs.is_active THEN cs.id END) as active_sessions
		FROM chat_sessions cs
		LEFT JOIN chat_messages cm ON cs.id = cm.session_id
		WHERE cs.created_at BETWEEN $1 AND $2`

	analytics := &model.ChatAnalytics{}
	var avgResponseTime sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, startDate, endDate).Scan(
		&analytics.TotalSessions,
		&analytics.TotalMessages,
		&analytics.TotalTokensUsed,
		&avgResponseTime,
		&analytics.ActiveSessions,
	)

	if err != nil {
		return nil, err
	}

	if avgResponseTime.Valid {
		analytics.AverageResponseTime = avgResponseTime.Float64
	}

	// Calculate user engagement (messages per session)
	if analytics.TotalSessions > 0 {
		analytics.UserEngagement = float64(analytics.TotalMessages) / float64(analytics.TotalSessions)
	}

	return analytics, nil
}

// GetUserAnalytics retrieves analytics for a specific user
func (r *AnalyticsRepository) GetUserAnalytics(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*model.UserChatAnalytics, error) {
	query := `
		SELECT 
			COUNT(DISTINCT cs.id) as total_sessions,
			COUNT(cm.id) as total_messages,
			COALESCE(SUM(cm.tokens_used), 0) as total_tokens,
			COALESCE(AVG(cm.response_time), 0) as avg_response_time,
			MAX(cs.updated_at) as last_active,
			COUNT(DISTINCT CASE WHEN cs.is_active THEN cs.id END) as active_sessions
		FROM chat_sessions cs
		LEFT JOIN chat_messages cm ON cs.id = cm.session_id
		WHERE cs.user_id = $1 AND cs.created_at BETWEEN $2 AND $3`

	analytics := &model.UserChatAnalytics{UserID: userID}
	var avgResponseTime sql.NullFloat64
	var lastActive sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID, startDate, endDate).Scan(
		&analytics.TotalSessions,
		&analytics.TotalMessages,
		&analytics.TotalTokensUsed,
		&avgResponseTime,
		&lastActive,
		&analytics.ActiveSessions,
	)

	if err != nil {
		return nil, err
	}

	if avgResponseTime.Valid {
		analytics.AverageResponseTime = avgResponseTime.Float64
	}

	if lastActive.Valid {
		analytics.LastActive = lastActive.Time
	}

	// Calculate messages per session
	if analytics.TotalSessions > 0 {
		analytics.MessagesPerSession = float64(analytics.TotalMessages) / float64(analytics.TotalSessions)
	}

	return analytics, nil
}

// GetCourseAnalytics retrieves analytics for a specific course
func (r *AnalyticsRepository) GetCourseAnalytics(ctx context.Context, courseID uuid.UUID, startDate, endDate time.Time) (*model.CourseAnalytics, error) {
	query := `
		SELECT 
			COUNT(DISTINCT cs.id) as total_sessions,
			COUNT(cm.id) as total_messages,
			COUNT(DISTINCT cs.user_id) as unique_users,
			COALESCE(AVG(cm.response_time), 0) as avg_response_time
		FROM chat_sessions cs
		LEFT JOIN chat_messages cm ON cs.id = cm.session_id
		WHERE cs.course_id = $1 AND cs.created_at BETWEEN $2 AND $3`

	analytics := &model.CourseAnalytics{CourseID: courseID}
	var avgResponseTime sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, courseID, startDate, endDate).Scan(
		&analytics.TotalSessions,
		&analytics.TotalMessages,
		&analytics.UniqueUsers,
		&avgResponseTime,
	)

	if err != nil {
		return nil, err
	}

	if avgResponseTime.Valid {
		analytics.AverageResponseTime = avgResponseTime.Float64
	}

	// Get top questions for the course
	topQuestions, err := r.getTopQuestions(ctx, &courseID, startDate, endDate, 5)
	if err == nil {
		analytics.TopQuestions = topQuestions
	}

	return analytics, nil
}

// GetTimeBasedAnalytics retrieves analytics data over time periods
func (r *AnalyticsRepository) GetTimeBasedAnalytics(ctx context.Context, period string, startDate, endDate time.Time) ([]*model.TimeBasedAnalytics, error) {

	query := `
		SELECT 
			DATE_TRUNC($1, cs.created_at) as period_date,
			COUNT(DISTINCT cs.id) as total_sessions,
			COUNT(cm.id) as total_messages,
			COALESCE(SUM(cm.tokens_used), 0) as total_tokens,
			COUNT(DISTINCT cs.user_id) as unique_users,
			COALESCE(AVG(cm.response_time), 0) as avg_response_time
		FROM chat_sessions cs
		LEFT JOIN chat_messages cm ON cs.id = cm.session_id
		WHERE cs.created_at BETWEEN $2 AND $3
		GROUP BY DATE_TRUNC($1, cs.created_at)
		ORDER BY period_date`

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
	}

	rows, err := r.db.QueryContext(ctx, query, truncateParam, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []*model.TimeBasedAnalytics
	for rows.Next() {
		analytic := &model.TimeBasedAnalytics{Period: period}
		var avgResponseTime sql.NullFloat64

		err := rows.Scan(
			&analytic.Date,
			&analytic.TotalSessions,
			&analytic.TotalMessages,
			&analytic.TotalTokensUsed,
			&analytic.UniqueUsers,
			&avgResponseTime,
		)
		if err != nil {
			return nil, err
		}

		if avgResponseTime.Valid {
			analytic.AverageResponseTime = avgResponseTime.Float64
		}

		analytics = append(analytics, analytic)
	}

	return analytics, nil
}

// GetTopQuestions retrieves the most common questions
func (r *AnalyticsRepository) getTopQuestions(ctx context.Context, courseID *uuid.UUID, startDate, endDate time.Time, limit int) ([]string, error) {
	query := `
		SELECT cm.content, COUNT(*) as question_count
		FROM chat_messages cm
		JOIN chat_sessions cs ON cm.session_id = cs.id
		WHERE cm.role = 'user' 
		AND cs.created_at BETWEEN $1 AND $2`

	params := []interface{}{startDate, endDate}

	if courseID != nil {
		query += " AND cs.course_id = $3"
		params = append(params, *courseID)
	}

	query += `
		GROUP BY cm.content
		ORDER BY question_count DESC
		LIMIT $` + string(rune(len(params)+1))

	params = append(params, limit)

	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []string
	for rows.Next() {
		var content string
		var count int
		err := rows.Scan(&content, &count)
		if err != nil {
			continue
		}

		// Truncate long questions and clean up
		if len(content) > 100 {
			content = content[:97] + "..."
		}
		questions = append(questions, strings.TrimSpace(content))
	}

	return questions, nil
}

// GetRealTimeMetrics retrieves current system metrics
func (r *AnalyticsRepository) GetRealTimeMetrics(ctx context.Context) (*model.RealTimeMetrics, error) {
	now := time.Now()
	oneHourAgo := now.Add(-time.Hour)

	query := `
		SELECT 
			COUNT(DISTINCT CASE WHEN cs.is_active THEN cs.id END) as active_sessions,
			COUNT(CASE WHEN cm.created_at >= $1 THEN cm.id END) as messages_last_hour,
			COALESCE(SUM(CASE WHEN cm.created_at >= $1 THEN cm.tokens_used ELSE 0 END), 0) as tokens_last_hour,
			COALESCE(AVG(CASE WHEN cm.created_at >= $1 THEN cm.response_time END), 0) as avg_response_time
		FROM chat_sessions cs
		LEFT JOIN chat_messages cm ON cs.id = cm.session_id`

	metrics := &model.RealTimeMetrics{LastUpdated: now}
	var avgResponseTime sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, oneHourAgo).Scan(
		&metrics.ActiveSessions,
		&metrics.MessagesInLastHour,
		&metrics.TokensInLastHour,
		&avgResponseTime,
	)

	if err != nil {
		return nil, err
	}

	if avgResponseTime.Valid {
		metrics.AverageResponseTime = avgResponseTime.Float64
	}

	// Calculate error rate (placeholder - would need error tracking)
	metrics.ErrorRate = 0.0

	return metrics, nil
}

// GetSessionMetrics retrieves detailed metrics for specific sessions
func (r *AnalyticsRepository) GetSessionMetrics(ctx context.Context, userID uuid.UUID, limit int) ([]*model.SessionMetrics, error) {
	query := `
		SELECT 
			cs.id,
			cs.user_id,
			cs.course_id,
			EXTRACT(EPOCH FROM (cs.updated_at - cs.created_at))::int as duration,
			COUNT(cm.id) as message_count,
			COALESCE(SUM(cm.tokens_used), 0) as tokens_used,
			COALESCE(AVG(cm.response_time), 0) as avg_response_time,
			cs.created_at,
			cs.updated_at
		FROM chat_sessions cs
		LEFT JOIN chat_messages cm ON cs.id = cm.session_id
		WHERE cs.user_id = $1
		GROUP BY cs.id, cs.user_id, cs.course_id, cs.created_at, cs.updated_at
		ORDER BY cs.created_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*model.SessionMetrics
	for rows.Next() {
		session := &model.SessionMetrics{}
		var courseID sql.NullString
		var avgResponseTime sql.NullFloat64

		err := rows.Scan(
			&session.SessionID,
			&session.UserID,
			&courseID,
			&session.Duration,
			&session.MessageCount,
			&session.TokensUsed,
			&avgResponseTime,
			&session.StartedAt,
			&session.EndedAt,
		)
		if err != nil {
			return nil, err
		}

		if courseID.Valid {
			courseUUID, err := uuid.Parse(courseID.String)
			if err == nil {
				session.CourseID = &courseUUID
			}
		}

		if avgResponseTime.Valid {
			session.AverageResponseTime = avgResponseTime.Float64
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetUsageStats retrieves usage statistics for billing/quota tracking
func (r *AnalyticsRepository) GetUsageStats(ctx context.Context, userID uuid.UUID, period string) (*model.UsageStats, error) {
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
			COUNT(DISTINCT cs.id) as total_sessions,
			COUNT(cm.id) as total_messages,
			COALESCE(SUM(cm.tokens_used), 0) as total_tokens
		FROM chat_sessions cs
		LEFT JOIN chat_messages cm ON cs.id = cm.session_id
		WHERE cs.user_id = $1 AND cs.created_at BETWEEN $2 AND $3`

	stats := &model.UsageStats{
		UserID: userID,
		Period: period,
	}

	err := r.db.QueryRowContext(ctx, query, userID, startDate, now).Scan(
		&stats.TotalSessions,
		&stats.TotalMessages,
		&stats.TotalTokens,
	)

	if err != nil {
		return nil, err
	}

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