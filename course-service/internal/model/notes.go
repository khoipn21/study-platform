package model

import (
	"time"

	"github.com/google/uuid"
)

// Note represents a user note for a specific lecture
type Note struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	CourseID         uuid.UUID  `json:"course_id" db:"course_id"`
	LectureID        uuid.UUID  `json:"lecture_id" db:"lecture_id"`
	Title            string     `json:"title" db:"title"`
	Content          string     `json:"content" db:"content"`
	TimestampSeconds int        `json:"timestamp_seconds" db:"timestamp_seconds"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateNoteRequest represents the request to create a note
type CreateNoteRequest struct {
	Title            string `json:"title" validate:"required,max=500"`
	Content          string `json:"content" validate:"required"`
	TimestampSeconds int    `json:"timestamp_seconds,omitempty"`
}

// UpdateNoteRequest represents the request to update a note
type UpdateNoteRequest struct {
	Title            string `json:"title,omitempty" validate:"omitempty,max=500"`
	Content          string `json:"content,omitempty"`
	TimestampSeconds int    `json:"timestamp_seconds,omitempty"`
}

// NoteResponse represents the response for note operations
type NoteResponse struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	TimestampSeconds int       `json:"timestamp_seconds"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}