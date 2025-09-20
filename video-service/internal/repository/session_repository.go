package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"video-service/internal/model"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// CreateViewingSession creates a new viewing session
func (sr *SessionRepository) CreateViewingSession(session *model.ViewingSession) error {
	query := `
		INSERT INTO viewing_sessions (
			id, session_id, user_id, video_id, started_at, last_heartbeat,
			current_time_seconds, current_quality, total_watch_time_seconds,
			completed, user_agent, ip_address
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at
	`

	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}

	if session.StartedAt.IsZero() {
		session.StartedAt = time.Now()
	}

	if session.LastHeartbeat.IsZero() {
		session.LastHeartbeat = time.Now()
	}

	err := sr.db.QueryRow(
		query, session.ID, session.SessionID, session.UserID, session.VideoID,
		session.StartedAt, session.LastHeartbeat, session.CurrentTimeSeconds,
		session.CurrentQuality, session.TotalWatchTimeSeconds, session.Completed,
		session.UserAgent, session.IPAddress,
	).Scan(&session.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create viewing session: %w", err)
	}

	return nil
}

// GetViewingSessionBySessionID retrieves a viewing session by session ID
func (sr *SessionRepository) GetViewingSessionBySessionID(sessionID string) (*model.ViewingSession, error) {
	query := `
		SELECT id, session_id, user_id, video_id, started_at, last_heartbeat,
			   current_time_seconds, current_quality, total_watch_time_seconds,
			   completed, user_agent, ip_address, created_at
		FROM viewing_sessions 
		WHERE session_id = $1
	`

	session := &model.ViewingSession{}

	err := sr.db.QueryRow(query, sessionID).Scan(
		&session.ID, &session.SessionID, &session.UserID, &session.VideoID,
		&session.StartedAt, &session.LastHeartbeat, &session.CurrentTimeSeconds,
		&session.CurrentQuality, &session.TotalWatchTimeSeconds, &session.Completed,
		&session.UserAgent, &session.IPAddress, &session.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("viewing session not found")
		}
		return nil, fmt.Errorf("failed to get viewing session: %w", err)
	}

	return session, nil
}

// UpdateSessionProgress updates the progress of a viewing session
func (sr *SessionRepository) UpdateSessionProgress(sessionID string, currentTime int, quality string, watchTime int) error {
	query := `
		UPDATE viewing_sessions 
		SET current_time_seconds = $1, current_quality = $2, 
			total_watch_time_seconds = $3, last_heartbeat = NOW()
		WHERE session_id = $4
	`

	result, err := sr.db.Exec(query, currentTime, quality, watchTime, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update session progress: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("viewing session not found")
	}

	return nil
}

// UpdateSessionHeartbeat updates the last heartbeat of a session
func (sr *SessionRepository) UpdateSessionHeartbeat(sessionID string) error {
	query := `
		UPDATE viewing_sessions 
		SET last_heartbeat = NOW()
		WHERE session_id = $1
	`

	result, err := sr.db.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update session heartbeat: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("viewing session not found")
	}

	return nil
}

// MarkSessionCompleted marks a viewing session as completed
func (sr *SessionRepository) MarkSessionCompleted(sessionID string) error {
	query := `
		UPDATE viewing_sessions 
		SET completed = true, last_heartbeat = NOW()
		WHERE session_id = $1
	`

	result, err := sr.db.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to mark session completed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("viewing session not found")
	}

	return nil
}

