package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"video-service/internal/model"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin for development
		// In production, implement proper origin checking
		return true
	},
}

// Client is a middleman between the websocket connection and the hub
type Client struct {
	hub *Hub

	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Session information
	SessionID string
	UserID    string
	VideoID   string

	// Connection metadata
	UserAgent string
	IPAddress string

	// Activity tracking
	LastActivity time.Time
}

// NewClient creates a new WebSocket client
func NewClient(hub *Hub, conn *websocket.Conn, sessionID, userID, videoID, userAgent, ipAddress string) *Client {
	return &Client{
		hub:          hub,
		conn:         conn,
		send:         make(chan []byte, 256),
		SessionID:    sessionID,
		UserID:       userID,
		VideoID:      videoID,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		LastActivity: time.Now(),
	}
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.UpdateActivity()
		return nil
	})

	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		c.UpdateActivity()

		// Process the message through the hub
		c.hub.ProcessClientMessage(c, messageType, message)
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// UpdateActivity updates the client's last activity time
func (c *Client) UpdateActivity() {
	c.LastActivity = time.Now()
}

// SendMessage sends a structured message to the client
func (c *Client) SendMessage(message model.WSMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	select {
	case c.send <- data:
	default:
		// Channel is full, close it
		close(c.send)
	}
}

// SendJSON sends a JSON message to the client
func (c *Client) SendJSON(data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		return
	}

	select {
	case c.send <- jsonData:
	default:
		// Channel is full, close it
		close(c.send)
	}
}

// SendText sends a text message to the client
func (c *Client) SendText(text string) {
	select {
	case c.send <- []byte(text):
	default:
		// Channel is full, close it
		close(c.send)
	}
}

// IsActive checks if the client connection is still active
func (c *Client) IsActive(threshold time.Duration) bool {
	return time.Since(c.LastActivity) < threshold
}

// GetConnectionInfo returns connection information
func (c *Client) GetConnectionInfo() map[string]interface{} {
	return map[string]interface{}{
		"session_id":    c.SessionID,
		"user_id":       c.UserID,
		"video_id":      c.VideoID,
		"user_agent":    c.UserAgent,
		"ip_address":    c.IPAddress,
		"last_activity": c.LastActivity,
		"connected_at":  time.Now().Sub(c.LastActivity),
	}
}

// Close gracefully closes the client connection
func (c *Client) Close() {
	close(c.send)
	c.conn.Close()
}

// ServeWS handles websocket requests from the peer
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request, sessionID, userID, videoID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	// Get client metadata
	userAgent := r.Header.Get("User-Agent")
	ipAddress := getClientIP(r)

	// Create new client
	client := NewClient(hub, conn, sessionID, userID, videoID, userAgent, ipAddress)

	// Register client with hub
	client.hub.register <- client

	// Start pumps in goroutines
	// Allow collection of memory referenced by the caller by doing all work in new goroutines
	go client.writePump()
	go client.readPump()
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Get first IP if there are multiple
		return forwarded
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}