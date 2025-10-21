package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/study-platform/pkg/logger"
)

type NotesHandler struct {
	courseServiceURL string
	httpClient       *http.Client
	logger           logger.Logger
}

func NewNotesHandler(courseServiceURL string, logger logger.Logger) *NotesHandler {
	return &NotesHandler{
		courseServiceURL: courseServiceURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// CreateNote creates a new note for a lecture
// CreateNote godoc
// @Summary      Create note
// @Description  Create a note for a lecture
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Param        course_id path string true "Course ID"
// @Param        lecture_id path string true "Lecture ID"
// @Success      201 {object} APIResponse "Note created"
// @Security     BearerAuth
// @Router       /notes/courses/{course_id}/lectures/{lecture_id} [post]
func (h *NotesHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	courseID := vars["course_id"]
	lectureID := vars["lecture_id"]

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Forward request to course service
	url := fmt.Sprintf("%s/courses/%s/lectures/%s/notes", h.courseServiceURL, courseID, lectureID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		h.logger.Errorf("Failed to create request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		h.logger.Errorf("Failed to call course service: %v", err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// GetNotesByLecture retrieves all notes for a specific lecture
// GetNotesByLecture godoc
// @Summary      Get lecture notes
// @Description  Get all notes for a lecture
// @Tags         Notes
// @Produce      json
// @Param        course_id path string true "Course ID"
// @Param        lecture_id path string true "Lecture ID"
// @Success      200 {object} APIResponse "List of notes"
// @Security     BearerAuth
// @Router       /notes/courses/{course_id}/lectures/{lecture_id} [get]
func (h *NotesHandler) GetNotesByLecture(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	courseID := vars["course_id"]
	lectureID := vars["lecture_id"]

	// Forward request to course service
	url := fmt.Sprintf("%s/courses/%s/lectures/%s/notes", h.courseServiceURL, courseID, lectureID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		h.logger.Errorf("Failed to create request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("X-User-ID", userID)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		h.logger.Errorf("Failed to call course service: %v", err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// GetNotesByCourse retrieves all notes for a specific course
func (h *NotesHandler) GetNotesByCourse(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	courseID := vars["course_id"]

	// Forward request to course service
	url := fmt.Sprintf("%s/courses/%s/notes", h.courseServiceURL, courseID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		h.logger.Errorf("Failed to create request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("X-User-ID", userID)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		h.logger.Errorf("Failed to call course service: %v", err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// UpdateNote updates an existing note
// UpdateNote godoc
// @Summary      Update note
// @Description  Update a note
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Param        note_id path string true "Note ID"
// @Success      200 {object} APIResponse "Note updated"
// @Security     BearerAuth
// @Router       /notes/{note_id} [put]
func (h *NotesHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	noteID := vars["note_id"]

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Forward request to course service
	url := fmt.Sprintf("%s/notes/%s", h.courseServiceURL, noteID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		h.logger.Errorf("Failed to create request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		h.logger.Errorf("Failed to call course service: %v", err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// DeleteNote deletes a note
// DeleteNote godoc
// @Summary      Delete note
// @Description  Delete a note
// @Tags         Notes
// @Param        note_id path string true "Note ID"
// @Success      200 {object} APIResponse "Note deleted"
// @Security     BearerAuth
// @Router       /notes/{note_id} [delete]
func (h *NotesHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	noteID := vars["note_id"]

	// Forward request to course service
	url := fmt.Sprintf("%s/notes/%s", h.courseServiceURL, noteID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		h.logger.Errorf("Failed to create request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("X-User-ID", userID)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		h.logger.Errorf("Failed to call course service: %v", err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// GetNote retrieves a specific note by ID
func (h *NotesHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	noteID := vars["note_id"]

	// Forward request to course service
	url := fmt.Sprintf("%s/notes/%s", h.courseServiceURL, noteID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		h.logger.Errorf("Failed to create request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("X-User-ID", userID)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		h.logger.Errorf("Failed to call course service: %v", err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}