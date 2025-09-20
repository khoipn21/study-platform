package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type VideoHandler struct {
	videoServiceURL string
	upgrader        websocket.Upgrader
}

func NewVideoHandler() *VideoHandler {
	videoServiceURL := os.Getenv("VIDEO_SERVICE_URL")
	if videoServiceURL == "" {
		videoServiceURL = "http://localhost:8084"
	}

	return &VideoHandler{
		videoServiceURL: videoServiceURL,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for development - in production, implement proper CORS
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// GetUploadURL proxies get upload URL requests to video service
func (vh *VideoHandler) GetUploadURL(w http.ResponseWriter, r *http.Request) {
	vh.proxyRequest(w, r, "POST", "/api/videos/upload-url", true)
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

// WebSocketProxy handles WebSocket connections to video service
func (vh *VideoHandler) WebSocketProxy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["session_id"]

	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	// Build WebSocket URL for video service
	wsURL, err := vh.buildWebSocketURL(r, fmt.Sprintf("/api/videos/ws/%s", sessionID))
	if err != nil {
		http.Error(w, "Failed to build WebSocket URL", http.StatusInternalServerError)
		return
	}

	// Upgrade the client connection
	clientConn, err := vh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	// Create headers for backend connection
	headers := http.Header{}

	// Forward authentication headers
	if auth := r.Header.Get("Authorization"); auth != "" {
		headers.Set("Authorization", auth)
	}

	// Forward relevant headers
	if origin := r.Header.Get("Origin"); origin != "" {
		headers.Set("Origin", origin)
	}

	// Connect to video service WebSocket
	backendConn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to video service: %v", err), http.StatusInternalServerError)
		return
	}
	defer backendConn.Close()

	// Start proxying messages bidirectionally
	vh.proxyWebSocketMessages(clientConn, backendConn)
}

// WebSocketStats proxies WebSocket stats requests to video service
func (vh *VideoHandler) WebSocketStats(w http.ResponseWriter, r *http.Request) {
	vh.proxyRequest(w, r, "GET", "/api/videos/ws/stats", false)
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

	// Copy response headers (excluding CORS headers to prevent duplication)
	for key, values := range resp.Header {
		// Skip ALL CORS headers since API Gateway already handles them
		if strings.HasPrefix(key, "Access-Control-") {
			continue
		}
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set status and return response
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}


// buildWebSocketURL converts HTTP URL to WebSocket URL for video service
func (vh *VideoHandler) buildWebSocketURL(r *http.Request, path string) (string, error) {
	// Parse the video service URL
	baseURL, err := url.Parse(vh.videoServiceURL)
	if err != nil {
		return "", err
	}

	// Convert HTTP(S) to WS(S)
	scheme := "ws"
	if baseURL.Scheme == "https" {
		scheme = "wss"
	}

	// Build WebSocket URL
	wsURL := fmt.Sprintf("%s://%s%s", scheme, baseURL.Host, path)

	// Add query parameters from original request
	if r.URL.RawQuery != "" {
		wsURL += "?" + r.URL.RawQuery
	}

	return wsURL, nil
}

// proxyWebSocketMessages handles bidirectional message proxying
func (vh *VideoHandler) proxyWebSocketMessages(clientConn, backendConn *websocket.Conn) {
	// Create channels for error handling
	clientDone := make(chan error, 1)
	backendDone := make(chan error, 1)

	// Proxy messages from client to backend
	go func() {
		defer close(clientDone)
		for {
			messageType, message, err := clientConn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					clientDone <- fmt.Errorf("client read error: %v", err)
				}
				return
			}

			err = backendConn.WriteMessage(messageType, message)
			if err != nil {
				clientDone <- fmt.Errorf("backend write error: %v", err)
				return
			}
		}
	}()

	// Proxy messages from backend to client
	go func() {
		defer close(backendDone)
		for {
			messageType, message, err := backendConn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					backendDone <- fmt.Errorf("backend read error: %v", err)
				}
				return
			}

			err = clientConn.WriteMessage(messageType, message)
			if err != nil {
				backendDone <- fmt.Errorf("client write error: %v", err)
				return
			}
		}
	}()

	// Wait for either connection to close or error
	select {
	case err := <-clientDone:
		if err != nil {
			// Log error if needed
			fmt.Printf("Client connection closed: %v\n", err)
		}
	case err := <-backendDone:
		if err != nil {
			// Log error if needed
			fmt.Printf("Backend connection closed: %v\n", err)
		}
	}

	// Close both connections gracefully
	clientConn.WriteMessage(websocket.CloseMessage, []byte{})
	backendConn.WriteMessage(websocket.CloseMessage, []byte{})
}

// UpdateVideoStatus proxies video status update requests to video service
func (vh *VideoHandler) UpdateVideoStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID := vars["video_id"]
	path := fmt.Sprintf("/api/videos/%s/status", videoID)
	vh.proxyRequest(w, r, "PUT", path, true)
}

// GetVideoHealth checks video service health
func (vh *VideoHandler) GetVideoHealth(w http.ResponseWriter, r *http.Request) {
	vh.proxyRequest(w, r, "GET", "/health", false)
}