package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/study-platform/course-service/internal/model"
	"github.com/study-platform/course-service/internal/repository"
	"github.com/study-platform/pkg/logger"
)

type NotesService struct {
	notesRepo *repository.NotesRepository
	logger    logger.Logger
}

func NewNotesService(notesRepo *repository.NotesRepository, logger logger.Logger) *NotesService {
	return &NotesService{
		notesRepo: notesRepo,
		logger:    logger,
	}
}

// CreateNote creates a new note for a user on a specific lecture
func (s *NotesService) CreateNote(userID, courseID, lectureID uuid.UUID, req *model.CreateNoteRequest) (*model.NoteResponse, error) {
	s.logger.Infof("Creating note for user %s on lecture %s", userID, lectureID)

	// Validate input
	if req.Title == "" {
		return nil, fmt.Errorf("note title is required")
	}
	if req.Content == "" {
		return nil, fmt.Errorf("note content is required")
	}

	note, err := s.notesRepo.CreateNote(userID, courseID, lectureID, req)
	if err != nil {
		s.logger.Errorf("Failed to create note: %v", err)
		return nil, err
	}

	return &model.NoteResponse{
		ID:               note.ID,
		Title:            note.Title,
		Content:          note.Content,
		TimestampSeconds: note.TimestampSeconds,
		CreatedAt:        note.CreatedAt,
		UpdatedAt:        note.UpdatedAt,
	}, nil
}

// GetNotesByLecture retrieves all notes for a user on a specific lecture
func (s *NotesService) GetNotesByLecture(userID, courseID, lectureID uuid.UUID) ([]*model.NoteResponse, error) {
	s.logger.Infof("Getting notes for user %s on lecture %s", userID, lectureID)

	notes, err := s.notesRepo.GetNotesByLecture(userID, courseID, lectureID)
	if err != nil {
		s.logger.Errorf("Failed to get notes by lecture: %v", err)
		return nil, err
	}

	var responses []*model.NoteResponse
	for _, note := range notes {
		responses = append(responses, &model.NoteResponse{
			ID:               note.ID,
			Title:            note.Title,
			Content:          note.Content,
			TimestampSeconds: note.TimestampSeconds,
			CreatedAt:        note.CreatedAt,
			UpdatedAt:        note.UpdatedAt,
		})
	}

	return responses, nil
}

// GetNotesByCourse retrieves all notes for a user on a specific course
func (s *NotesService) GetNotesByCourse(userID, courseID uuid.UUID) ([]*model.NoteResponse, error) {
	s.logger.Infof("Getting notes for user %s on course %s", userID, courseID)

	notes, err := s.notesRepo.GetNotesByCourse(userID, courseID)
	if err != nil {
		s.logger.Errorf("Failed to get notes by course: %v", err)
		return nil, err
	}

	var responses []*model.NoteResponse
	for _, note := range notes {
		responses = append(responses, &model.NoteResponse{
			ID:               note.ID,
			Title:            note.Title,
			Content:          note.Content,
			TimestampSeconds: note.TimestampSeconds,
			CreatedAt:        note.CreatedAt,
			UpdatedAt:        note.UpdatedAt,
		})
	}

	return responses, nil
}

// GetNoteByID retrieves a specific note by its ID
func (s *NotesService) GetNoteByID(userID, noteID uuid.UUID) (*model.NoteResponse, error) {
	s.logger.Infof("Getting note %s for user %s", noteID, userID)

	note, err := s.notesRepo.GetNoteByID(userID, noteID)
	if err != nil {
		s.logger.Errorf("Failed to get note: %v", err)
		return nil, err
	}

	return &model.NoteResponse{
		ID:               note.ID,
		Title:            note.Title,
		Content:          note.Content,
		TimestampSeconds: note.TimestampSeconds,
		CreatedAt:        note.CreatedAt,
		UpdatedAt:        note.UpdatedAt,
	}, nil
}

// UpdateNote updates an existing note
func (s *NotesService) UpdateNote(userID, noteID uuid.UUID, req *model.UpdateNoteRequest) (*model.NoteResponse, error) {
	s.logger.Infof("Updating note %s for user %s", noteID, userID)

	note, err := s.notesRepo.UpdateNote(userID, noteID, req)
	if err != nil {
		s.logger.Errorf("Failed to update note: %v", err)
		return nil, err
	}

	return &model.NoteResponse{
		ID:               note.ID,
		Title:            note.Title,
		Content:          note.Content,
		TimestampSeconds: note.TimestampSeconds,
		CreatedAt:        note.CreatedAt,
		UpdatedAt:        note.UpdatedAt,
	}, nil
}

// DeleteNote deletes a note
func (s *NotesService) DeleteNote(userID, noteID uuid.UUID) error {
	s.logger.Infof("Deleting note %s for user %s", noteID, userID)

	err := s.notesRepo.DeleteNote(userID, noteID)
	if err != nil {
		s.logger.Errorf("Failed to delete note: %v", err)
		return err
	}

	return nil
}