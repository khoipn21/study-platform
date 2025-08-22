package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type VideoHandler struct {
	videoServiceURL string
}

func NewVideoHandler() *VideoHandler {
	videoServiceURL := os.Getenv("VIDEO_SERVICE_URL")
	if videoServiceURL == "" {
		videoServiceURL = "http://localhost:8084"
	}

	return &VideoHandler{
		videoServiceURL: videoServiceURL,
	}
}

// UploadVideo proxies video upload requests to video service
func (vh *VideoHandler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	vh.proxyRequest(w, r, "POST", "/api/videos/upload", true)
}

// GetVideo proxies get video requests to video service
func (vh *VideoHandler) GetVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID := vars["video_id"]
	path := fmt.Sprintf("/api/videos/%s", videoID)
	vh.proxyRequest(w, r, "GET", path, false)
}

// UpdateVideo proxies video update requests to video service
func (vh *VideoHandler) UpdateVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID := vars["video_id"]
	path := fmt.Sprintf("/api/videos/%s", videoID)
	vh.proxyRequest(w, r, "PUT", path, true)
}

// DeleteVideo proxies video deletion requests to video service
func (vh *VideoHandler) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID := vars["video_id"]
	path := fmt.Sprintf("/api/videos/%s", videoID)
	vh.proxyRequest(w, r, "DELETE", path, true)
}

// CreateViewingSession proxies session creation requests to video service
func (vh *VideoHandler) CreateViewingSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID := vars["video_id"]
	path := fmt.Sprintf("/api/videos/%s/sessions", videoID)
	vh.proxyRequest(w, r, "POST", path, true)
}

// UpdateSessionProgress proxies session progress updates to video service
func (vh *VideoHandler) UpdateSessionProgress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["session_id"]
	path := fmt.Sprintf("/api/videos/sessions/%s/progress", sessionID)
	vh.proxyRequest(w, r, "PUT", path, true)
}

// UpdateNetworkStatus proxies network status updates to video service
func (vh *VideoHandler) UpdateNetworkStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["session_id"]
	path := fmt.Sprintf("/api/videos/sessions/%s/network", sessionID)
	vh.proxyRequest(w, r, "POST", path, true)
}

// GetVideoAnalytics proxies video analytics requests to video service
func (vh *VideoHandler) GetVideoAnalytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID := vars["video_id"]
	path := fmt.Sprintf("/api/videos/%s/analytics", videoID)
	vh.proxyRequest(w, r, "GET", path, true)
}

// SearchVideos proxies video search requests to video service
func (vh *VideoHandler) SearchVideos(w http.ResponseWriter, r *http.Request) {
	// Pass query parameters
	queryParams := r.URL.RawQuery
	path := "/api/videos/search"
	if queryParams != "" {
		path += "?" + queryParams
	}
	vh.proxyRequest(w, r, "GET", path, false)
}

// ListUserVideos proxies user video listing requests to video service
func (vh *VideoHandler) ListUserVideos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	queryParams := r.URL.RawQuery
	path := fmt.Sprintf("/api/videos/user/%s", userID)
	if queryParams != "" {
		path += "?" + queryParams
	}
	vh.proxyRequest(w, r, "GET", path, true)
}

// ListCourseVideos proxies course video listing requests to video service
func (vh *VideoHandler) ListCourseVideos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["course_id"]
	queryParams := r.URL.RawQuery
	path := fmt.Sprintf("/api/videos/course/%s", courseID)
	if queryParams != "" {
		path += "?" + queryParams
	}
	vh.proxyRequest(w, r, "GET", path, false)
}

// CloudflareWebhook proxies Cloudflare webhook requests to video service
func (vh *VideoHandler) CloudflareWebhook(w http.ResponseWriter, r *http.Request) {
	vh.proxyRequest(w, r, "POST", "/api/videos/webhooks/cloudflare", false)
}

// WebSocketProxy handles WebSocket connections to video service (placeholder)
func (vh *VideoHandler) WebSocketProxy(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement WebSocket proxying when gorilla/websocket is added
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error": "WebSocket functionality not yet implemented in API Gateway"}`))
}

// WebSocketStats proxies WebSocket stats requests to video service
func (vh *VideoHandler) WebSocketStats(w http.ResponseWriter, r *http.Request) {
	vh.proxyRequest(w, r, "GET", "/api/videos/ws/stats", true)
}

// BroadcastMessage proxies broadcast requests to video service
func (vh *VideoHandler) BroadcastMessage(w http.ResponseWriter, r *http.Request) {
	vh.proxyRequest(w, r, "POST", "/api/videos/ws/broadcast", true)
}

// SendToSession proxies session message requests to video service
func (vh *VideoHandler) SendToSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["session_id"]
	path := fmt.Sprintf("/api/videos/ws/session/%s/send", sessionID)
	vh.proxyRequest(w, r, "POST", path, true)
}

// GetSessionInfo proxies session info requests to video service
func (vh *VideoHandler) GetSessionInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["session_id"]
	path := fmt.Sprintf("/api/videos/ws/session/%s", sessionID)
	vh.proxyRequest(w, r, "GET", path, true)
}

// proxyRequest is a helper function to proxy HTTP requests to the video service
func (vh *VideoHandler) proxyRequest(w http.ResponseWriter, r *http.Request, method, path string, requireAuth bool) {
	// Build target URL
	targetURL := vh.videoServiceURL + path

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
	req, err := http.NewRequest(method, targetURL, bytes.NewBuffer(body))
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

	// Forward authentication if present
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	// Set content type if body is present
	if len(body) > 0 {
		req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
	}

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to connect to video service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set status and return response
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}


// GetVideoHealth checks video service health
func (vh *VideoHandler) GetVideoHealth(w http.ResponseWriter, r *http.Request) {
	vh.proxyRequest(w, r, "GET", "/health", false)
}