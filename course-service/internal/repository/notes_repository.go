package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/study-platform/course-service/internal/model"
)

type NotesRepository struct {
	db *sqlx.DB
}

func NewNotesRepository(db *sqlx.DB) *NotesRepository {
	return &NotesRepository{db: db}
}

// CreateNote creates a new note for a user on a specific lecture
func (r *NotesRepository) CreateNote(userID, courseID, lectureID uuid.UUID, req *model.CreateNoteRequest) (*model.Note, error) {
	note := &model.Note{
		ID:               uuid.New(),
		UserID:           userID,
		CourseID:         courseID,
		LectureID:        lectureID,
		Title:            req.Title,
		Content:          req.Content,
		TimestampSeconds: req.TimestampSeconds,
	}

	query := `
		INSERT INTO notes (id, user_id, course_id, lecture_id, title, content, timestamp_seconds)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at`

	err := r.db.QueryRow(
		query,
		note.ID, note.UserID, note.CourseID, note.LectureID,
		note.Title, note.Content, note.TimestampSeconds,
	).Scan(&note.CreatedAt, &note.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	return note, nil
}

// GetNotesByLecture retrieves all notes for a user on a specific lecture
func (r *NotesRepository) GetNotesByLecture(userID, courseID, lectureID uuid.UUID) ([]*model.Note, error) {
	var notes []*model.Note

	query := `
		SELECT id, user_id, course_id, lecture_id, title, content, timestamp_seconds, created_at, updated_at
		FROM notes
		WHERE user_id = $1 AND course_id = $2 AND lecture_id = $3
		ORDER BY timestamp_seconds ASC, created_at ASC`

	err := r.db.Select(&notes, query, userID, courseID, lectureID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes by lecture: %w", err)
	}

	return notes, nil
}

// GetNotesByCourse retrieves all notes for a user on a specific course
func (r *NotesRepository) GetNotesByCourse(userID, courseID uuid.UUID) ([]*model.Note, error) {
	var notes []*model.Note

	query := `
		SELECT id, user_id, course_id, lecture_id, title, content, timestamp_seconds, created_at, updated_at
		FROM notes
		WHERE user_id = $1 AND course_id = $2
		ORDER BY lecture_id, timestamp_seconds ASC, created_at ASC`

	err := r.db.Select(&notes, query, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes by course: %w", err)
	}

	return notes, nil
}

// GetNoteByID retrieves a specific note by its ID
func (r *NotesRepository) GetNoteByID(userID, noteID uuid.UUID) (*model.Note, error) {
	var note model.Note

	query := `
		SELECT id, user_id, course_id, lecture_id, title, content, timestamp_seconds, created_at, updated_at
		FROM notes
		WHERE id = $1 AND user_id = $2`

	err := r.db.Get(&note, query, noteID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("note not found")
		}
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	return &note, nil
}

// UpdateNote updates an existing note
func (r *NotesRepository) UpdateNote(userID, noteID uuid.UUID, req *model.UpdateNoteRequest) (*model.Note, error) {
	// First get the existing note
	existingNote, err := r.GetNoteByID(userID, noteID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Title != "" {
		existingNote.Title = req.Title
	}
	if req.Content != "" {
		existingNote.Content = req.Content
	}
	if req.TimestampSeconds != 0 {
		existingNote.TimestampSeconds = req.TimestampSeconds
	}

	query := `
		UPDATE notes
		SET title = $3, content = $4, timestamp_seconds = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2
		RETURNING updated_at`

	err = r.db.QueryRow(
		query,
		noteID, userID, existingNote.Title, existingNote.Content, existingNote.TimestampSeconds,
	).Scan(&existingNote.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	return existingNote, nil
}

// DeleteNote deletes a note
func (r *NotesRepository) DeleteNote(userID, noteID uuid.UUID) error {
	query := `DELETE FROM notes WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, noteID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}