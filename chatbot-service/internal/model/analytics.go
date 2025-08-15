package model

import (
	"time"

	"github.com/google/uuid"
)

// ChatAnalytics represents comprehensive chat analytics data
type ChatAnalytics struct {
	TotalSessions       int     `json:"total_sessions"`
	TotalMessages       int     `json:"total_messages"`
	TotalTokensUsed     int     `json:"total_tokens_used"`
	AverageResponseTime float64 `json:"average_response_time_ms"`
	ActiveSessions      int     `json:"active_sessions"`
	UserEngagement      float64 `json:"user_engagement"` // messages per session
}

// UserChatAnalytics represents analytics for a specific user
type UserChatAnalytics struct {
	UserID              uuid.UUID `json:"user_id"`
	TotalSessions       int       `json:"total_sessions"`
	TotalMessages       int       `json:"total_messages"`
	TotalTokensUsed     int       `json:"total_tokens_used"`
	AverageResponseTime float64   `json:"average_response_time_ms"`
	LastActive          time.Time `json:"last_active"`
	ActiveSessions      int       `json:"active_sessions"`
	MessagesPerSession  float64   `json:"messages_per_session"`
}

// CourseAnalytics represents analytics for course-specific chats
type CourseAnalytics struct {
	CourseID            uuid.UUID `json:"course_id"`
	TotalSessions       int       `json:"total_sessions"`
	TotalMessages       int       `json:"total_messages"`
	UniqueUsers         int       `json:"unique_users"`
	AverageResponseTime float64   `json:"average_response_time_ms"`
	TopQuestions        []string  `json:"top_questions"`
}

// TimeBasedAnalytics represents analytics data over a specific time period
type TimeBasedAnalytics struct {
	Period              string    `json:"period"` // daily, weekly, monthly
	Date                time.Time `json:"date"`
	TotalSessions       int       `json:"total_sessions"`
	TotalMessages       int       `json:"total_messages"`
	TotalTokensUsed     int       `json:"total_tokens_used"`
	UniqueUsers         int       `json:"unique_users"`
	AverageResponseTime float64   `json:"average_response_time_ms"`
}

// TopicAnalytics represents analytics for common chat topics
type TopicAnalytics struct {
	Topic               string  `json:"topic"`
	MessageCount        int     `json:"message_count"`
	UserCount           int     `json:"user_count"`
	AverageResponseTime float64 `json:"average_response_time_ms"`
	SentimentScore      float64 `json:"sentiment_score,omitempty"`
}

// AnalyticsRequest represents a request for analytics data
type AnalyticsRequest struct {
	StartDate   *time.Time  `json:"start_date,omitempty"`
	EndDate     *time.Time  `json:"end_date,omitempty"`
	UserID      *uuid.UUID  `json:"user_id,omitempty"`
	CourseID    *uuid.UUID  `json:"course_id,omitempty"`
	Period      string      `json:"period,omitempty"` // daily, weekly, monthly
	Granularity string      `json:"granularity,omitempty"` // hour, day, week, month
}

// AnalyticsResponse wraps analytics data with metadata
type AnalyticsResponse struct {
	Period      string      `json:"period"`
	StartDate   time.Time   `json:"start_date"`
	EndDate     time.Time   `json:"end_date"`
	GeneratedAt time.Time   `json:"generated_at"`
	Data        interface{} `json:"data"`
}

// QuestionPattern represents patterns in user questions
type QuestionPattern struct {
	Pattern     string  `json:"pattern"`
	Count       int     `json:"count"`
	Percentage  float64 `json:"percentage"`
	LastSeen    time.Time `json:"last_seen"`
}

// ResponseQuality represents metrics about response quality
type ResponseQuality struct {
	AverageLength       float64 `json:"average_length"`
	AverageTokens       float64 `json:"average_tokens"`
	AverageResponseTime float64 `json:"average_response_time_ms"`
	SuccessRate         float64 `json:"success_rate"`
	RetryRate           float64 `json:"retry_rate"`
}

// SessionMetrics represents detailed metrics for a chat session
type SessionMetrics struct {
	SessionID           uuid.UUID `json:"session_id"`
	UserID              uuid.UUID `json:"user_id"`
	CourseID            *uuid.UUID `json:"course_id,omitempty"`
	Duration            int64     `json:"duration_seconds"`
	MessageCount        int       `json:"message_count"`
	TokensUsed          int       `json:"tokens_used"`
	AverageResponseTime float64   `json:"average_response_time_ms"`
	StartedAt           time.Time `json:"started_at"`
	EndedAt             time.Time `json:"ended_at"`
}

// RealTimeMetrics represents current system metrics
type RealTimeMetrics struct {
	ActiveSessions      int       `json:"active_sessions"`
	MessagesInLastHour  int       `json:"messages_in_last_hour"`
	TokensInLastHour    int       `json:"tokens_in_last_hour"`
	AverageResponseTime float64   `json:"average_response_time_ms"`
	ErrorRate           float64   `json:"error_rate"`
	LastUpdated         time.Time `json:"last_updated"`
}

// UsageStats represents usage statistics for billing/quota tracking
type UsageStats struct {
	UserID           uuid.UUID `json:"user_id"`
	Period           string    `json:"period"`
	TotalTokens      int       `json:"total_tokens"`
	TotalMessages    int       `json:"total_messages"`
	TotalSessions    int       `json:"total_sessions"`
	EstimatedCost    float64   `json:"estimated_cost"`
	QuotaUsed        float64   `json:"quota_used"`
	QuotaLimit       int       `json:"quota_limit"`
	RemainingQuota   int       `json:"remaining_quota"`
}