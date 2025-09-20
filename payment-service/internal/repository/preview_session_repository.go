package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PreviewSessionRepository struct {
	db *sql.DB
}

func NewPreviewSessionRepository(db *sql.DB) *PreviewSessionRepository {
	return &PreviewSessionRepository{
		db: db,
	}
}

type PreviewSession struct {
	ID                     string     `json:"id"`
	UserID                 string     `json:"user_id"`
	CourseID               string     `json:"course_id"`
	LectureID              string     `json:"lecture_id"`
	SessionStartedAt       time.Time  `json:"session_started_at"`
	SessionDurationSeconds int        `json:"session_duration_seconds"`
	PreviewLimitSeconds    int        `json:"preview_limit_seconds"`
	PreviewExhausted       bool       `json:"preview_exhausted"`
	IPAddress              string     `json:"ip_address"`
	LastAccessedAt         *time.Time `json:"last_accessed_at,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

// CreatePreviewSession creates a new preview session
func (r *PreviewSessionRepository) CreatePreviewSession(session *PreviewSession) error {
	query := `
		INSERT INTO preview_sessions (id, user_id, course_id, lecture_id, session_started_at,
			session_duration_seconds, preview_limit_seconds, preview_exhausted, ip_address,
			last_accessed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	if session.ID == "" {
		session.ID = uuid.New().String()
	}

	now := time.Now()
	if session.CreatedAt.IsZero() {
		session.CreatedAt = now
	}
	if session.UpdatedAt.IsZero() {
		session.UpdatedAt = now
	}

	_, err := r.db.Exec(query, session.ID, session.UserID, session.CourseID, session.LectureID,
		session.SessionStartedAt, session.SessionDurationSeconds, session.PreviewLimitSeconds,
		session.PreviewExhausted, session.IPAddress, session.LastAccessedAt,
		session.CreatedAt, session.UpdatedAt)

	return err
}

// GetPreviewSession retrieves a preview session by ID
func (r *PreviewSessionRepository) GetPreviewSession(sessionID string) (*PreviewSession, error) {
	query := `
		SELECT id, user_id, course_id, lecture_id, session_started_at,
			session_duration_seconds, preview_limit_seconds, preview_exhausted,
			ip_address, last_accessed_at, created_at, updated_at
		FROM preview_sessions
		WHERE id = $1
	`

	session := &PreviewSession{}
	err := r.db.QueryRow(query, sessionID).Scan(
		&session.ID, &session.UserID, &session.CourseID, &session.LectureID,
		&session.SessionStartedAt, &session.SessionDurationSeconds, &session.PreviewLimitSeconds,
		&session.PreviewExhausted, &session.IPAddress, &session.LastAccessedAt,
		&session.CreatedAt, &session.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetActivePreviewSession retrieves active preview session for user and course
func (r *PreviewSessionRepository) GetActivePreviewSession(userID, courseID string) (*PreviewSession, error) {
	query := `
		SELECT id, user_id, course_id, lecture_id, session_started_at,
			session_duration_seconds, preview_limit_seconds, preview_exhausted,
			ip_address, last_accessed_at, created_at, updated_at
		FROM preview_sessions
		WHERE user_id = $1 AND course_id = $2 AND NOT preview_exhausted
		ORDER BY created_at DESC
		LIMIT 1
	`

	session := &PreviewSession{}
	err := r.db.QueryRow(query, userID, courseID).Scan(
		&session.ID, &session.UserID, &session.CourseID, &session.LectureID,
		&session.SessionStartedAt, &session.SessionDurationSeconds, &session.PreviewLimitSeconds,
		&session.PreviewExhausted, &session.IPAddress, &session.LastAccessedAt,
		&session.CreatedAt, &session.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No active session found
		}
		return nil, err
	}

	return session, nil
}

// UpdatePreviewSession updates a preview session
func (r *PreviewSessionRepository) UpdatePreviewSession(session *PreviewSession) error {
	query := `
		UPDATE preview_sessions
		SET session_duration_seconds = $1, preview_exhausted = $2, last_accessed_at = $3, updated_at = $4
		WHERE id = $5
	`

	session.UpdatedAt = time.Now()
	now := time.Now()
	session.LastAccessedAt = &now
	_, err := r.db.Exec(query, session.SessionDurationSeconds, session.PreviewExhausted, session.LastAccessedAt, session.UpdatedAt, session.ID)
	return err
}

// GetByUserAndLecture retrieves a preview session by user and lecture
func (r *PreviewSessionRepository) GetByUserAndLecture(ctx context.Context, userID, lectureID string) (*PreviewSession, error) {
	query := `
		SELECT id, user_id, course_id, lecture_id, session_started_at,
			session_duration_seconds, preview_limit_seconds, preview_exhausted,
			ip_address, last_accessed_at, created_at, updated_at
		FROM preview_sessions
		WHERE user_id = $1 AND lecture_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	session := &PreviewSession{}
	err := r.db.QueryRowContext(ctx, query, userID, lectureID).Scan(
		&session.ID, &session.UserID, &session.CourseID, &session.LectureID,
		&session.SessionStartedAt, &session.SessionDurationSeconds, &session.PreviewLimitSeconds,
		&session.PreviewExhausted, &session.IPAddress, &session.LastAccessedAt,
		&session.CreatedAt, &session.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("preview session not found")
		}
		return nil, err
	}

	return session, nil
}

// ExpirePreviewSession marks a preview session as expired
func (r *PreviewSessionRepository) ExpirePreviewSession(sessionID string) error {
	query := `
		UPDATE preview_sessions
		SET preview_exhausted = true, updated_at = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(query, time.Now(), sessionID)
	return err
}

// CleanupExpiredSessions removes expired preview sessions
func (r *PreviewSessionRepository) CleanupExpiredSessions() error {
	query := `
		DELETE FROM preview_sessions
		WHERE preview_exhausted = true AND updated_at < NOW() - INTERVAL '24 hours'
	`

	_, err := r.db.Exec(query)
	return err
}