// GetUserViewingSessions retrieves all viewing sessions for a user
func (sr *SessionRepository) GetUserViewingSessions(userID uuid.UUID, limit, offset int) ([]*model.ViewingSession, error) {
	query := `
		SELECT id, session_id, user_id, video_id, started_at, last_heartbeat,
			   current_time_seconds, current_quality, total_watch_time_seconds,
			   completed, user_agent, ip_address, created_at
		FROM viewing_sessions 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := sr.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query viewing sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*model.ViewingSession
	for rows.Next() {
		session := &model.ViewingSession{}

		err := rows.Scan(
			&session.ID, &session.SessionID, &session.UserID, &session.VideoID,
			&session.StartedAt, &session.LastHeartbeat, &session.CurrentTimeSeconds,
			&session.CurrentQuality, &session.TotalWatchTimeSeconds, &session.Completed,
			&session.UserAgent, &session.IPAddress, &session.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan viewing session: %w", err)
		}

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return sessions, nil
}

// GetVideoViewingSessions retrieves all viewing sessions for a video
func (sr *SessionRepository) GetVideoViewingSessions(videoID uuid.UUID, limit, offset int) ([]*model.ViewingSession, error) {
	query := `
		SELECT id, session_id, user_id, video_id, started_at, last_heartbeat,
			   current_time_seconds, current_quality, total_watch_time_seconds,
			   completed, user_agent, ip_address, created_at
		FROM viewing_sessions 
		WHERE video_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := sr.db.Query(query, videoID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query viewing sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*model.ViewingSession
	for rows.Next() {
		session := &model.ViewingSession{}

		err := rows.Scan(
			&session.ID, &session.SessionID, &session.UserID, &session.VideoID,
			&session.StartedAt, &session.LastHeartbeat, &session.CurrentTimeSeconds,
			&session.CurrentQuality, &session.TotalWatchTimeSeconds, &session.Completed,
			&session.UserAgent, &session.IPAddress, &session.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan viewing session: %w", err)
		}

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return sessions, nil
}

// CreateNetworkMetrics creates a new network metrics record
func (sr *SessionRepository) CreateNetworkMetrics(metrics *model.NetworkMetrics) error {
	query := `
		INSERT INTO network_metrics (
			id, session_id, user_id, timestamp, bandwidth_mbps, latency_ms,
			packet_loss_percent, connection_type, quality_score, 
			recommended_quality, buffer_health_seconds
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING created_at
	`

	if metrics.ID == uuid.Nil {
		metrics.ID = uuid.New()
	}

	if metrics.Timestamp.IsZero() {
		metrics.Timestamp = time.Now()
	}

	err := sr.db.QueryRow(
		query, metrics.ID, metrics.SessionID, metrics.UserID, metrics.Timestamp,
		metrics.BandwidthMbps, metrics.LatencyMs, metrics.PacketLossPercent,
		metrics.ConnectionType, metrics.QualityScore, metrics.RecommendedQuality,
		metrics.BufferHealthSeconds,
	).Scan(&metrics.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create network metrics: %w", err)
	}

	return nil
}

// GetNetworkMetricsBySession retrieves network metrics for a session
func (sr *SessionRepository) GetNetworkMetricsBySession(sessionID string, limit, offset int) ([]*model.NetworkMetrics, error) {
	query := `
		SELECT id, session_id, user_id, timestamp, bandwidth_mbps, latency_ms,
			   packet_loss_percent, connection_type, quality_score, 
			   recommended_quality, buffer_health_seconds, created_at
		FROM network_metrics 
		WHERE session_id = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := sr.db.Query(query, sessionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query network metrics: %w", err)
	}
	defer rows.Close()

	var metrics []*model.NetworkMetrics
	for rows.Next() {
		metric := &model.NetworkMetrics{}

		err := rows.Scan(
			&metric.ID, &metric.SessionID, &metric.UserID, &metric.Timestamp,
			&metric.BandwidthMbps, &metric.LatencyMs, &metric.PacketLossPercent,
			&metric.ConnectionType, &metric.QualityScore, &metric.RecommendedQuality,
			&metric.BufferHealthSeconds, &metric.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan network metrics: %w", err)
		}

		metrics = append(metrics, metric)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return metrics, nil
}

// GetActiveSessionsByUser retrieves active sessions for a user (sessions with recent heartbeat)
func (sr *SessionRepository) GetActiveSessionsByUser(userID uuid.UUID, heartbeatThreshold time.Duration) ([]*model.ViewingSession, error) {
	query := `
		SELECT id, session_id, user_id, video_id, started_at, last_heartbeat,
			   current_time_seconds, current_quality, total_watch_time_seconds,
			   completed, user_agent, ip_address, created_at
		FROM viewing_sessions 
		WHERE user_id = $1 AND last_heartbeat > $2 AND completed = false
		ORDER BY last_heartbeat DESC
	`

	thresholdTime := time.Now().Add(-heartbeatThreshold)

	rows, err := sr.db.Query(query, userID, thresholdTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query active sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*model.ViewingSession
	for rows.Next() {
		session := &model.ViewingSession{}

		err := rows.Scan(
			&session.ID, &session.SessionID, &session.UserID, &session.VideoID,
			&session.StartedAt, &session.LastHeartbeat, &session.CurrentTimeSeconds,
			&session.CurrentQuality, &session.TotalWatchTimeSeconds, &session.Completed,
			&session.UserAgent, &session.IPAddress, &session.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan viewing session: %w", err)
		}

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return sessions, nil
}

// CleanupStaleSessions removes stale sessions (sessions with old heartbeat)
func (sr *SessionRepository) CleanupStaleSessions(heartbeatThreshold time.Duration) (int, error) {
	query := `
		DELETE FROM viewing_sessions 
		WHERE last_heartbeat < $1 AND completed = false
	`

	thresholdTime := time.Now().Add(-heartbeatThreshold)

	result, err := sr.db.Exec(query, thresholdTime)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup stale sessions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

// GetSessionWatchTime calculates total watch time for a session
func (sr *SessionRepository) GetSessionWatchTime(sessionID string) (int, error) {
	query := `
		SELECT COALESCE(total_watch_time_seconds, 0)
		FROM viewing_sessions 
		WHERE session_id = $1
	`

	var watchTime int
	err := sr.db.QueryRow(query, sessionID).Scan(&watchTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("viewing session not found")
		}
		return 0, fmt.Errorf("failed to get session watch time: %w", err)
	}

	return watchTime, nil
}

// GetUserVideoProgress gets user's progress for a specific video
func (sr *SessionRepository) GetUserVideoProgress(userID, videoID uuid.UUID) (*model.ViewingSession, error) {
	query := `
		SELECT id, session_id, user_id, video_id, started_at, last_heartbeat,
			   current_time_seconds, current_quality, total_watch_time_seconds,
			   completed, user_agent, ip_address, created_at
		FROM viewing_sessions 
		WHERE user_id = $1 AND video_id = $2
		ORDER BY last_heartbeat DESC
		LIMIT 1
	`

	session := &model.ViewingSession{}

	err := sr.db.QueryRow(query, userID, videoID).Scan(
		&session.ID, &session.SessionID, &session.UserID, &session.VideoID,
		&session.StartedAt, &session.LastHeartbeat, &session.CurrentTimeSeconds,
		&session.CurrentQuality, &session.TotalWatchTimeSeconds, &session.Completed,
		&session.UserAgent, &session.IPAddress, &session.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No progress found, not an error
		}
		return nil, fmt.Errorf("failed to get user video progress: %w", err)
	}

	return session, nil
}

// GetActiveSessionByUserVideo retrieves active session for a user and video combination
func (sr *SessionRepository) GetActiveSessionByUserVideo(userID, videoID uuid.UUID) (*model.ViewingSession, error) {
	query := `
		SELECT id, session_id, user_id, video_id, started_at, last_heartbeat,
			   current_time_seconds, current_quality, total_watch_time_seconds,
			   completed, user_agent, ip_address, created_at
		FROM viewing_sessions
		WHERE user_id = $1 AND video_id = $2 AND completed = false
			  AND last_heartbeat > NOW() - INTERVAL '1 hour'
		ORDER BY last_heartbeat DESC
		LIMIT 1
	`

	session := &model.ViewingSession{}

	err := sr.db.QueryRow(query, userID, videoID).Scan(
		&session.ID, &session.SessionID, &session.UserID, &session.VideoID,
		&session.StartedAt, &session.LastHeartbeat, &session.CurrentTimeSeconds,
		&session.CurrentQuality, &session.TotalWatchTimeSeconds, &session.Completed,
		&session.UserAgent, &session.IPAddress, &session.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get active session: %w", err)
	}

	return session, nil
}