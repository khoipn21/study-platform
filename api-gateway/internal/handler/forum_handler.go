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
func (h *ForumHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	h.proxyRequest(w, r, "/api/v1/topics", "POST")
}

func (h *ForumHandler) GetTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["topicId"]
	path := fmt.Sprintf("/api/v1/topics/%s", topicID)
	h.proxyRequest(w, r, path, "GET")
}

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
	if uid, ok := r.Context().Value("user_id").(string); ok {
		userID = uid
	}
	if role, ok := r.Context().Value("user_role").(string); ok {
		userRole = role
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
func (h *ForumHandler) ApproveTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["topicId"]
	path := fmt.Sprintf("/api/v1/topics/%s/approve", topicID)
	h.proxyRequest(w, r, path, "PUT")
}

func (h *ForumHandler) ApprovePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]
	path := fmt.Sprintf("/api/v1/posts/%s/approve", postID)
	h.proxyRequest(w, r, path, "PUT")
}

func (h *ForumHandler) GetPendingTopics(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query().Encode()
	path := "/api/v1/pending/topics"
	if queryParams != "" {
		path += "?" + queryParams
	}
	h.proxyRequest(w, r, path, "GET")
}

// Pin management handlers
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