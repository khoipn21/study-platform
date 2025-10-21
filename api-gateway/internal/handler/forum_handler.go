package handler

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ForumHandler struct {
	forumServiceURL string
	client          *http.Client
}

func NewForumHandler(forumServiceURL string) *ForumHandler {
	return &ForumHandler{
		forumServiceURL: forumServiceURL,
		client:          &http.Client{},
	}
}

// Topic handlers
// CreateTopic godoc
// @Summary      Create topic
// @Description  Create new forum topic
// @Tags         Forum
// @Accept       json
// @Produce      json
// @Success      201 {object} APIResponse "Topic created"
// @Failure      400 {object} APIResponse "Invalid request"
// @Security     BearerAuth
// @Router       /forum/topics [post]
func (h *ForumHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	h.proxyRequest(w, r, "/api/v1/topics", "POST")
}

// GetTopic godoc
// @Summary      Get topic
// @Description  Get forum topic by ID
// @Tags         Forum
// @Produce      json
// @Param        topicId path string true "Topic ID"
// @Success      200 {object} APIResponse "Topic details"
// @Failure      404 {object} APIResponse "Topic not found"
// @Router       /forum/topics/{topicId} [get]
func (h *ForumHandler) GetTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["topicId"]
	path := fmt.Sprintf("/api/v1/topics/%s", topicID)
	h.proxyRequest(w, r, path, "GET")
}

// ListTopics godoc
// @Summary      List forum topics
// @Description  Get list of forum topics with pagination
// @Tags         Forum
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Items per page"
// @Param        course_id query string false "Filter by course"
// @Success      200 {object} APIResponse "List of topics"
// @Router       /forum/topics [get]
func (h *ForumHandler) ListTopics(w http.ResponseWriter, r *http.Request) {
	// Forward query parameters
	queryParams := r.URL.Query().Encode()
	path := "/api/v1/topics"
	if queryParams != "" {
		path += "?" + queryParams
	}
	h.proxyRequest(w, r, path, "GET")
}

func (h *ForumHandler) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["topicId"]
	path := fmt.Sprintf("/api/v1/topics/%s", topicID)
	h.proxyRequest(w, r, path, "PUT")
}

func (h *ForumHandler) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["topicId"]
	path := fmt.Sprintf("/api/v1/topics/%s", topicID)
	h.proxyRequest(w, r, path, "DELETE")
}

func (h *ForumHandler) ToggleTopicSticky(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["topicId"]
	path := fmt.Sprintf("/api/v1/topics/%s/sticky", topicID)
	h.proxyRequest(w, r, path, "PUT")
}

func (h *ForumHandler) ToggleTopicLock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["topicId"]
	path := fmt.Sprintf("/api/v1/topics/%s/lock", topicID)
	h.proxyRequest(w, r, path, "PUT")
}

// Post handlers
func (h *ForumHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	h.proxyRequest(w, r, "/api/v1/posts", "POST")
}

func (h *ForumHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]
	path := fmt.Sprintf("/api/v1/posts/%s", postID)
	h.proxyRequest(w, r, path, "GET")
}

func (h *ForumHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["topicId"]
	queryParams := r.URL.Query().Encode()
	path := fmt.Sprintf("/api/v1/topics/%s/posts", topicID)
	if queryParams != "" {
		path += "?" + queryParams
	}
	h.proxyRequest(w, r, path, "GET")
}

func (h *ForumHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]
	path := fmt.Sprintf("/api/v1/posts/%s", postID)
	h.proxyRequest(w, r, path, "PUT")
}

func (h *ForumHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]
	path := fmt.Sprintf("/api/v1/posts/%s", postID)
	h.proxyRequest(w, r, path, "DELETE")
}

func (h *ForumHandler) MarkPostAsAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]
	path := fmt.Sprintf("/api/v1/posts/%s/answer", postID)
	h.proxyRequest(w, r, path, "PUT")
}

func (h *ForumHandler) TogglePostPin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]
	path := fmt.Sprintf("/api/v1/posts/%s/pin", postID)
	h.proxyRequest(w, r, path, "PUT")
}

// Voting handlers
func (h *ForumHandler) VotePost(w http.ResponseWriter, r *http.Request) {
	h.proxyRequest(w, r, "/api/v1/votes", "POST")
}

func (h *ForumHandler) RemoveVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]
	path := fmt.Sprintf("/api/v1/posts/%s/vote", postID)
	h.proxyRequest(w, r, path, "DELETE")
}

// Search handler
func (h *ForumHandler) SearchTopics(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query().Encode()
	path := "/api/v1/search"
	if queryParams != "" {
		path += "?" + queryParams
	}
	h.proxyRequest(w, r, path, "GET")
}

// Course-specific forum handlers
func (h *ForumHandler) ListCourseTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["courseId"]
	queryParams := r.URL.Query()
	queryParams.Set("course_id", courseID)
	path := "/api/v1/topics?" + queryParams.Encode()
	h.proxyRequest(w, r, path, "GET")
}

func (h *ForumHandler) proxyRequest(w http.ResponseWriter, r *http.Request, path, method string) {
	// Get user info from context (set by auth middleware)
	var userID, userRole string

	// Extract user ID from context
	if userIDValue := r.Context().Value("user_id"); userIDValue != nil {
		userID = fmt.Sprintf("%v", userIDValue)
	}

	// Extract user role from context
	if userRoleValue := r.Context().Value("user_role"); userRoleValue != nil {
		userRole = fmt.Sprintf("%v", userRoleValue)
	}


	// Create target URL
	targetURL := h.forumServiceURL + path

	// Read request body
	var body []byte
	if r.Body != nil {
		var err error
		body, err = ReadAll(r.Body)
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

	// Add user information headers (forum service expects these)
	if userID != "" {
		req.Header.Set("X-User-ID", userID)
	}
	if userRole != "" {
		req.Header.Set("X-User-Role", userRole)
	}

	// Make request
	resp, err := h.client.Do(req)
	if err != nil {
		http.Error(w, "Failed to connect to forum service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	responseBody, err := ReadAll(resp.Body)
	if err != nil {
		return
	}
	w.Write(responseBody)
}
// Approval handlers
func (h *ForumHandler) GetPendingTopics(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query().Encode()
	path := "/api/v1/pending/topics"
	if queryParams != "" {
		path += "?" + queryParams
	}
	h.proxyRequest(w, r, path, "GET")
}

func (h *ForumHandler) ApproveTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["topicId"]
	path := fmt.Sprintf("/api/v1/topics/%s/approve", topicID)
	h.proxyRequest(w, r, path, "PUT")
}

func (h *ForumHandler) RejectTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["topicId"]
	path := fmt.Sprintf("/api/v1/topics/%s/reject", topicID)
	h.proxyRequest(w, r, path, "PUT")
}

func (h *ForumHandler) ApprovePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]
	path := fmt.Sprintf("/api/v1/posts/%s/approve", postID)
	h.proxyRequest(w, r, path, "PUT")
}

func (h *ForumHandler) RejectPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]
	path := fmt.Sprintf("/api/v1/posts/%s/reject", postID)
	h.proxyRequest(w, r, path, "PUT")
}

func (h *ForumHandler) SetTopicPinOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["topicId"]
	path := fmt.Sprintf("/api/v1/topics/%s/pin-order", topicID)
	h.proxyRequest(w, r, path, "PUT")
}

func (h *ForumHandler) SetPostPinOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]
	path := fmt.Sprintf("/api/v1/posts/%s/pin-order", postID)
	h.proxyRequest(w, r, path, "PUT")
}
