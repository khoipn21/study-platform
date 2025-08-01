package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type ChatbotHandler struct {
	chatbotServiceURL string
	client           *http.Client
	upgrader         websocket.Upgrader
}

func NewChatbotHandler(chatbotServiceURL string) *ChatbotHandler {
	return &ChatbotHandler{
		chatbotServiceURL: chatbotServiceURL,
		client:           &http.Client{},
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow connections from any origin in development
			},
		},
	}
}

func (h *ChatbotHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	h.proxyRequest(w, r, "/api/v1/sessions", "POST")
}

func (h *ChatbotHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]
	path := fmt.Sprintf("/api/v1/sessions/%s", sessionID)
	h.proxyRequest(w, r, path, "GET")
}

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

func (h *ChatbotHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	h.proxyRequest(w, r, "/api/v1/chat", "POST")
}

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
	// Get user info from context (set by auth middleware)
	userID := r.Header.Get("X-User-ID")
	userRole := r.Header.Get("X-User-Role")

	if userID == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Create WebSocket URL for chatbot service
	wsURL := h.chatbotServiceURL + "/api/v1/ws"
	u, err := url.Parse(wsURL)
	if err != nil {
		http.Error(w, "Invalid WebSocket URL", http.StatusInternalServerError)
		return
	}

	// Change scheme to ws/wss
	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else if u.Scheme == "https" {
		u.Scheme = "wss" 
	}

	// Upgrade connection to WebSocket
	clientConn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	// Connect to chatbot service WebSocket
	serviceHeaders := http.Header{}
	serviceHeaders.Set("X-User-ID", userID)
	if userRole != "" {
		serviceHeaders.Set("X-User-Role", userRole)
	}

	serviceConn, _, err := websocket.DefaultDialer.Dial(u.String(), serviceHeaders)
	if err != nil {
		clientConn.WriteMessage(websocket.TextMessage, []byte(`{"type":"error","data":{"error":"Failed to connect to chat service"}}`))
		return
	}
	defer serviceConn.Close()

	// Proxy messages between client and service
	go h.proxyWebSocketMessages(clientConn, serviceConn)
	h.proxyWebSocketMessages(serviceConn, clientConn)
}

func (h *ChatbotHandler) proxyWebSocketMessages(from, to *websocket.Conn) {
	defer to.Close()
	for {
		messageType, message, err := from.ReadMessage()
		if err != nil {
			break
		}
		
		err = to.WriteMessage(messageType, message)
		if err != nil {
			break
		}
	}
}

func (h *ChatbotHandler) proxyRequest(w http.ResponseWriter, r *http.Request, path, method string) {
	// Get user info from context (set by auth middleware)
	userID := r.Header.Get("X-User-ID")
	userRole := r.Header.Get("X-User-Role")

	// Create target URL
	targetURL := h.chatbotServiceURL + path

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