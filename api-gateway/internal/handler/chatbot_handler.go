package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type ChatbotHandler struct {
	chatbotServiceURL string
	client           *http.Client
}

func NewChatbotHandler(chatbotServiceURL string) *ChatbotHandler {
	return &ChatbotHandler{
		chatbotServiceURL: chatbotServiceURL,
		client:           &http.Client{},
	}
}

// CreateSession godoc
// @Summary      Create chat session
// @Description  Create new chatbot session
// @Tags         Chatbot
// @Accept       json
// @Produce      json
// @Success      201 {object} APIResponse "Session created"
// @Security     BearerAuth
// @Router       /chat/sessions [post]
func (h *ChatbotHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	h.proxyRequest(w, r, "/api/v1/sessions", "POST")
}

func (h *ChatbotHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]
	path := fmt.Sprintf("/api/v1/sessions/%s", sessionID)
	h.proxyRequest(w, r, path, "GET")
}

// GetUserSessions godoc
// @Summary      Get user sessions
// @Description  Get all chat sessions for user
// @Tags         Chatbot
// @Produce      json
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Success      200 {object} APIResponse "List of sessions"
// @Security     BearerAuth
// @Router       /chat/sessions [get]
func (h *ChatbotHandler) GetUserSessions(w http.ResponseWriter, r *http.Request) {
	// Forward query parameters
	queryParams := r.URL.Query().Encode()
	path := "/api/v1/sessions"
	if queryParams != "" {
		path += "?" + queryParams
	}
	h.proxyRequest(w, r, path, "GET")
}

func (h *ChatbotHandler) UpdateSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]
	path := fmt.Sprintf("/api/v1/sessions/%s", sessionID)
	h.proxyRequest(w, r, path, "PUT")
}

func (h *ChatbotHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]
	path := fmt.Sprintf("/api/v1/sessions/%s", sessionID)
	h.proxyRequest(w, r, path, "DELETE")
}

// SendMessage godoc
// @Summary      Send message
// @Description  Send message to chatbot
// @Tags         Chatbot
// @Accept       json
// @Produce      json
// @Success      200 {object} APIResponse "Message sent"
// @Security     BearerAuth
// @Router       /chat/message [post]
func (h *ChatbotHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	h.proxyRequest(w, r, "/api/v1/chat", "POST")
}

// GetMessages godoc
// @Summary      Get session messages
// @Description  Get messages for a session
// @Tags         Chatbot
// @Produce      json
// @Param        sessionId path string true "Session ID"
// @Success      200 {object} APIResponse "List of messages"
// @Security     BearerAuth
// @Router       /chat/sessions/{sessionId}/messages [get]
func (h *ChatbotHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]
	queryParams := r.URL.Query().Encode()
	path := fmt.Sprintf("/api/v1/sessions/%s/messages", sessionID)
	if queryParams != "" {
		path += "?" + queryParams
	}
	h.proxyRequest(w, r, path, "GET")
}

func (h *ChatbotHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// WebSocket functionality disabled for simplicity
	http.Error(w, "WebSocket not available through gateway", http.StatusNotImplemented)
}

// ListUserSessions godoc
// @Summary      List user chat sessions
// @Description  Get all chat sessions with history for the current user
// @Tags         Chatbot
// @Produce      json
// @Success      200 {object} APIResponse "List of chat sessions"
// @Security     BearerAuth
// @Router       /chat/history [get]
func (h *ChatbotHandler) ListUserSessions(w http.ResponseWriter, r *http.Request) {
	h.proxyRequest(w, r, "/api/v1/chat/history", "GET")
}

// GetSessionHistory godoc
// @Summary      Get session history
// @Description  Get all messages for a specific chat session
// @Tags         Chatbot
// @Produce      json
// @Param        sessionId path string true "Session ID"
// @Success      200 {object} APIResponse "Session messages"
// @Security     BearerAuth
// @Router       /chat/history/{sessionId} [get]
func (h *ChatbotHandler) GetSessionHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]
	path := fmt.Sprintf("/api/v1/chat/history/%s", sessionID)
	h.proxyRequest(w, r, path, "GET")
}

// GetRateLimit godoc
// @Summary      Get rate limit
// @Description  Get current rate limit status for the user
// @Tags         Chatbot
// @Produce      json
// @Success      200 {object} APIResponse "Rate limit info"
// @Security     BearerAuth
// @Router       /rate-limit [get]
func (h *ChatbotHandler) GetRateLimit(w http.ResponseWriter, r *http.Request) {
	h.proxyRequest(w, r, "/api/v1/rate-limit", "GET")
}

// DeleteHistorySession godoc
// @Summary      Delete history session
// @Description  Delete a chat history session
// @Tags         Chatbot
// @Produce      json
// @Param        sessionId path string true "Session ID"
// @Success      200 {object} APIResponse "Session deleted"
// @Security     BearerAuth
// @Router       /chat/history/{sessionId} [delete]
func (h *ChatbotHandler) DeleteHistorySession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]
	path := fmt.Sprintf("/api/v1/chat/history/%s", sessionID)
	h.proxyRequest(w, r, path, "DELETE")
}

func (h *ChatbotHandler) proxyRequest(w http.ResponseWriter, r *http.Request, path, method string) {
	// CRITICAL FIX for BUG-003: Get user info from context (set by auth middleware)
	var userID, userRole string

	// Extract user ID from context
	if userIDValue := r.Context().Value("user_id"); userIDValue != nil {
		userID = userIDValue.(string)
	}

	// Extract user role from context
	if userRoleValue := r.Context().Value("user_role"); userRoleValue != nil {
		userRole = userRoleValue.(string)
	}

	// Create target URL
	targetURL := h.chatbotServiceURL + path

	// Read request body
	var body []byte
	if r.Body != nil {
		var err error
		body, err = io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
	}

	// Create new request
	req, err := http.NewRequest(method, targetURL, bytes.NewReader(body))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Add user information headers
	if userID != "" {
		req.Header.Set("X-User-ID", userID)
	}
	if userRole != "" {
		req.Header.Set("X-User-Role", userRole)
	}

	// Make request
	resp, err := h.client.Do(req)
	if err != nil {
		http.Error(w, "Failed to connect to chat service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response headers (skip CORS headers to prevent duplication - API Gateway handles CORS)
	for key, values := range resp.Header {
		// Skip CORS headers since API Gateway middleware already sets them
		if strings.HasPrefix(key, "Access-Control-") {
			continue
		}
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	io.Copy(w, resp.Body)
}